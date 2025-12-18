// repository/cart_repository.go
package repo

import (
	"backend/model" // sesuaikan dengan path project lu
	"errors"

	"gorm.io/gorm"
)

type CartRepository interface {
	// Cart methods
	GetOrCreateCart(userID uint) (*model.Cart, error)
	GetCartByUserID(userID uint) (*model.Cart, error)
	GetCartWithItems(userID uint) (*model.Cart, error)
	DeleteCart(cartID uint) error

	// CartItem methods
	AddItemToCart(cartItem *model.CartItem) error
	GetCartItem(cartID, barangID uint) (*model.CartItem, error)
	GetCartItemByID(cartItemID, cartID uint) (*model.CartItem, error)
	GetCartItemByIDWithBarang(cartItemID, cartID uint) (*model.CartItem, error)
	UpdateCartItem(cartItem *model.CartItem) error
	DeleteCartItem(cartItemID uint) error
	DeleteCartItemByBarangID(cartID, barangID uint) error
	ClearCart(cartID uint) error
	GetCartItemCount(cartID uint) (int64, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

// ========== CART METHODS ==========

// GetOrCreateCart - Get cart atau create kalau belum ada
func (r *cartRepository) GetOrCreateCart(userID uint) (*model.Cart, error) {
	var cart model.Cart
	
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Belum ada cart, create baru
			cart = model.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}
	
	return &cart, nil
}

// GetCartByUserID - Get cart tanpa items
func (r *cartRepository) GetCartByUserID(userID uint) (*model.Cart, error) {
	var cart model.Cart
	
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	
	return &cart, nil
}

// GetCartWithItems - Get cart dengan semua items & barang details
func (r *cartRepository) GetCartWithItems(userID uint) (*model.Cart, error) {
	var cart model.Cart
	
	err := r.db.
		Preload("CartItems.Barang").
		Where("user_id = ?", userID).
		First(&cart).Error
	
	if err != nil {
		return nil, err
	}
	
	return &cart, nil
}

// DeleteCart - Hapus cart (cart items juga ke-delete karena CASCADE)
func (r *cartRepository) DeleteCart(cartID uint) error {
	return r.db.Delete(&model.Cart{}, cartID).Error
}

// ========== CART ITEM METHODS ==========

// AddItemToCart - Tambah item ke cart
func (r *cartRepository) AddItemToCart(cartItem *model.CartItem) error {
	return r.db.Create(cartItem).Error
}

// GetCartItem - Get specific cart item
func (r *cartRepository) GetCartItem(cartID, barangID uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	
	err := r.db.
		Where("cart_id = ? AND barang_id = ?", cartID, barangID).
		First(&cartItem).Error
	
	if err != nil {
		return nil, err
	}
	
	return &cartItem, nil
}

// UpdateCartItem - Update cart item (quantity, tanggal, dll)
func (r *cartRepository) UpdateCartItem(cartItem *model.CartItem) error {
	return r.db.Save(cartItem).Error
}

// DeleteCartItem - Hapus cart item by ID
func (r *cartRepository) DeleteCartItem(cartItemID uint) error {
	return r.db.Delete(&model.CartItem{}, cartItemID).Error
}

// DeleteCartItemByBarangID - Hapus cart item by barang ID
func (r *cartRepository) DeleteCartItemByBarangID(cartID, barangID uint) error {
	return r.db.
		Where("cart_id = ? AND barang_id = ?", cartID, barangID).
		Delete(&model.CartItem{}).Error
}

// ClearCart - Hapus semua items di cart
func (r *cartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}

// GetCartItemCount - Hitung jumlah item di cart
func (r *cartRepository) GetCartItemCount(cartID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.CartItem{}).Where("cart_id = ?", cartID).Count(&count).Error
	return count, err
}

// GetCartItemByID - Get cart item by ID dengan validasi cart ownership
func (r *cartRepository) GetCartItemByID(cartItemID, cartID uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	
	err := r.db.
		Where("id = ? AND cart_id = ?", cartItemID, cartID).
		First(&cartItem).Error
	
	if err != nil {
		return nil, err
	}
	
	return &cartItem, nil
}

// GetCartItemByIDWithBarang - Get cart item dengan preload barang
func (r *cartRepository) GetCartItemByIDWithBarang(cartItemID, cartID uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	
	err := r.db.
		Preload("Barang").
		Where("id = ? AND cart_id = ?", cartItemID, cartID).
		First(&cartItem).Error
	
	if err != nil {
		return nil, err
	}
	
	return &cartItem, nil
}