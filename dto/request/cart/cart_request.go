// dto/request/cart/cart_request.go
package cart

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ========== REQUEST DTOs ==========

// AddToCartRequest - ✅ SEKARANG PAKAI KODE BARANG
type AddToCartRequest struct {
	KodeBarang string `json:"kode_barang"` // ✅ Ganti dari BarangID
	Quantity   int    `json:"quantity"`
}

// BindAndValidate - Bind dan validasi AddToCartRequest
func (r *AddToCartRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON dulu
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	// Validation Ozzo
	return validation.ValidateStruct(r,
		validation.Field(&r.KodeBarang,
			validation.Required.Error("kode_barang wajib diisi"),
		),
		validation.Field(&r.Quantity,
			validation.Required.Error("quantity wajib diisi"),
			validation.Min(1).Error("quantity minimal 1"),
		),
	)
}

// UpdateCartItemRequest - Request untuk update cart item
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
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

// ✅ ENDPOINT BARU: Accept Suggestion
type AcceptSuggestionRequest struct {
	Adjustments []AdjustmentItem `json:"adjustments"`
}

// BindAndValidate - Bind dan validasi AcceptSuggestionRequest
func (r *AcceptSuggestionRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON dulu
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	// Validation Ozzo
	return validation.ValidateStruct(r,
		validation.Field(&r.Adjustments,
			validation.Required.Error("adjustments wajib diisi"),
			validation.Length(1, 100).Error("adjustments minimal 1 item, maksimal 100"),
		),
	)
}

type AdjustmentItem struct {
	CartItemID  uint `json:"cart_item_id"`
	AcceptedQty int  `json:"accepted_qty"`
}

// Validate AdjustmentItem
func (a AdjustmentItem) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CartItemID,
			validation.Required.Error("cart_item_id wajib diisi"),
		),
		validation.Field(&a.AcceptedQty,
			validation.Required.Error("accepted_qty wajib diisi"),
			validation.Min(0).Error("accepted_qty minimal 0"),
		),
	)
}