// model/loan.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null;index"`

	// Identifier
	LoanCode string `json:"loan_code" gorm:"type:varchar(30);uniqueIndex"`

	// Decision Status (ADMIN)
	// PENDING  -> belum diputuskan / revisi
	// APPROVED -> disetujui admin
	// REJECTED -> ditolak final
	LoanStatus string `json:"loan_status" gorm:"type:enum('PENDING','APPROVED','REJECTED');default:'PENDING'"`

	// Alasan reject (WAJIB kalau REJECTED)
	RejectReason *string `json:"reject_reason" gorm:"type:text"`

	// Flow Status (OPERASIONAL)
	// REQUESTED -> user submit
	// REVISION  -> perlu perbaikan user
	// READY     -> units sudah di-reserve
	// TAKEN     -> barang diambil
	// FINISHED  -> dikembalikan
	LoanFlowStatus string `json:"loan_flow_status" gorm:"type:enum('REQUESTED','REVISION','READY','TAKEN','FINISHED');default:'REQUESTED'"`

	// Date Range
	StartDate time.Time `json:"start_date" gorm:"not null;index"`
	EndDate   time.Time `json:"end_date" gorm:"not null;index"`

	// Metadata
	Reason string `json:"reason" gorm:"type:text"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"` // ✅ TAMBAH INI
	LoanItems []LoanItem `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE" json:"loan_items,omitempty"`
}

func (Loan) TableName() string {
	return "loans"
}

type LoanItem struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	LoanID   uint `json:"loan_id" gorm:"not null;index"`
	BarangID uint `json:"barang_id" gorm:"not null;index"`

	// ❌ REMOVE: Quantity (diganti jadi count AssignedUnits)
	// Quantity int `json:"quantity" gorm:"not null"`

	// Status per item
	// REQUESTED -> ikut loan
	// READY     -> units di-reserve
	// BORROWED  -> sudah diambil
	// RETURNED  -> dikembalikan
	// REJECTED  -> item ditolak (parsial)
	Status string `json:"status" gorm:"type:enum('REQUESTED','READY','BORROWED','RETURNED','REJECTED');default:'REQUESTED'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Loan          *Loan        `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
	Barang        *Barang      `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
	AssignedUnits []BarangUnit `gorm:"many2many:loan_item_units;" json:"assigned_units,omitempty"` // ✅ RENAME dari Units
}

func (LoanItem) TableName() string {
	return "loan_items"
}

// ✅ Junction table untuk many-to-many
type LoanItemUnit struct {
	LoanItemID uint      `gorm:"primaryKey"`
	UnitID     uint      `gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
}

func (LoanItemUnit) TableName() string {
	return "loan_item_units"
}

// ✅ HISTORY MODELS (OPTIONAL tapi RECOMMENDED)
type LoanStatusHistory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`

	FromStatus *string `json:"from_status"`
	ToStatus   string  `json:"to_status"`

	ChangedBy     uint   `json:"changed_by"`
	ChangedByRole string `json:"changed_by_role" gorm:"type:enum('ADMIN','USER','SYSTEM')"`

	Note *string `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`

	Loan *Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

func (LoanStatusHistory) TableName() string {
	return "loan_status_history"
}

type LoanFlowHistory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`

	FromFlow *string `json:"from_flow"`
	ToFlow   string  `json:"to_flow" gorm:"not null"`

	ChangedBy     uint   `json:"changed_by"`
	ChangedByRole string `json:"changed_by_role" gorm:"type:enum('ADMIN','USER','SYSTEM')"`

	Note *string `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`

	Loan *Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

func (LoanFlowHistory) TableName() string {
	return "loan_flow_history"
}