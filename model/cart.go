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
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	
	CartItems []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"cart_items,omitempty"`
}

func (Cart) TableName() string {
	return "carts"
}

// CartItem model - Simple: cuma barang + quantity
type CartItem struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	CartID   uint      `json:"cart_id" gorm:"not null;index"`
	BarangID uint      `json:"barang_id" gorm:"not null;index"`
	Quantity int       `json:"quantity" gorm:"not null;default:1"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Barang *Barang `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
	
	// ❌ REMOVE: StartDate, EndDate, Reason
	// Ini nanti diisi pas checkout jadi loan
}

func (CartItem) TableName() string {
	return "cart_items"
}

// BeforeCreate - validasi sebelum create cart item
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.Quantity <= 0 {
		return errors.New("quantity harus lebih dari 0")
	}
	return nil
}

// BeforeUpdate - validasi sebelum update cart item
func (ci *CartItem) BeforeUpdate(tx *gorm.DB) error {
	if ci.Quantity <= 0 {
		return errors.New("quantity harus lebih dari 0")
	}
	return nil
}