// usecase/cart_usecase.go
package usecase

import (
	req "backend/dto/request/cart"
	res "backend/dto/response/cart"
	"backend/model"
	repository "backend/repo"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CartUsecase interface {
	GetMyCart(userID uint) (*res.CartResponse, error)
	AddToCart(userID uint, req req.AddToCartRequest) (*res.CartItemResponse, error)
	UpdateCartItem(userID, cartItemID uint, req req.UpdateCartItemRequest) (*res.CartItemResponse, error)
	RemoveFromCart(userID, cartItemID uint) error
	ClearMyCart(userID uint) error
	GetCartItemCount(userID uint) (*res.CartSummaryResponse, error)
	
	// ✅ ENDPOINT BARU: Accept Suggestion
	AcceptSuggestion(userID uint, req req.AcceptSuggestionRequest) (*res.CartResponse, error)
}

type cartUsecase struct {
	cartRepo   repository.CartRepository
	barangRepo repository.BarangRepository
	unitRepo   repository.BarangUnitRepository
}

func NewCartUsecase(
	cartRepo repository.CartRepository,
	barangRepo repository.BarangRepository,
	unitRepo repository.BarangUnitRepository,
) CartUsecase {
	return &cartUsecase{
		cartRepo:   cartRepo,
		barangRepo: barangRepo,
		unitRepo:   unitRepo,
	}
}

// ========================================
// GetMyCart - DENGAN DETAIL UNIT CHECKING
// ========================================
func (u *cartUsecase) GetMyCart(userID uint) (*res.CartResponse, error) {
	cart, err := u.cartRepo.GetCartWithItems(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart, err = u.cartRepo.GetOrCreateCart(userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	removed := []res.RemovedItem{}
	suggestions := []res.ItemSuggestion{}
	validItems := make([]model.CartItem, 0)

	for _, item := range cart.CartItems {

		// ===============================
		// 1. AUTO-DELETE (Tanpa ACC)
		// ===============================
		
		// Barang dihapus/soft deleted dari database
		if item.Barang == nil {
			_ = u.cartRepo.DeleteCartItem(item.ID)
			removed = append(removed, res.RemovedItem{
				BarangID: item.BarangID,
				Nama:     "Unknown",
				Reason:   "BARANG_DIHAPUS",
			})
			continue
		}

		// Tidak ada unit sama sekali (semua dihapus)
		totalUnits, _ := u.unitRepo.CountTotal(item.BarangID)
		if totalUnits == 0 {
			_ = u.cartRepo.DeleteCartItem(item.ID)
			removed = append(removed, res.RemovedItem{
				BarangID: item.BarangID,
				Nama:     item.Barang.Nama,
				Reason:   "TIDAK_ADA_UNIT",
			})
			continue
		}

		// ===============================
		// 2. SUGGESTION (Perlu ACC User)
		// ===============================
		
		// Barang jadi nonaktif (ITEM TETAP ADA!)
		if item.Barang.Status == "nonaktif" {
			suggestions = append(suggestions, res.ItemSuggestion{
				CartItemID:   item.ID,
				BarangID:     item.BarangID,
				Nama:         item.Barang.Nama,
				RequestedQty: item.Quantity,
				AvailableQty: 0,
				Status:       "NOT_AVAILABLE",
				Message:      "Barang sedang nonaktif",
			})
			validItems = append(validItems, item)
			continue
		}

		// Cek unit availability
		unitCheck := u.checkUnitAvailability(item.BarangID, item.Quantity)

		// Kalau qty request > available (maintenance/dipinjam)
		if item.Quantity > unitCheck.AvailableQty {
			suggestions = append(suggestions, res.ItemSuggestion{
				CartItemID:   item.ID,
				BarangID:     item.BarangID,
				Nama:         item.Barang.Nama,
				RequestedQty: item.Quantity,
				AvailableQty: unitCheck.AvailableQty,
				Status:       unitCheck.Status,
				Message:      unitCheck.Message,
			})
		}

		// Item tetap di cart (ga dihapus)
		validItems = append(validItems, item)
	}

	cart.CartItems = validItems
	resp := toCartResponse(cart, u.unitRepo)

	// Attach adjustments kalau ada
	if len(removed) > 0 || len(suggestions) > 0 {
		resp.Adjustments = &res.CartAdjustmentResponse{
			RemovedItems: removed,
			Suggestions:  suggestions,
		}
	}

	return resp, nil
}

// ========================================
// AddToCart - PAKAI KODE BARANG
// ========================================
func (u *cartUsecase) AddToCart(
	userID uint,
	req req.AddToCartRequest,
) (*res.CartItemResponse, error) {

	// ✅ 1. VALIDASI: Cari barang by KODE (bukan ID)
	if req.KodeBarang == "" {
		return nil, errors.New("kode barang harus diisi")
	}

	barang, err := u.barangRepo.FindByKode(req.KodeBarang)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("barang tidak ditemukan")
		}
		return nil, err
	}

	if barang.Status == "nonaktif" {
		return nil, errors.New("barang tidak aktif")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity tidak valid")
	}

	// ✅ 2. CEK UNIT AVAILABILITY (detailed)
	unitCheck := u.checkUnitAvailability(barang.ID, req.Quantity)

	if unitCheck.TotalUnits == 0 {
		return nil, errors.New("tidak ada unit tersedia untuk barang ini")
	}

	if unitCheck.AvailableQty < req.Quantity {
		return nil, fmt.Errorf(
			"hanya tersedia %d unit (%s)",
			unitCheck.AvailableQty,
			unitCheck.Message,
		)
	}

	// ✅ 3. Get or create cart
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// ✅ 4. Cek apakah item sudah ada
	existing, _ := u.cartRepo.GetCartItem(cart.ID, barang.ID)

	if existing != nil {
		// Item sudah ada, update quantity
		newQty := existing.Quantity + req.Quantity

		// Re-check dengan qty baru
		unitCheck = u.checkUnitAvailability(barang.ID, newQty)
		if unitCheck.AvailableQty < newQty {
			return nil, fmt.Errorf(
				"hanya tersedia %d unit (%s)",
				unitCheck.AvailableQty,
				unitCheck.Message,
			)
		}

		existing.Quantity = newQty
		if err := u.cartRepo.UpdateCartItem(existing); err != nil {
			return nil, err
		}

		item, _ := u.cartRepo.GetCartItemByIDWithBarang(existing.ID, cart.ID)
		return toCartItemResponse(item, u.unitRepo), nil
	}

	// ✅ 5. Create baru
	cartItem := &model.CartItem{
		CartID:   cart.ID,
		BarangID: barang.ID,
		Quantity: req.Quantity,
	}

	if err := u.cartRepo.AddItemToCart(cartItem); err != nil {
		return nil, err
	}

	item, _ := u.cartRepo.GetCartItemByIDWithBarang(cartItem.ID, cart.ID)
	return toCartItemResponse(item, u.unitRepo), nil
}

// ========================================
// UpdateCartItem - FIX LOGIC
// ========================================
func (u *cartUsecase) UpdateCartItem(
	userID, cartItemID uint,
	req req.UpdateCartItemRequest,
) (*res.CartItemResponse, error) {

	if req.Quantity <= 0 {
		return nil, errors.New("quantity harus lebih dari 0")
	}

	// 1. Get cart user
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart tidak ditemukan")
		}
		return nil, err
	}

	// 2. Get cart item (ownership check)
	cartItem, err := u.cartRepo.GetCartItemByID(cartItemID, cart.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item tidak ditemukan di cart")
		}
		return nil, err
	}

	// 3. Get barang
	barang, err := u.barangRepo.FindByID(cartItem.BarangID)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	if barang.Status == "nonaktif" {
		return nil, errors.New("barang tidak aktif")
	}

	// ✅ 4. CEK UNIT AVAILABILITY
	unitCheck := u.checkUnitAvailability(cartItem.BarangID, req.Quantity)

	if unitCheck.AvailableQty < req.Quantity {
		return nil, fmt.Errorf(
			"hanya tersedia %d unit (%s)",
			unitCheck.AvailableQty,
			unitCheck.Message,
		)
	}

	// 5. Update quantity
	cartItem.Quantity = req.Quantity
	if err := u.cartRepo.UpdateCartItem(cartItem); err != nil {
		return nil, err
	}

	// 6. Return updated item
	updatedItem, err := u.cartRepo.GetCartItemByIDWithBarang(cartItem.ID, cart.ID)
	if err != nil {
		return nil, err
	}

	return toCartItemResponse(updatedItem, u.unitRepo), nil
}

