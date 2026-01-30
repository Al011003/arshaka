// ========== RESPONSE DTOs ==========
package cart

// dto/response/cart/cart.go

type BarangInCartResponse struct {
	ID       uint   `json:"id"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	Merk     string `json:"merk"`
	Kategori string `json:"kategori"`
	CoverURL string `json:"cover_url"`
	Status   string `json:"status"`
	
	// ✅ Unit-based (bukan stok lagi)
	TotalUnits     int  `json:"total_units"`
	AvailableUnits int  `json:"available_units"`
	IsAvailable    bool `json:"is_available"`
	Issue          string `json:"issue,omitempty"`
}

// CartItemResponse - Response untuk single cart item
type CartItemResponse struct {
	ID       uint                  `json:"id"`
	CartID   uint                  `json:"cart_id"`
	BarangID uint                  `json:"barang_id"`
	Barang   *BarangInCartResponse `json:"barang"`
	Quantity int                   `json:"quantity"`

	

}

// CartResponse - Response untuk cart dengan semua items
type CartResponse struct {
	ID         uint               `json:"id"`
	UserID     uint               `json:"user_id"`
	TotalItems int                `json:"total_items"`
	CartItems  []CartItemResponse `json:"cart_items"`

	Adjustments *CartAdjustmentResponse `json:"adjustments,omitempty"`
}

// CartSummaryResponse - Response ringkas untuk cart (tanpa detail items)
type CartSummaryResponse struct {
	ID         uint  `json:"id"`
	UserID     uint  `json:"user_id"`
	TotalItems int64 `json:"total_items"`
}

type CartAdjustmentResponse struct {
    RemovedItems    []RemovedItem       `json:"removed_items"`
    QuantityChanged []QuantityChanged   `json:"quantity_changed"`
}

type QuantityChanged struct {
	BarangID uint `json:"barang_id"`
	Nama     string `json:"nama"`
	From     int `json:"from"`
	To       int `json:"to"`
	Reason   string `json:"reason"`
}

type RemovedItem struct {
	BarangID uint `json:"barang_id"`
	Nama     string `json:"nama"`
	Reason   string `json:"reason"`
}