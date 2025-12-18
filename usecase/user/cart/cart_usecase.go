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
			// Kalau belum ada cart, create baru
			cart, err = u.cartRepo.GetOrCreateCart(userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	return toCartResponse(cart), nil
}

// AddToCart - Tambah barang ke cart
func (u *cartUsecase) AddToCart(userID uint, req req.AddToCartRequest) (*res.CartItemResponse, error) {
	// 1. Validasi barang exists
	barang, err := u.barangRepo.FindByID(req.BarangID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("barang tidak ditemukan")
		}
		return nil, err
	}

	// 2. Validasi status barang
	if barang.Status == "nonaktif" {
		return nil, errors.New("barang tidak aktif")
	}

	if req.Quantity > barang.StokTotal {
		return nil, errors.New("stok barang tidak mencukupi")
		}

	// 3. Get or create cart
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// 4. Cek apakah barang sudah ada di cart
	existingItem, err := u.cartRepo.GetCartItem(cart.ID, req.BarangID)
	if err == nil {
		// Barang sudah ada, update quantity
		existingItem.Quantity += req.Quantity
		totalQty := existingItem.Quantity + req.Quantity
			if totalQty > barang.StokTotal {
			return nil, errors.New("stok barang tidak mencukupi")
			}

	existingItem.Quantity = totalQty
		
		if err := u.cartRepo.UpdateCartItem(existingItem); err != nil {
			return nil, err
		}
		
		// Get updated item with barang
		updatedItem, err := u.cartRepo.GetCartItemByIDWithBarang(existingItem.ID, cart.ID)
		if err != nil {
			return nil, err
		}
		
		return toCartItemResponse(updatedItem), nil
	}

	// 5. Barang belum ada, create new cart item
	cartItem := &model.CartItem{
		CartID:   cart.ID,
		BarangID: req.BarangID,
		Quantity: req.Quantity,
	}

	if err := u.cartRepo.AddItemToCart(cartItem); err != nil {
		return nil, err
	}

	// Get cart item with barang info
	createdItem, err := u.cartRepo.GetCartItemByIDWithBarang(cartItem.ID, cart.ID)
	if err != nil {
		return nil, err
	}

	return toCartItemResponse(createdItem), nil
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
func toCartItemResponse(item *model.CartItem) *res.CartItemResponse {
	var barangResp *res.BarangInCartResponse
	
	if item.Barang != nil {
		barangResp = &res.BarangInCartResponse{
			ID:        item.Barang.ID,
			Kode:      item.Barang.Kode,
			Nama:      item.Barang.Nama,
			Merk:      item.Barang.Merk,
			Kategori:  item.Barang.Kategori,
			StokTotal: item.Barang.StokTotal,
			StokSisa:  item.Barang.StokSisa,
			CoverURL:  item.Barang.CoverURL,
			Status:    item.Barang.Status,
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