// ========================================
// AcceptSuggestion - ENDPOINT BARU
// ========================================
func (u *cartUsecase) AcceptSuggestion(
	userID uint,
	req req.AcceptSuggestionRequest,
) (*res.CartResponse, error) {

	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart tidak ditemukan")
	}

	// Loop semua adjustments yang user terima
	for _, adjustment := range req.Adjustments {
		cartItem, err := u.cartRepo.GetCartItemByID(adjustment.CartItemID, cart.ID)
		if err != nil {
			continue // Skip kalau tidak ketemu
		}

		// Update ke qty yang disarankan
		cartItem.Quantity = adjustment.AcceptedQty
		_ = u.cartRepo.UpdateCartItem(cartItem)
	}

	// Return updated cart
	return u.GetMyCart(userID)
}

// RemoveFromCart - Hapus item dari cart
func (u *cartUsecase) RemoveFromCart(userID, cartItemID uint) error {
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return errors.New("cart tidak ditemukan")
	}

	_, err = u.cartRepo.GetCartItemByID(cartItemID, cart.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item tidak ditemukan di cart")
		}
		return err
	}

	return u.cartRepo.DeleteCartItem(cartItemID)
}

// ClearMyCart - Hapus semua items di cart
func (u *cartUsecase) ClearMyCart(userID uint) error {
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return u.cartRepo.ClearCart(cart.ID)
}

