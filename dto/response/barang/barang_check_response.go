package barang

// Calendar Response
type AvailabilityCalendarResponse struct {
	BarangID   uint              `json:"barang_id"`
	BarangNama string            `json:"barang_nama"`
	BarangKode string            `json:"barang_kode"`
	TotalStock int               `json:"total_stock"`
	Month      string            `json:"month"`
	Days       []DayAvailability `json:"days"`
}

type DayAvailability struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"day_of_week"`
	Confirmed int    `json:"confirmed"`
	Pending   int    `json:"pending"`
	Revision  int    `json:"revision"`
	Available int    `json:"available"`
	Status    string `json:"status"`
}

// Date Detail Response
type DateDetailResponse struct {
	BarangID   uint          `json:"barang_id"`
	BarangNama string        `json:"barang_nama"`
	BarangKode string        `json:"barang_kode"`
	TotalStock int           `json:"total_stock"`
	Date       string        `json:"date"`
	Confirmed  int           `json:"confirmed"`
	Pending    int           `json:"pending"`
	Revision   int           `json:"revision"`
	Available  int           `json:"available"`
	Loans      []LoanDetail  `json:"loans"`
	Summary    DetailSummary `json:"summary"`
}

type LoanDetail struct {
	LoanCode       string `json:"loan_code"`
	UserName       string `json:"user_name"`
	Quantity       int    `json:"quantity"`
	LoanStatus     string `json:"loan_status"`
	LoanFlowStatus string `json:"loan_flow_status"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	IsRevised      bool   `json:"is_revised"`
	RevisionReason string `json:"revision_reason,omitempty"`
	Reason         string `json:"reason"`
	CreatedAt      string `json:"created_at"`
	Priority       string `json:"priority"`
}

type DetailSummary struct {
	TotalRequests      int    `json:"total_requests"`
	RevisionCount      int    `json:"revision_count"`
	PendingCount       int    `json:"pending_count"`
	CanAccommodateAll  bool   `json:"can_accommodate_all"`
	Notes              string `json:"notes"`
}