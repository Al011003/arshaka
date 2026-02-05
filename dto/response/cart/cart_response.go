// ========== RESPONSE DTOs ==========
package cart

// dto/response/cart/cart.go



type QuantityChanged struct {
	BarangID uint `json:"barang_id"`
	Nama     string `json:"nama"`
	From     int `json:"from"`
	To       int `json:"to"`
	Reason   string `json:"reason"`
}


// dto/response/cart/cart_response.go

// CartResponse - Response utama cart
type CartResponse struct {
	ID          uint                      `json:"id"`
	UserID      uint                      `json:"user_id"`
	TotalItems  int                       `json:"total_items"`
	CartItems   []CartItemResponse        `json:"cart_items"`
	Adjustments *CartAdjustmentResponse   `json:"adjustments,omitempty"` // ✅ Nullable
}

// CartItemResponse - Response per item
type CartItemResponse struct {
	ID       uint                  `json:"id"`
	CartID   uint                  `json:"cart_id"`
	BarangID uint                  `json:"barang_id"`
	Barang   *BarangInCartResponse `json:"barang"`
	Quantity int                   `json:"quantity"`
}

// BarangInCartResponse - Info barang di cart
type BarangInCartResponse struct {
	ID             uint   `json:"id"`
	Kode           string `json:"kode"`
	Nama           string `json:"nama"`
	Merk           string `json:"merk"`
	Kategori       string `json:"kategori"`
	CoverURL       string `json:"cover_url"`
	Status         string `json:"status"`
	TotalUnits     int    `json:"total_units"`
	AvailableUnits int    `json:"available_units"`
	IsAvailable    bool   `json:"is_available"`
	Issue          string `json:"issue,omitempty"`
}

// ✅ ADJUSTMENT RESPONSE (Auto-adjust + Suggestions)
type CartAdjustmentResponse struct {
	RemovedItems []RemovedItem     `json:"removed_items,omitempty"`
	Suggestions  []ItemSuggestion  `json:"suggestions,omitempty"`
}

// RemovedItem - Barang yang dihapus otomatis
type RemovedItem struct {
	BarangID uint   `json:"barang_id"`
	Nama     string `json:"nama"`
	Reason   string `json:"reason"` // BARANG_DIHAPUS, BARANG_NONAKTIF, TIDAK_ADA_UNIT
}

// ✅ ItemSuggestion - Saran qty adjustment
type ItemSuggestion struct {
	CartItemID   uint   `json:"cart_item_id"`
	BarangID     uint   `json:"barang_id"`
	Nama         string `json:"nama"`
	RequestedQty int    `json:"requested_qty"`
	AvailableQty int    `json:"available_qty"`
	Status       string `json:"status"`  // AVAILABLE, PARTIAL_AVAILABLE, NOT_AVAILABLE
	Message      string `json:"message"` // "3 unit sedang maintenance, 2 unit sedang dipinjam"
}

// CartSummaryResponse - Summary cart (count only)
type CartSummaryResponse struct {
	ID         uint `json:"id"`
	UserID     uint `json:"user_id"`
	TotalItems int64  `json:"total_items"`
}