// dto/response/barang/availability.go
package barang

type UnitDetail struct {
	UnitID   uint   `json:"unit_id"`
	KodeUnit string `json:"kode_unit"`
	Kondisi  string `json:"kondisi"`
	Status   string `json:"status"`
}

type LoanDetail struct {
	LoanCode       string       `json:"loan_code"`
	UserName       string       `json:"user_name"`
	Quantity       int          `json:"quantity"` // jumlah units
	Units          []UnitDetail `json:"units"`    // ← TAMBAH INI
	LoanStatus     string       `json:"loan_status"`
	LoanFlowStatus string       `json:"loan_flow_status"`
	StartDate      string       `json:"start_date"`
	EndDate        string       `json:"end_date"`
	IsRevised      bool         `json:"is_revised"`
	RevisionReason string       `json:"revision_reason,omitempty"`
	Reason         string       `json:"reason"`
	CreatedAt      string       `json:"created_at"`
	Priority       string       `json:"priority"`
}

type DayAvailability struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"day_of_week"`
	Confirmed int    `json:"confirmed"` // units confirmed
	Pending   int    `json:"pending"`   // units pending
	Revision  int    `json:"revision"`  // units revision
	Available int    `json:"available"` // units available
	Status    string `json:"status"`
}

type AvailabilityCalendarResponse struct {
	BarangID   uint              `json:"barang_id"`
	BarangNama string            `json:"barang_nama"`
	BarangKode string            `json:"barang_kode"`
	TotalStock int               `json:"total_stock"` // total units
	Month      string            `json:"month"`
	Days       []DayAvailability `json:"days"`
}

type DetailSummary struct {
	TotalRequests     int    `json:"total_requests"`
	RevisionCount     int    `json:"revision_count"`
	PendingCount      int    `json:"pending_count"`
	CanAccommodateAll bool   `json:"can_accommodate_all"`
	Notes             string `json:"notes"`
}

type DateDetailResponse struct {
	BarangID   uint          `json:"barang_id"`
	BarangNama string        `json:"barang_nama"`
	BarangKode string        `json:"barang_kode"`
	TotalStock int           `json:"total_stock"` // total units
	Date       string        `json:"date"`
	Confirmed  int           `json:"confirmed"`
	Pending    int           `json:"pending"`
	Revision   int           `json:"revision"`
	Available  int           `json:"available"`
	Loans      []LoanDetail  `json:"loans"`
	Summary    DetailSummary `json:"summary"`
}