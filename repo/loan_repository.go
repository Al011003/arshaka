// repository/loan_repository.go
package repo

import (
	request "backend/dto/request/loan"
	"backend/model"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ========== STRUCT HELPER ==========

type LoanReturnInfo struct {
	LoanID    uint
	EndDate   time.Time
	BarangID  uint
	UnitCount int
	UnitIDs   []uint
}

// ========== INTERFACE ==========

type LoanRepository interface {
	// ===== LOAN (HEADER) =====
	CreateLoan(loan *model.Loan) error
	UpdateLoan(loan *model.Loan) error
	GetLoanByID(loanID uint) (*model.Loan, error)
	GetLoanByIDWithItems(loanID uint) (*model.Loan, error)
	GetLoanByUserID(userID uint) ([]model.Loan, error)
	GetLoanByCode(code string) (*model.Loan, error)
	GetLoanByCodeWithItems(code string) (*model.Loan, error)
	GetLoanForUpdate(loanID uint) (*model.Loan, error)
	SoftDelete(id uint) error

	// ===== LOAN ITEM =====
	CreateLoanItem(item *model.LoanItem) error
	CreateLoanItems(items []model.LoanItem) error
	GetLoanItemsByLoanID(loanID uint) ([]model.LoanItem, error)
	UpdateLoanItem(item *model.LoanItem) error
	DeleteLoanItem(loanID uint, barangID uint) error
	DeleteLoanItemsByLoanID(loanID uint) error

	// ===== UNIT ASSIGNMENT =====
	AssignUnitsToLoanItem(loanItemID uint, unitIDs []uint) error
	GetAssignedUnits(loanItemID uint) ([]model.BarangUnit, error)
	RemoveUnitsFromLoanItem(loanItemID uint) error

	// ===== AVAILABILITY & OVERLAPPING =====
	GetLoanItemsInDateRange(barangID uint, startDate, endDate time.Time) ([]model.LoanItem, error)
	GetOverlappingUnits(barangID uint, startDate, endDate time.Time, excludeLoanID *uint) ([]uint, error)

	// 🔥 NEW: Gap 1 Hari Support
	GetUpcomingReturns(barangID uint, afterDate time.Time) ([]LoanReturnInfo, error)

	// 🔥 NEW: Auto-Revision Support
	GetActiveLoansByBarangID(barangID uint) ([]model.Loan, error)
	BulkUpdateLoanFlow(loanIDs []uint, newFlow string) error

	// ===== STATUS UPDATE =====
	UpdateLoanStatus(loanID uint, newStatus string) error
	UpdateLoanFlow(loanID uint, newFlow string) error

	// ===== TRANSACTION =====
	CreateLoanFromCartTx(loan *model.Loan, items []model.LoanItem, cartID uint) error

	// ===== HISTORY =====
	CreateStatusHistory(history *model.LoanStatusHistory) error
	CreateFlowHistory(history *model.LoanFlowHistory) error
	GetStatusHistory(loanID uint) ([]model.LoanStatusHistory, error)
	GetFlowHistory(loanID uint) ([]model.LoanFlowHistory, error)

	 GetAllLoansWithFilter(filter request.AdminLoanFilter) ([]model.Loan, error)
}

// ========== IMPLEMENTATION ==========

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

// ========== LOAN METHODS ==========

func (r *loanRepository) CreateLoan(loan *model.Loan) error {
	return r.db.Create(loan).Error
}

func (r *loanRepository) UpdateLoan(loan *model.Loan) error {
	return r.db.Save(loan).Error
}

func (r *loanRepository) GetLoanByID(loanID uint) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.First(&loan, loanID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan tidak ditemukan")
		}
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetLoanByIDWithItems(loanID uint) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.
		Preload("User").
		Preload("LoanItems.Barang").
		Preload("LoanItems.AssignedUnits").
		Preload("LoanItems.AssignedUnits.Komponen").
		First(&loan, loanID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan tidak ditemukan")
		}
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetLoanByUserID(userID uint) ([]model.Loan, error) {
	var loans []model.Loan
	err := r.db.
		Preload("LoanItems.AssignedUnits"). // 🔥 Tambah ini buat hitung units
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&loans).Error

	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *loanRepository) GetLoanByCode(code string) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.Where("loan_code = ?", code).First(&loan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan tidak ditemukan")
		}
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetLoanByCodeWithItems(code string) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.
		Preload("User").
		Preload("LoanItems.Barang").
		Preload("LoanItems.AssignedUnits").
		Preload("LoanItems.AssignedUnits.Komponen").
		Where("loan_code = ?", code).
		First(&loan).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan tidak ditemukan")
		}
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetLoanForUpdate(loanID uint) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("LoanItems").
		First(&loan, loanID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan tidak ditemukan")
		}
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) SoftDelete(id uint) error {
	return r.db.Delete(&model.Loan{}, id).Error
}

