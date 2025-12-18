// ========== RESPONSE DTOs ==========
package cart

// BarangInCartResponse - Info barang di cart (nested)
type BarangInCartResponse struct {
	ID       uint   `json:"id"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	Merk     string `json:"merk"`
	Kategori string `json:"kategori"`
	StokTotal int   `json:"stok_total"`
	StokSisa  int   `json:"stok_sisa"`
	CoverURL  string `json:"cover_url"`
	Status    string `json:"status"`
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
}

// CartSummaryResponse - Response ringkas untuk cart (tanpa detail items)
type CartSummaryResponse struct {
	ID         uint  `json:"id"`
	UserID     uint  `json:"user_id"`
	TotalItems int64 `json:"total_items"`
}