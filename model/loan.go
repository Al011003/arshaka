package model

import (
	"time"

	"gorm.io/gorm"
)

//
// ==========================
// Loan Decision Status (ADMIN)
// ==========================
//
const (
	LoanStatusPending  = "PENDING"
	LoanStatusApproved = "APPROVED"
	LoanStatusRejected = "REJECTED"
)

var LoanStatusEnum = []string{
	LoanStatusPending,
	LoanStatusApproved,
	LoanStatusRejected,
}

//
// ==========================
// Loan Flow Status (OPERATIONAL)
// ==========================
//
const (
	LoanFlowRequested = "REQUESTED" // user submit
	LoanFlowRevision  = "REVISION"  // admin minta revisi
	LoanFlowReady     = "READY"     // units di-reserve
	LoanFlowTaken     = "TAKEN"     // barang diambil
	LoanFlowFinished  = "FINISHED"  // dikembalikan
)

var LoanFlowEnum = []string{
	LoanFlowRequested,
	LoanFlowRevision,
	LoanFlowReady,
	LoanFlowTaken,
	LoanFlowFinished,
}

//
// ==========================
// Loan Item Status
// ==========================
//
const (
	LoanItemStatusRequested = "REQUESTED"
	LoanItemStatusReady     = "READY"
	LoanItemStatusBorrowed  = "BORROWED"
	LoanItemStatusReturned  = "RETURNED"
	LoanItemStatusRejected  = "REJECTED"
)

var LoanItemStatusEnum = []string{
	LoanItemStatusRequested,
	LoanItemStatusReady,
	LoanItemStatusBorrowed,
	LoanItemStatusReturned,
	LoanItemStatusRejected,
}

//
// ==========================
// History Actor Role
// ==========================
//
const (
	HistoryRoleAdmin  = "ADMIN"
	HistoryRoleUser   = "USER"
	HistoryRoleSystem = "SYSTEM"
)

var HistoryRoleEnum = []string{
	HistoryRoleAdmin,
	HistoryRoleUser,
	HistoryRoleSystem,
}

//
// ==========================
// Loan
// ==========================
//
type Loan struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null;index"`

	// Identifier
	LoanCode string `json:"loan_code" gorm:"type:varchar(30);uniqueIndex"`

	// Decision (ADMIN)
	LoanStatus string `json:"loan_status" gorm:"type:varchar(20);not null;default:'PENDING'"`

	RejectReason *string `json:"reject_reason" gorm:"type:text"`

	// Operational Flow
	LoanFlowStatus string `json:"loan_flow_status" gorm:"type:varchar(20);not null;default:'REQUESTED'"`

	// Date Range
	StartDate time.Time `json:"start_date" gorm:"not null;index"`
	EndDate   time.Time `json:"end_date" gorm:"not null;index"`

	// Metadata
	Reason string `json:"reason" gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LoanItems []LoanItem `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE" json:"loan_items,omitempty"`
}

func (Loan) TableName() string {
	return "loans"
}

//
// ==========================
// Loan Item
// ==========================
//
type LoanItem struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	LoanID   uint `json:"loan_id" gorm:"not null;index"`
	BarangID uint `json:"barang_id" gorm:"not null;index"`

	Status string `json:"status" gorm:"type:varchar(20);not null;default:'REQUESTED'"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations
	Loan          *Loan        `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
	Barang        *Barang      `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
	AssignedUnits []BarangUnit `gorm:"many2many:loan_item_units;joinForeignKey:LoanItemID;joinReferences:UnitID" json:"assigned_units,omitempty"`
}

func (LoanItem) TableName() string {
	return "loan_items"
}

//
// ==========================
// Junction Table (LoanItem <-> Unit)
// ==========================
//
type LoanItemUnit struct {
	LoanItemID uint      `gorm:"primaryKey;autoIncrement:false;index:idx_loan_unit,unique,priority:1"`
	UnitID     uint      `gorm:"primaryKey;autoIncrement:false;index:idx_loan_unit,unique,priority:2"`
	CreatedAt  time.Time `json:"created_at"`
}

func (LoanItemUnit) TableName() string {
	return "loan_item_units"
}

//
// ==========================
// Loan Status History (DECISION)
// ==========================
//
type LoanStatusHistory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`

	FromStatus *string `json:"from_status"`
	ToStatus   string  `json:"to_status"`

	ChangedBy     uint   `json:"changed_by"`
	ChangedByRole string `json:"changed_by_role" gorm:"type:varchar(10);not null"`

	Note *string `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time

	Loan *Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

func (LoanStatusHistory) TableName() string {
	return "loan_status_history"
}

//
// ==========================
// Loan Flow History (OPERATIONAL)
// ==========================
//
type LoanFlowHistory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	LoanID uint `json:"loan_id" gorm:"not null;index"`

	FromFlow *string `json:"from_flow" gorm:"type:varchar(20)"`
	ToFlow   string  `json:"to_flow" gorm:"type:varchar(20);not null"`

	TriggeredBy     uint   `json:"triggered_by"`
	TriggeredByRole string `json:"triggered_by_role" gorm:"type:varchar(10);not null"`

	Note string `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time

	Loan *Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

func (LoanFlowHistory) TableName() string {
	return "loan_flow_history"
}
