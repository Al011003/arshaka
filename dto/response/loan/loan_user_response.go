package loan

import "time"

// Untuk CARD / TABLE
type LoanListResponse struct {
	LoanCode       string    `json:"loan_code"`
	LoanStatus     string    `json:"loan_status"`
	LoanFlowStatus string    `json:"loan_flow_status"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	TotalItems     int       `json:"total_items"`
	CreatedAt      time.Time `json:"created_at"`
}

type LoanDetailResponse struct {
	LoanCode       string                   `json:"loan_code"`
	UserID         uint                     `json:"user_id"`
	LoanStatus     string                   `json:"loan_status"`
	LoanFlowStatus string                   `json:"loan_flow_status"`
	StartDate      time.Time                `json:"start_date"`
	EndDate        time.Time                `json:"end_date"`
	Reason         string                   `json:"reason"`
	TotalItems     int                      `json:"total_items"`
	Items          []LoanItemDetailResponse `json:"items"`
	CreatedAt      time.Time                `json:"created_at"`
}

type LoanItemDetailResponse struct {
	BarangID uint   `json:"barang_id"`
	Nama     string `json:"nama"`
	Merk     string `json:"merk"`
	Kategori string `json:"kategori"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