// GetCartItemCount - Get jumlah item di cart
func (u *cartUsecase) GetCartItemCount(userID uint) (*res.CartSummaryResponse, error) {
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &res.CartSummaryResponse{
				UserID:     userID,
				TotalItems: 0,
			}, nil
		}
		return nil, err
	}

	count, err := u.cartRepo.GetCartItemCount(cart.ID)
	if err != nil {
		return nil, err
	}

	return &res.CartSummaryResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		TotalItems: count,
	}, nil
}

// ========================================
// HELPER: CHECK UNIT AVAILABILITY
// ========================================
type UnitAvailabilityCheck struct {
	TotalUnits   int64
	AvailableQty int
	Status       string // "AVAILABLE", "PARTIAL_AVAILABLE", "NOT_AVAILABLE"
	Message      string
}

func (u *cartUsecase) checkUnitAvailability(
	barangID uint,
	requestedQty int,
) UnitAvailabilityCheck {

	result := UnitAvailabilityCheck{
		TotalUnits:   0,
		AvailableQty: 0,
		Status:       "NOT_AVAILABLE",
		Message:      "",
	}

	// 1. Get total units
	totalUnits, err := u.unitRepo.CountTotal(barangID)
	if err != nil || totalUnits == 0 {
		result.Message = "Tidak ada unit tersedia"
		return result
	}
	result.TotalUnits = totalUnits

	// 2. Get counts by status (SESUAIKAN DENGAN STATUS LO!)
	availableCount, _ := u.unitRepo.CountAvailable(barangID) // status = "siap"
	
	// Cek status lain yang lo punya (sesuaikan!)
	maintCount, _ := u.unitRepo.CountByStatus(barangID, "maintenance_required")
	
	// Kalau lo ada status "dipinjam" atau "rusak", tambahin di sini
	// rusak, _ := u.unitRepo.CountByStatus(barangID, "rusak")
	
	result.AvailableQty = int(availableCount)

	// 3. Build message
	messages := []string{}
	if maintCount > 0 {
		messages = append(messages, fmt.Sprintf("%d unit sedang maintenance", maintCount))
	}
	
	// Hitung yang "tidak tersedia" (selain siap & maintenance)
	unavailableCount := totalUnits - availableCount - maintCount
	if unavailableCount > 0 {
		messages = append(messages, fmt.Sprintf("%d unit tidak tersedia", unavailableCount))
	}

	// 4. Determine status
	if availableCount >= int64(requestedQty) {
		result.Status = "AVAILABLE"
		if len(messages) > 0 {
			result.Message = fmt.Sprintf("Tersedia %d unit", availableCount)
		}
	} else if availableCount > 0 {
		result.Status = "PARTIAL_AVAILABLE"
		if len(messages) > 0 {
			result.Message = fmt.Sprintf("Hanya tersedia %d unit, %s",
				availableCount,
				joinMessages(messages),
			)
		} else {
			result.Message = fmt.Sprintf("Hanya tersedia %d unit", availableCount)
		}
	} else {
		result.Status = "NOT_AVAILABLE"
		result.Message = joinMessages(messages)
		if result.Message == "" {
			result.Message = "Tidak ada unit tersedia"
		}
	}

	return result
}