// ========== LOAN ITEM METHODS ==========

func (r *loanRepository) CreateLoanItem(item *model.LoanItem) error {
	return r.db.Create(item).Error
}

func (r *loanRepository) CreateLoanItems(items []model.LoanItem) error {
	if len(items) == 0 {
		return errors.New("loan items kosong")
	}
	return r.db.Create(&items).Error
}

func (r *loanRepository) GetLoanItemsByLoanID(loanID uint) ([]model.LoanItem, error) {
	var items []model.LoanItem
	err := r.db.
		Preload("Barang").
		Preload("AssignedUnits").
		Where("loan_id = ?", loanID).
		Find(&items).Error

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *loanRepository) UpdateLoanItem(item *model.LoanItem) error {
	return r.db.Save(item).Error
}

func (r *loanRepository) DeleteLoanItem(loanID uint, barangID uint) error {
	return r.db.
		Where("loan_id = ? AND barang_id = ?", loanID, barangID).
		Delete(&model.LoanItem{}).
		Error
}

func (r *loanRepository) DeleteLoanItemsByLoanID(loanID uint) error {
	return r.db.
		Where("loan_id = ?", loanID).
		Delete(&model.LoanItem{}).Error
}

// ========== UNIT ASSIGNMENT METHODS ==========

func (r *loanRepository) AssignUnitsToLoanItem(loanItemID uint, unitIDs []uint) error {
	var loanItem model.LoanItem

	if err := r.db.First(&loanItem, loanItemID).Error; err != nil {
		return err
	}

	var units []model.BarangUnit
	if err := r.db.Where("id IN ?", unitIDs).Find(&units).Error; err != nil {
		return err
	}

	// Replace associations (clear old + add new)
	if err := r.db.Model(&loanItem).Association("AssignedUnits").Replace(&units); err != nil {
		return err
	}

	return nil
}

func (r *loanRepository) GetAssignedUnits(loanItemID uint) ([]model.BarangUnit, error) {
	var loanItem model.LoanItem

	err := r.db.
		Preload("AssignedUnits").
		Preload("AssignedUnits.Komponen").
		First(&loanItem, loanItemID).Error

	if err != nil {
		return nil, err
	}

	return loanItem.AssignedUnits, nil
}

func (r *loanRepository) RemoveUnitsFromLoanItem(loanItemID uint) error {
	var loanItem model.LoanItem

	if err := r.db.First(&loanItem, loanItemID).Error; err != nil {
		return err
	}

	if err := r.db.Model(&loanItem).Association("AssignedUnits").Clear(); err != nil {
		return err
	}

	return nil
}

// ========== AVAILABILITY & OVERLAPPING ==========

func (r *loanRepository) GetLoanItemsInDateRange(
	barangID uint,
	startDate, endDate time.Time,
) ([]model.LoanItem, error) {
	var loanItems []model.LoanItem

	err := r.db.
		Preload("Loan").
		Preload("Loan.User").
		Preload("AssignedUnits").
		Preload("AssignedUnits.Komponen").
		Joins("JOIN loans ON loans.id = loan_items.loan_id").
		Where("loan_items.barang_id = ?", barangID).
		Where("loan_items.status NOT IN ?", []string{"REJECTED", "RETURNED"}).
		Where("loans.loan_status != ?", "REJECTED").
		Where("loans.start_date <= ? AND loans.end_date >= ?", endDate, startDate).
		Order("loans.loan_status DESC, loans.created_at ASC").
		Find(&loanItems).Error

	if err != nil {
		return nil, err
	}

	return loanItems, nil
}

func (r *loanRepository) GetOverlappingUnits(
	barangID uint,
	startDate time.Time,
	endDate time.Time,
	excludeLoanID *uint,
) ([]uint, error) {
	var unitIDs []uint

	query := r.db.Table("loan_item_units liu").
		Select("DISTINCT liu.unit_id").
		Joins("JOIN loan_items li ON li.id = liu.loan_item_id").
		Joins("JOIN loans l ON l.id = li.loan_id").
		Where("li.barang_id = ?", barangID).
		Where("l.start_date <= ? AND l.end_date >= ?", endDate, startDate).
		Where("l.loan_status = ?", "APPROVED").
		Where("l.loan_flow_status IN ?", []string{"READY", "TAKEN"})

	if excludeLoanID != nil {
		query = query.Where("l.id != ?", *excludeLoanID)
	}

	err := query.Pluck("liu.unit_id", &unitIDs).Error

	return unitIDs, err
}

// ========== 🔥 NEW: GAP 1 HARI SUPPORT ==========

