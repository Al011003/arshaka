// usecase/cart_usecase.go
package usecase

import (
	req "backend/dto/request/cart"
	res "backend/dto/response/cart"
	"backend/model"
	repository "backend/repo"
	"errors"

	"gorm.io/gorm"
)

type CartUsecase interface {
	GetMyCart(userID uint) (*res.CartResponse, error)
	AddToCart(userID uint, req req.AddToCartRequest) (*res.CartItemResponse, error)
	UpdateCartItem(userID, cartItemID uint, req req.UpdateCartItemRequest) (*res.CartItemResponse, error)
	RemoveFromCart(userID, cartItemID uint) error
	ClearMyCart(userID uint) error
	GetCartItemCount(userID uint) (*res.CartSummaryResponse, error)
}

type cartUsecase struct {
	cartRepo   repository.CartRepository
	barangRepo repository.BarangRepository
	unitRepo   repository.BarangUnitRepository // ← TAMBAH INI
}

func NewCartUsecase(
	cartRepo repository.CartRepository,
	barangRepo repository.BarangRepository,
	unitRepo repository.BarangUnitRepository, // ← TAMBAH INI
) CartUsecase {
	return &cartUsecase{
		cartRepo:   cartRepo,
		barangRepo: barangRepo,
		unitRepo:   unitRepo,
	}
}

// GetMyCart - Get cart user dengan semua items
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
	adjusted := []res.QuantityChanged{}
	validItems := make([]model.CartItem, 0)

	for _, item := range cart.CartItems {

		// ===============================
		// BARANG DIHAPUS / NONAKTIF
		// ===============================
		if item.Barang == nil || item.Barang.Status == "nonaktif" {
			_ = u.cartRepo.DeleteCartItem(item.ID)

			barangNama := "Unknown"
			if item.Barang != nil {
				barangNama = item.Barang.Nama
			}

			removed = append(removed, res.RemovedItem{
				BarangID: item.BarangID,
				Nama:     barangNama,
				Reason:   "BARANG_TIDAK_TERSEDIA",
			})
			continue
		}

		// ===============================
		// CEK TOTAL UNITS (pengganti StokTotal)
		// ===============================
		totalUnits, err := u.unitRepo.CountTotal(item.BarangID)
		if err != nil {
			continue // Skip kalau error
		}

		// Total units = 0 → AUTO ADJUST
		if totalUnits == 0 && item.Quantity > 0 {
			oldQty := item.Quantity
			item.Quantity = 0
			_ = u.cartRepo.UpdateCartItem(&item)

			adjusted = append(adjusted, res.QuantityChanged{
				BarangID: item.BarangID,
				Nama:     item.Barang.Nama,
				From:     oldQty,
				To:       0,
				Reason:   "TIDAK_ADA_UNIT",
			})
		}

		// Total units berkurang
		if item.Quantity > int(totalUnits) {
			oldQty := item.Quantity
			item.Quantity = int(totalUnits)
			_ = u.cartRepo.UpdateCartItem(&item)

			adjusted = append(adjusted, res.QuantityChanged{
				BarangID: item.BarangID,
				Nama:     item.Barang.Nama,
				From:     oldQty,
				To:       item.Quantity,
				Reason:   "JUMLAH_UNIT_BERUBAH",
			})
		}

		validItems = append(validItems, item)
	}

	cart.CartItems = validItems
	resp := toCartResponse(cart, u.unitRepo)

	if len(removed) > 0 || len(adjusted) > 0 {
		resp.Adjustments = &res.CartAdjustmentResponse{
			RemovedItems:    removed,
			QuantityChanged: adjusted,
		}
	}

	return resp, nil
}

// AddToCart - Tambah barang ke cart
func (u *cartUsecase) AddToCart(
	userID uint,
	req req.AddToCartRequest,
) (*res.CartItemResponse, error) {

	// 1. Validasi barang exists & aktif
	barang, err := u.barangRepo.FindByID(req.BarangID)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	if barang.Status == "nonaktif" {
		return nil, errors.New("barang tidak aktif")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity tidak valid")
	}

	// 2. Cek total units (pengganti stok total)
	totalUnits, err := u.unitRepo.CountTotal(req.BarangID)
	if err != nil {
		return nil, err
	}

	if totalUnits == 0 {
		return nil, errors.New("tidak ada unit tersedia untuk barang ini")
	}

	if req.Quantity > int(totalUnits) {
		return nil, errors.New("jumlah unit tidak mencukupi")
	}

	// 3. Get or create cart
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// 4. Cek apakah item sudah ada
	existing, _ := u.cartRepo.GetCartItem(cart.ID, barang.ID)
	
	if existing != nil {
		// Item sudah ada, update quantity
		newQty := existing.Quantity + req.Quantity

		if newQty > int(totalUnits) {
			return nil, errors.New("jumlah unit tidak mencukupi")
		}

		existing.Quantity = newQty
		if err := u.cartRepo.UpdateCartItem(existing); err != nil {
			return nil, err
		}

		item, _ := u.cartRepo.GetCartItemByIDWithBarang(existing.ID, cart.ID)
		return toCartItemResponse(item, u.unitRepo), nil
	}

	// 5. Create baru
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

// UpdateCartItem - Update cart item
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
		return nil, errors.New("cart tidak ditemukan")
	}

	// 2. Get cart item (ownership check)
	cartItem, err := u.cartRepo.GetCartItemByID(cartItemID, cart.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item tidak ditemukan di cart")
		}
		return nil, err
	}

	// 3. Get barang (buat cek status)
	barang, err := u.barangRepo.FindByID(cartItem.BarangID)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	if barang.Status == "nonaktif" {
		return nil, errors.New("barang tidak aktif")
	}

	// 4. Cek total units
	totalUnits, err := u.unitRepo.CountTotal(cartItem.BarangID)
	if err != nil {
		return nil, err
	}

	if req.Quantity > int(totalUnits) {
		return nil, errors.New("jumlah unit tidak mencukupi")
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

// ========== HELPER FUNCTIONS (Converters) ==========

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
		// Get unit stats
		totalUnits, _ := unitRepo.CountTotal(item.BarangID)
		availableUnits, _ := unitRepo.CountAvailable(item.BarangID)
		
		barangResp = &res.BarangInCartResponse{
			ID:       item.Barang.ID,
			Kode:     item.Barang.Kode,
			Nama:     item.Barang.Nama,
			Merk:     item.Barang.Merk,
			Kategori: item.Barang.Kategori,
			CoverURL: item.Barang.CoverURL,
			Status:   item.Barang.Status,
			
			// ✅ Unit-based info
			TotalUnits:     int(totalUnits),
			AvailableUnits: int(availableUnits),
			IsAvailable:    true,
		}

		// ===============================
		// RUNTIME ISSUE DETECTION
		// ===============================
		if item.Barang.Status == "nonaktif" {
			barangResp.IsAvailable = false
			barangResp.Issue = "BARANG_NONAKTIF"
		} else if totalUnits == 0 {
			barangResp.IsAvailable = false
			barangResp.Issue = "TIDAK_ADA_UNIT"
		} else if item.Quantity > int(totalUnits) {
			barangResp.IsAvailable = false
			barangResp.Issue = "JUMLAH_UNIT_TIDAK_CUKUP"
		} else if availableUnits == 0 {
			// Semua unit dipinjam (tapi ini cuma warning, bukan error)
			// Karena belum tau tanggal pinjam
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