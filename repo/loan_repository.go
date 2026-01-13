package repo

import (
	"backend/model"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
=================================================
INTERFACE
=================================================
*/

type LoanRepository interface {

	// ===== LOAN (HEADER) =====
	CreateLoan(loan *model.Loan) error
	UpdateLoan(loan *model.Loan) error
	GetLoanByID(loanID uint) (*model.Loan, error)
	GetLoanByIDWithItems(loanID uint) (*model.Loan, error)
	GetLoanByUserID(userID uint) ([]model.Loan, error)

	// ===== LOAN ITEM =====
	CreateLoanItem(item *model.LoanItem) error
	CreateLoanItems(items []model.LoanItem) error
	GetLoanItemsByLoanID(loanID uint) ([]model.LoanItem, error)
	DeleteLoanItemsByLoanID(loanID uint) error
	GetTotalBorrowedQtyByDate(barangID uint, startDate time.Time, endDate time.Time) (int64, error)

	// ===== TRANSACTION =====
	CreateLoanFromCartTx(loan *model.Loan, items []model.LoanItem, cartID uint) error
	GetLoanByCode(code string) (*model.Loan, error)
	GetLoanByCodeWithItems(code string) (*model.Loan, error)
	GetLoanItemsInDateRange(barangID uint, startDate, endDate time.Time) ([]model.LoanItem, error)

	// ===== STATE UPDATE =====
	UpdateLoanStatus(loanID uint, newStatus string) error
	UpdateLoanFlow(loanID uint, newFlow string) error

	// ===== LOCKING =====
	GetLoanForUpdate(loanID uint) (*model.Loan, error)

	// ====== BARANG ======
	UpdateLoanItem(item *model.LoanItem) error
	DeleteLoanItem(loanID uint, barangID uint) error
	FindByIDWithItems(loanID uint) (*model.Loan, error)
	FindByIDAndUserIDWithItems(loanID uint, userID uint) (*model.Loan, error)
	SoftDelete(id uint) error
	GetOverlappingLoanItems(barangID uint, startDate time.Time, endDate time.Time, excludeLoanID *uint,) (int, error)

}


type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

/*
=================================================
LOAN METHODS
=================================================
*/

// CreateLoan - buat loan baru (PENDING)
func (r *loanRepository) CreateLoan(loan *model.Loan) error {
	return r.db.Create(loan).Error
}

// UpdateLoan - update loan (HANYA PENDING / ADMIN ACTION)
func (r *loanRepository) UpdateLoan(loan *model.Loan) error {
	return r.db.Save(loan).Error
}

// GetLoanByID - ambil loan tanpa item
func (r *loanRepository) GetLoanByID(loanID uint) (*model.Loan, error) {
	var loan model.Loan

	err := r.db.First(&loan, loanID).Error
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

// GetLoanByIDWithItems - ambil loan + items + barang
func (r *loanRepository) GetLoanByIDWithItems(loanID uint) (*model.Loan, error) {
	var loan model.Loan

	err := r.db.
		Preload("LoanItems.Barang").
		Where("id = ?", loanID).
		First(&loan).Error

	if err != nil {
		return nil, err
	}

	return &loan, nil
}

// GetLoanByUserID - ambil semua loan milik user
func (r *loanRepository) GetLoanByUserID(userID uint) ([]model.Loan, error) {
	var loans []model.Loan

	err := r.db.
		Preload("LoanItems").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&loans).Error

	if err != nil {
		return nil, err
	}

	return loans, nil
}


/*
=================================================
LOAN ITEM METHODS
=================================================
*/

// CreateLoanItem - create satu loan item
func (r *loanRepository) CreateLoanItem(item *model.LoanItem) error {
	return r.db.Create(item).Error
}

// CreateLoanItems - bulk insert loan items
func (r *loanRepository) CreateLoanItems(items []model.LoanItem) error {
	if len(items) == 0 {
		return errors.New("loan items kosong")
	}
	return r.db.Create(&items).Error
}

// GetLoanItemsByLoanID - ambil semua item dari loan
func (r *loanRepository) GetLoanItemsByLoanID(loanID uint) ([]model.LoanItem, error) {
	var items []model.LoanItem

	err := r.db.
		Preload("Barang").
		Where("loan_id = ?", loanID).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteLoanItemsByLoanID - hapus semua item loan (dipakai pas update PENDING)
func (r *loanRepository) DeleteLoanItemsByLoanID(loanID uint) error {
	return r.db.
		Where("loan_id = ?", loanID).
		Delete(&model.LoanItem{}).Error
}

// GetTotalBorrowedQtyByDate - hitung total qty yang dipinjam di range tanggal
func (r *loanRepository) GetTotalBorrowedQtyByDate(
	barangID uint,
	startDate time.Time,
	endDate time.Time,
) (int64, error) {

	var total int64

	err := r.db.
		Table("loan_items").
		Joins("JOIN loans ON loans.id = loan_items.loan_id").
		Where("loan_items.barang_id = ?", barangID).
		Where("loan_items.status IN ?", []string{"READY", "BORROWED"}).
		Where(`
			NOT (
				loans.end_date < ? OR loans.start_date > ?
			)
		`, startDate, endDate).
		Select("COALESCE(SUM(loan_items.quantity), 0)").
		Scan(&total).Error

	return total, err
}


/*
=================================================
TRANSACTION METHODS
=================================================
*/

// CreateLoanFromCartTx - Transaction: Create loan + items + clear cart
func (r *loanRepository) CreateLoanFromCartTx(loan *model.Loan, items []model.LoanItem, cartID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create loan
		if err := tx.Create(loan).Error; err != nil {
			return err
		}

		// 2. Set loan_id untuk semua items
		for i := range items {
			items[i].LoanID = loan.ID
		}

		// 3. Create loan items (bulk insert)
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		// 4. Clear cart items
		if err := tx.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *loanRepository) GetLoanByCode(code string) (*model.Loan, error) {
	var loan model.Loan

	err := r.db.
		Where("loan_code = ?", code).
		First(&loan).Error

	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) GetLoanByCodeWithItems(code string) (*model.Loan, error) {
	var loan model.Loan

	err := r.db.
		Preload("LoanItems.Barang").
		Where("loan_code = ?", code).
		First(&loan).Error

	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) GetLoanItemsInDateRange(barangID uint, startDate, endDate time.Time) ([]model.LoanItem, error) {
	var loanItems []model.LoanItem
	err := r.db.
		Preload("Loan").
		Joins("JOIN loans ON loans.id = loan_items.loan_id").
		Where("loan_items.barang_id = ?", barangID).
		Where("loan_items.status NOT IN ?", []string{"REJECTED", "RETURNED"}).
		Where("loans.loan_status != ?", "REJECTED").
		Where("loans.start_date <= ? AND loans.end_date >= ?", endDate, startDate).
		Order("loans.loan_status DESC, loans.created_at ASC"). // APPROVED first, then by time
		Find(&loanItems).Error

	if err != nil {
		return nil, err
	}

	return loanItems, nil
}

func (r *loanRepository) GetLoanForUpdate(loanID uint) (*model.Loan, error) {
	var loan model.Loan
	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("LoanItems").
		First(&loan, loanID).Error

	if err != nil {
		return nil, err
	}
	return &loan, nil
}

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

func (r *loanRepository) SetLoanRevision(
	loanID uint,
	note string,
) error {
	return r.db.Model(&model.Loan{}).
		Where("id = ?", loanID).
		Updates(map[string]interface{}{
			"loan_status":      "PENDING",
			"loan_flow_status": "REVISION",
		}).Error
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

func (r *loanRepository) FindByIDWithItems(
	loanID uint,
) (*model.Loan, error) {

	var loan model.Loan

	err := r.db.
		Preload("Items").
		Preload("Items.Barang").
		First(&loan, loanID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("peminjaman tidak ditemukan")
		}
		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) FindByIDAndUserIDWithItems(
	loanID uint,
	userID uint,
) (*model.Loan, error) {

	var loan model.Loan

	err := r.db.
		Where("id = ? AND user_id = ?", loanID, userID).
		Preload("Items").
		Preload("Items.Barang").
		First(&loan).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("peminjaman tidak ditemukan")
		}
		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) SoftDelete(id uint) error {
	return r.db.Delete(&model.Loan{}, id).Error
}

func (r *loanRepository) GetOverlappingLoanItems(
	barangID uint,
	startDate time.Time,
	endDate time.Time,
	excludeLoanID *uint,
) (int, error) {

	var total int

	query := r.db.
		Table("loan_items li").
		Select("COALESCE(SUM(li.quantity),0)").
		Joins("JOIN loans l ON l.id = li.loan_id").
		Where("li.barang_id = ?", barangID).
		Where(`
			l.start_date <= ? 
			AND l.end_date >= ?
		`, endDate, startDate).
		Where("l.loan_status = ?", "APPROVED").
		Where("l.loan_flow_status IN ?", []string{"READY", "TAKEN"})

	if excludeLoanID != nil {
		query = query.Where("l.id <> ?", *excludeLoanID)
	}

	if err := query.Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
