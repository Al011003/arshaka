// model/cart.go
package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Cart model
type Cart struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	CartItems []CartItem `json:"cart_items,omitempty" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// CartItem model - Simplified (NO dates)
type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"not null;index"`
	Cart      *Cart     `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	BarangID  uint      `json:"barang_id" gorm:"not null;index"`
	Barang    *Barang   `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Cart) TableName() string {
	return "carts"
}

func (CartItem) TableName() string {
	return "cart_items"
}

// BeforeCreate - validasi sebelum create cart item
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	// Validasi quantity
	if ci.Quantity <= 0 {
		return errors.New("quantity harus lebih dari 0")
	}

	return nil
}

// BeforeUpdate - validasi sebelum update cart item
func (ci *CartItem) BeforeUpdate(tx *gorm.DB) error {
	// Validasi quantity
	if ci.Quantity <= 0 {
		return errors.New("quantity harus lebih dari 0")
	}

	return nil
}