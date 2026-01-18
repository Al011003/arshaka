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
}

func NewCartUsecase(
	cartRepo repository.CartRepository,
	barangRepo repository.BarangRepository,
) CartUsecase {
	return &cartUsecase{
		cartRepo:   cartRepo,
		barangRepo: barangRepo,
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

			removed = append(removed, res.RemovedItem{
				BarangID: item.BarangID,
				Nama:     item.BarangNama,
				Reason:   "BARANG_TIDAK_TERSEDIA",
			})
			continue
		}

		// ===============================
		// STOK TOTAL = 0 → AUTO ADJUST
		// ===============================
		if item.Barang.StokTotal == 0 && item.Quantity > 0 {
			oldQty := item.Quantity
			item.Quantity = 0
			_ = u.cartRepo.UpdateCartItem(&item)

			adjusted = append(adjusted, res.QuantityChanged{
				BarangID: item.BarangID,
				Nama:     item.BarangNama,
				From:     oldQty,
				To:       0,
				Reason:   "STOK_KOSONG",
			})
		}

		// ===============================
		// STOK TOTAL BERKURANG
		// ===============================
		if item.Quantity > item.Barang.StokTotal {
			oldQty := item.Quantity
			item.Quantity = item.Barang.StokTotal
			_ = u.cartRepo.UpdateCartItem(&item)

			adjusted = append(adjusted, res.QuantityChanged{
				BarangID: item.BarangID,
				Nama:     item.BarangNama,
				From:     oldQty,
				To:       item.Quantity,
				Reason:   "STOK_TOTAL_BERUBAH",
			})
		}

		validItems = append(validItems, item)
	}

	cart.CartItems = validItems
	resp := toCartResponse(cart)

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

	// ===============================
	// VALIDASI STOK (STRICT)
	// ===============================
	if req.Quantity > barang.StokTotal {
		return nil, errors.New("stok barang tidak mencukupi")
	}

	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// ===============================
	// ITEM SUDAH ADA
	// ===============================
	existing, err := u.cartRepo.GetCartItem(cart.ID, barang.ID)
	if err == nil {

		newQty := existing.Quantity + req.Quantity

		if newQty > barang.StokTotal {
			return nil, errors.New("stok barang tidak mencukupi")
		}

		existing.Quantity = newQty
		if err := u.cartRepo.UpdateCartItem(existing); err != nil {
			return nil, err
		}

		item, _ := u.cartRepo.GetCartItemByIDWithBarang(existing.ID, cart.ID)
		return toCartItemResponse(item), nil
	}

	// ===============================
	// CREATE BARU
	// ===============================
	cartItem := &model.CartItem{
		CartID:     cart.ID,
		BarangID:   barang.ID,
		BarangNama: barang.Nama, // snapshot
		BarangKode: barang.Kode, // snapshot
		Quantity:   req.Quantity,
	}

	if err := u.cartRepo.AddItemToCart(cartItem); err != nil {
		return nil, err
	}

	item, _ := u.cartRepo.GetCartItemByIDWithBarang(cartItem.ID, cart.ID)
	return toCartItemResponse(item), nil
}



// UpdateCartItem - Update cart item
func (u *cartUsecase) UpdateCartItem(userID, cartItemID uint, req req.UpdateCartItemRequest) (*res.CartItemResponse, error) {

	if req.Quantity <= 0 {
		return nil, errors.New("quantity harus lebih dari 0")
	}

	// 2. Get cart user
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart tidak ditemukan")
	}

	// 3. Get cart item (ownership check)
	cartItem, err := u.cartRepo.GetCartItemByID(cartItemID, cart.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item tidak ditemukan di cart")
		}
		return nil, err
	}

	// 4. Get barang (buat cek stok & status)
	barang, err := u.barangRepo.FindByID(cartItem.BarangID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("barang tidak ditemukan")
		}
		return nil, err
	}

	// 5. Validasi status barang
	if barang.Status == "nonaktif" {
		return nil, errors.New("barang tidak aktif")
	}

	// 6. VALIDASI STOK
	if req.Quantity > barang.StokTotal {
		return nil, errors.New("stok barang tidak mencukupi")
	}

	// 7. Update quantity
	cartItem.Quantity = req.Quantity
	if err := u.cartRepo.UpdateCartItem(cartItem); err != nil {
		return nil, err
	}

	// 8. Return updated item
	updatedItem, err := u.cartRepo.GetCartItemByIDWithBarang(cartItem.ID, cart.ID)
	if err != nil {
		return nil, err
	}

	return toCartItemResponse(updatedItem), nil
}

