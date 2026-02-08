// dto/response/loan/loan.go
package loan

import "time"

// ========== LIST RESPONSE (untuk card/table) ==========
type LoanListResponse struct {
	LoanCode       string    `json:"loan_code"`
	LoanStatus     string    `json:"loan_status"`
	LoanFlowStatus string    `json:"loan_flow_status"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	TotalItems     int       `json:"total_items"` // jumlah barang berbeda
	TotalUnits     int       `json:"total_units"` // total units
	CreatedAt      time.Time `json:"created_at"`
}

// ========== DETAIL RESPONSE ==========
type LoanDetailResponse struct {
	LoanCode       string                   `json:"loan_code"`
	UserID         uint                     `json:"user_id"`
	LoanStatus     string                   `json:"loan_status"`
	LoanFlowStatus string                   `json:"loan_flow_status"`
	StartDate      time.Time                `json:"start_date"`
	EndDate        time.Time                `json:"end_date"`
	Reason         string                   `json:"reason"`
	
	TotalItems     int                      `json:"total_items"` // jumlah barang berbeda
	TotalUnits     int                      `json:"total_units"` // ✅ TAMBAH INI
	
	Items          []LoanItemDetailResponse `json:"items"`
	
	RevisionNote   string                   `json:"revision_note,omitempty"`
	RejectReason   *string                   `json:"reject_reason,omitempty"` // ✅ TAMBAH INI (opsional)
	
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"` // ✅ TAMBAH INI (opsional)
}

// ========== LOAN ITEM DETAIL ==========
type LoanItemDetailResponse struct {
	BarangID   uint         `json:"barang_id"`
	BarangKode string       `json:"barang_kode"`
	BarangNama string       `json:"barang_nama"`
	Merk       string       `json:"merk"`
	Kategori   string       `json:"kategori"`
	Quantity   int          `json:"quantity"`   // len(AssignedUnits)
	Units      []UnitDetail `json:"units"`      // Detail units
	Status     string       `json:"status"`
}

// ========== UNIT DETAIL ==========
type UnitDetail struct {
	UnitID   uint   `json:"unit_id"`
	KodeUnit string `json:"kode_unit"`
	Kondisi  string `json:"kondisi"`
	Status   string `json:"status"`
}

// ========== SIMPLE RESPONSE (untuk create/update) ==========
type LoanResponse struct {
	LoanCode       string              `json:"loan_code"`
	LoanStatus     string              `json:"loan_status"`
	LoanFlowStatus string              `json:"loan_flow_status"`
	StartDate      time.Time           `json:"start_date"`
	EndDate        time.Time           `json:"end_date"`
	Reason         string              `json:"reason"`
	TotalItems     int                 `json:"total_items"`
	TotalUnits     int                 `json:"total_units"` // ✅ TAMBAH INI
	Items          []LoanItemResponse  `json:"items"`
	CreatedAt      time.Time           `json:"created_at"`
}

type LoanItemResponse struct {
	BarangID   uint   `json:"barang_id"`
	BarangKode string `json:"barang_kode"`
	BarangNama string `json:"barang_nama"`
	Quantity   int    `json:"quantity"` // len(AssignedUnits)
	Status     string `json:"status"`
}


// AvailabilityCheckResponse - Response untuk check availability dengan gap rule & suggestions
type AvailabilityCheckResponse struct {
	IsAvailable    bool       `json:"is_available"`
	AvailableCount int        `json:"available_count"`
	RequestedCount int        `json:"requested_count"`
	Message        string     `json:"message"`
	SuggestedDate  *time.Time `json:"suggested_date,omitempty"` // 🔥 Tanggal saran (dengan gap 1 hari)
	
	// 🔥 Rekomendasi barang alternatif (kategori sama, merk beda)
	AlternativeBarangs []AlternativeBarang `json:"alternative_barangs,omitempty"`
}

// AlternativeBarang - Barang alternatif yang bisa dipinjam
type AlternativeBarang struct {
	BarangID       uint   `json:"barang_id"`
	BarangKode     string `json:"barang_kode"`
	BarangNama     string `json:"barang_nama"`
	Merk           string `json:"merk"`
	Kategori       string `json:"kategori"`
	AvailableCount int    `json:"available_count"` // Jumlah unit available di tanggal yang diminta
}