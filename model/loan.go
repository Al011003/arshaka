package model

import (
	"time"

	"gorm.io/gorm"
)

/*
=========================
LOAN (HEADER)
=========================
*/

type Loan struct {
	gorm.Model
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null;index"`

	// =========================
	// IDENTIFIER
	// =========================
	LoanCode string `json:"loan_code" gorm:"type:varchar(30);uniqueIndex"`

	// =========================
	// DECISION STATUS (ADMIN)
	// =========================
	// PENDING  -> belum diputuskan / revisi
	// APPROVED -> disetujui admin
	// REJECTED -> ditolak final
	LoanStatus string `json:"loan_status" gorm:"type:enum('PENDING','APPROVED','REJECTED');default:'PENDING'"`

	// Alasan reject (WAJIB kalau REJECTED)
	RejectReason *string `json:"reject_reason" gorm:"type:text"`

	// =========================
	// FLOW STATUS (OPERASIONAL)
	// =========================
	// REQUESTED -> user submit
	// REVISION  -> perlu perbaikan user
	// READY     -> stok sudah di-reserve
	// TAKEN     -> barang diambil
	// FINISHED  -> dikembalikan
	LoanFlowStatus string `json:"loan_flow_status" gorm:"type:enum('REQUESTED','REVISION','READY','TAKEN','FINISHED');default:'REQUESTED'"`

	// =========================
	// DATE
	// =========================
	StartDate time.Time `json:"start_date" gorm:"not null;index"`
	EndDate   time.Time `json:"end_date" gorm:"not null;index"`

	// =========================
	// META
	// =========================
	Reason string `json:"reason" gorm:"type:text"`

	LoanItems []LoanItem `json:"loan_items,omitempty" gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}


type LoanItem struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`
	Loan   *Loan `json:"loan,omitempty" gorm:"foreignKey:LoanID"`

	BarangID uint    `json:"barang_id" gorm:"not null;index"`
	Barang   *Barang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`

	Quantity int `json:"quantity" gorm:"not null"`

	// REQUESTED -> ikut loan
	// READY     -> stok di-reserve
	// BORROWED  -> sudah diambil
	// RETURNED  -> dikembalikan
	// REJECTED  -> item ditolak (parsial)
	Status string `json:"status" gorm:"type:enum('REQUESTED','READY','BORROWED','RETURNED','REJECTED');default:'REQUESTED'"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

 
type LoanStatusHistory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`
	Loan   *Loan `json:"loan,omitempty" gorm:"foreignKey:LoanID"`

	FromStatus *string `json:"from_status"`
	ToStatus   string  `json:"to_status"`

	ChangedBy     uint   `json:"changed_by"`
	ChangedByRole string `json:"changed_by_role" gorm:"type:enum('ADMIN','SYSTEM')"`

	Note *string `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}


type LoanFlowHistory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`
	Loan   *Loan `json:"loan,omitempty" gorm:"foreignKey:LoanID"`

	FromFlow *string `json:"from_flow" gorm:"type:enum('REQUESTED','REVISION','READY','TAKEN','FINISHED')"`
	ToFlow   string  `json:"to_flow" gorm:"type:enum('REQUESTED','REVISION','READY','TAKEN','FINISHED');not null"`

	ChangedBy     uint   `json:"changed_by"`
	ChangedByRole string `json:"changed_by_role" gorm:"type:enum('ADMIN','SYSTEM')"`

	Note *string `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