func (r *loanRepository) GetUpcomingReturns(
	barangID uint,
	afterDate time.Time,
) ([]LoanReturnInfo, error) {
	
	type queryResult struct {
		LoanID   uint
		EndDate  time.Time
		BarangID uint
		UnitID   uint
	}

	var results []queryResult

	// Query: ambil semua loan yang end_date >= afterDate
	err := r.db.Table("loan_item_units liu").
		Select("l.id as loan_id, l.end_date, li.barang_id, liu.unit_id").
		Joins("JOIN loan_items li ON li.id = liu.loan_item_id").
		Joins("JOIN loans l ON l.id = li.loan_id").
		Where("li.barang_id = ?", barangID).
		Where("l.end_date >= ?", afterDate).
		Where("l.loan_status = ?", "APPROVED").
		Where("l.loan_flow_status IN ?", []string{"READY", "TAKEN"}).
		Order("l.end_date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Group by LoanID
	grouped := make(map[uint]*LoanReturnInfo)

	for _, row := range results {
		if _, exists := grouped[row.LoanID]; !exists {
			grouped[row.LoanID] = &LoanReturnInfo{
				LoanID:   row.LoanID,
				EndDate:  row.EndDate,
				BarangID: row.BarangID,
				UnitIDs:  []uint{},
			}
		}
		grouped[row.LoanID].UnitIDs = append(grouped[row.LoanID].UnitIDs, row.UnitID)
	}

	// Convert map to slice
	returnInfos := make([]LoanReturnInfo, 0, len(grouped))
	for _, info := range grouped {
		info.UnitCount = len(info.UnitIDs)
		returnInfos = append(returnInfos, *info)
	}

	// Sort by EndDate
	// (Go doesn't have built-in sort for custom structs, tapi bisa pake sort.Slice)
	// Untuk simplicity, udah di-order di query

	return returnInfos, nil
}

// ========== 🔥 NEW: AUTO-REVISION SUPPORT ==========

func (r *loanRepository) GetActiveLoansByBarangID(barangID uint) ([]model.Loan, error) {
	var loans []model.Loan

	err := r.db.
		Distinct("loans.*").
		Joins("JOIN loan_items ON loan_items.loan_id = loans.id").
		Where("loan_items.barang_id = ?", barangID).
		Where("loans.loan_status = ?", "APPROVED").
		Where("loans.loan_flow_status IN ?", []string{"REQUESTED", "READY", "TAKEN"}).
		Find(&loans).Error

	return loans, err
}

func (r *loanRepository) BulkUpdateLoanFlow(loanIDs []uint, newFlow string) error {
	if len(loanIDs) == 0 {
		return nil
	}

	return r.db.
		Model(&model.Loan{}).
		Where("id IN ?", loanIDs).
		Update("loan_flow_status", newFlow).
		Error
}

// ========== STATUS UPDATE ==========

func (r *loanRepository) UpdateLoanStatus(loanID uint, newStatus string) error {
	return r.db.
		Model(&model.Loan{}).
		Where("id = ?", loanID).
		Update("loan_status", newStatus).
		Error
}

func (r *loanRepository) UpdateLoanFlow(loanID uint, newFlow string) error {
	return r.db.
		Model(&model.Loan{}).
		Where("id = ?", loanID).
		Update("loan_flow_status", newFlow).
		Error
}

// ========== TRANSACTION ==========

func (r *loanRepository) CreateLoanFromCartTx(
	loan *model.Loan,
	items []model.LoanItem,
	cartID uint,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(loan).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].LoanID = loan.ID
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// ========== HISTORY METHODS ==========

func (r *loanRepository) CreateStatusHistory(history *model.LoanStatusHistory) error {
	return r.db.Create(history).Error
}

func (r *loanRepository) CreateFlowHistory(history *model.LoanFlowHistory) error {
	return r.db.Create(history).Error
}

func (r *loanRepository) GetStatusHistory(loanID uint) ([]model.LoanStatusHistory, error) {
	var history []model.LoanStatusHistory
	err := r.db.
		Where("loan_id = ?", loanID).
		Order("created_at DESC").
		Find(&history).Error
	return history, err
}

func (r *loanRepository) GetFlowHistory(loanID uint) ([]model.LoanFlowHistory, error) {
	var history []model.LoanFlowHistory
	err := r.db.
		Where("loan_id = ?", loanID).
		Order("created_at DESC").
		Find(&history).Error
	return history, err
}

func (r *loanRepository) GetAllLoansWithFilter(
    filter request.AdminLoanFilter,
) ([]model.Loan, error) {

    var loans []model.Loan

    q := r.db.
        Preload("User").
        Preload("LoanItems.Barang").
        Preload("LoanItems.AssignedUnits")

    if filter.LoanStatus != "" {
        q = q.Where("loan_status = ?", filter.LoanStatus)
    }

    if filter.LoanFlowStatus != "" {
        q = q.Where("loan_flow_status = ?", filter.LoanFlowStatus)
    }

    if filter.UserID != 0 {
        q = q.Where("user_id = ?", filter.UserID)
    }

    err := q.Order("created_at DESC").Find(&loans).Error
    return loans, err
}
