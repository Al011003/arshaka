// dto/request/cart/cart_request.go
package cart

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ========== REQUEST DTOs ==========

// AddToCartRequest - Request untuk tambah item ke cart
type AddToCartRequest struct {
	BarangID uint   `json:"barang_id"`
	Quantity int    `json:"quantity"`
}

// BindAndValidate - Bind dan validasi AddToCartRequest
func (r *AddToCartRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON dulu
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	// Validation Ozzo
	return validation.ValidateStruct(r,
		validation.Field(&r.BarangID,
			validation.Required.Error("barang_id wajib diisi"),
		),
		validation.Field(&r.Quantity,
			validation.Required.Error("quantity wajib diisi"),
			validation.Min(1).Error("quantity minimal 1"),
		),
	)
}

// UpdateCartItemRequest - Request untuk update cart item
type UpdateCartItemRequest struct {
	Quantity int    `json:"quantity"`
}

// BindAndValidate - Bind dan validasi UpdateCartItemRequest
func (r *UpdateCartItemRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON dulu
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	// Validation Ozzo
	return validation.ValidateStruct(r,
		validation.Field(&r.Quantity,
			validation.Required.Error("quantity wajib diisi"),
			validation.Min(1).Error("quantity minimal 1"),
		),
	)
}

// ========== RESPONS