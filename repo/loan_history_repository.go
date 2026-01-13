package repo

import (
	"backend/model"

	"gorm.io/gorm"
)

type LoanStatusHistoryRepository interface {
	Create(history *model.LoanStatusHistory) error
	GetByLoanID(loanID uint) ([]model.LoanStatusHistory, error)
}

type loanStatusHistoryRepository struct {
	db *gorm.DB
}

func NewLoanStatusHistoryRepository(db *gorm.DB) LoanStatusHistoryRepository {
	return &loanStatusHistoryRepository{db: db}
}

func (r *loanStatusHistoryRepository) Create(history *model.LoanStatusHistory) error {
	return r.db.Create(history).Error
}

func (r *loanStatusHistoryRepository) GetByLoanID(loanID uint) ([]model.LoanStatusHistory, error) {
	var histories []model.LoanStatusHistory

	err := r.db.
		Where("loan_id = ?", loanID).
		Order("created_at ASC").
		Find(&histories).Error

	return histories, err
}
