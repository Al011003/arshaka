package loan

import "time"

type LoanResponse struct {
	ID             uint                `json:"id"`
	UserID         uint                `json:"user_id"`
	StartDate      time.Time           `json:"start_date"`
	EndDate        time.Time           `json:"end_date"`
	Reason         string              `json:"reason"`
	LoanStatus     string              `json:"loan_status"`      // PENDING, APPROVED, REJECTED, etc
	LoanFlowStatus string              `json:"loan_flow_status"` // REQUESTED, APPROVED, BORROWED, RETURNED, etc
	TotalItems     int                 `json:"total_items"`
	LoanItems      []LoanItemResponse  `json:"loan_items"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	LoanCode string `json:"loan_code"`
}

type LoanItemResponse struct {
	ID       uint                   `json:"id"`
	LoanID   uint                   `json:"loan_id"`
	BarangID uint                   `json:"barang_id"`
	Barang   *BarangInLoanResponse  `json:"barang"`
	Quantity int                    `json:"quantity"`
	Status   string                 `json:"status"` // BORROWED, RETURNED, etc
}

type BarangInLoanResponse struct {
	ID       uint   `json:"id"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	Merk     string `json:"merk"`
	Kategori string `json:"kategori"`
	CoverURL string `json:"cover_url"`
}