func joinMessages(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	if len(messages) == 1 {
		return messages[0]
	}
	
	result := messages[0]
	for i := 1; i < len(messages); i++ {
		result += ", " + messages[i]
	}
	return result
}

// ========================================
// HELPER FUNCTIONS (Converters)
// ========================================

func toCartResponse(cart *model.Cart, unitRepo repository.BarangUnitRepository) *res.CartResponse {
	cartItems := make([]res.CartItemResponse, 0, len(cart.CartItems))

	for _, item := range cart.CartItems {
		cartItems = append(cartItems, *toCartItemResponse(&item, unitRepo))
	}

	return &res.CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		TotalItems: len(cart.CartItems),
		CartItems:  cartItems,
	}
}

func toCartItemResponse(
	item *model.CartItem,
	unitRepo repository.BarangUnitRepository,
) *res.CartItemResponse {
	var barangResp *res.BarangInCartResponse

	if item.Barang != nil {
		totalUnits, _ := unitRepo.CountTotal(item.BarangID)
		availableUnits, _ := unitRepo.CountAvailable(item.BarangID)

		barangResp = &res.BarangInCartResponse{
			ID:             item.Barang.ID,
			Kode:           item.Barang.Kode,
			Nama:           item.Barang.Nama,
			Merk:           item.Barang.Merk,
			Kategori:       item.Barang.Kategori,
			CoverURL:       item.Barang.CoverURL,
			Status:         item.Barang.Status,
			TotalUnits:     int(totalUnits),
			AvailableUnits: int(availableUnits),
			IsAvailable:    true,
		}

		// Runtime issue detection
		if item.Barang.Status == "nonaktif" {
			barangResp.IsAvailable = false
			barangResp.Issue = "BARANG_NONAKTIF"
		} else if totalUnits == 0 {
			barangResp.IsAvailable = false
			barangResp.Issue = "TIDAK_ADA_UNIT"
		} else if item.Quantity > int(availableUnits) {
			barangResp.IsAvailable = false
			barangResp.Issue = "JUMLAH_UNIT_TIDAK_CUKUP"
		} else if availableUnits == 0 {
			barangResp.Issue = "SEMUA_UNIT_SEDANG_DIPINJAM"
		}
	}

	return &res.CartItemResponse{
		ID:       item.ID,
		CartID:   item.CartID,
		BarangID: item.BarangID,
		Barang:   barangResp,
		Quantity: item.Quantity,
	}
}