// RemoveFromCart - Hapus item dari cart
func (u *cartUsecase) RemoveFromCart(userID, cartItemID uint) error {
	// 1. Get cart user
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return errors.New("cart tidak ditemukan")
	}

	// 2. Validasi cart item belongs to user
	_, err = u.cartRepo.GetCartItemByID(cartItemID, cart.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item tidak ditemukan di cart")
		}
		return err
	}

	// 3. Delete cart item
	return u.cartRepo.DeleteCartItem(cartItemID)
}

// ClearMyCart - Hapus semua items di cart
func (u *cartUsecase) ClearMyCart(userID uint) error {
	cart, err := u.cartRepo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Cart belum ada, gak perlu clear
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
			// Cart belum ada, return empty summary
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

// toCartResponse - Convert model.Cart ke dto.CartResponse
func toCartResponse(cart *model.Cart) *res.CartResponse {
	cartItems := make([]res.CartItemResponse, 0, len(cart.CartItems))
	
	for _, item := range cart.CartItems {
		cartItems = append(cartItems, *toCartItemResponse(&item))
	}

	return &res.CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		TotalItems: len(cart.CartItems),
		CartItems:  cartItems,
	}
}

// toCartItemResponse - Convert model.CartItem ke dto.CartItemResponse

func mapBarangToCartResponse(b *model.Barang, qty int) *res.BarangInCartResponse {
	resp := &res.BarangInCartResponse{
		ID:          b.ID,
		Kode:        b.Kode,
		Nama:        b.Nama,
		Merk:        b.Merk,
		Kategori:    b.Kategori,
		StokTotal:   b.StokTotal,
		StokSisa:    b.StokSisa,
		CoverURL:    b.CoverURL,
		Status:      b.Status,
		IsAvailable: true,
	}

	if b.Status != "aktif" {
		resp.IsAvailable = false
		resp.Issue = "BARANG_NONAKTIF"
	} else if b.StokSisa == 0 {
		resp.IsAvailable = false
		resp.Issue = "STOK_HABIS"
	} else if b.StokSisa < qty {
		resp.IsAvailable = false
		resp.Issue = "STOK_TIDAK_CUKUP"
	}

	return resp
}

func toCartItemResponse(item *model.CartItem) *res.CartItemResponse {
	var barangResp *res.BarangInCartResponse

	if item.Barang != nil {
		barangResp = &res.BarangInCartResponse{
			ID:          item.Barang.ID,
			Kode:        item.Barang.Kode,
			Nama:        item.Barang.Nama,
			Merk:        item.Barang.Merk,
			Kategori:    item.Barang.Kategori,
			StokTotal:   item.Barang.StokTotal,
			StokSisa:    item.Barang.StokSisa,
			CoverURL:    item.Barang.CoverURL,
			Status:      item.Barang.Status,

			// 🔥 PENTING
			IsAvailable: true,
		}

		// ===============================
		// RUNTIME ISSUE DETECTION
		// ===============================
		if item.Barang.Status == "nonaktif" {
			barangResp.IsAvailable = false
			barangResp.Issue = "BARANG_NONAKTIF"
		}

		if item.Quantity > item.Barang.StokTotal {
			barangResp.IsAvailable = false
			barangResp.Issue = "STOK_TOTAL_TIDAK_CUKUP"
		}

		if item.Barang.StokTotal == 0 {
			barangResp.IsAvailable = false
			barangResp.Issue = "STOK_KOSONG"
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
