package repo

import (
	"backend/model"

	"gorm.io/gorm"
)

type LoanFlowHistoryRepository interface {
	Create(history *model.LoanFlowHistory) error
	GetByLoanID(loanID uint) ([]model.LoanFlowHistory, error)
}

type loanFlowHistoryRepository struct {
	db *gorm.DB
}

func NewLoanFlowHistoryRepository(db *gorm.DB) LoanFlowHistoryRepository {
	return &loanFlowHistoryRepository{db: db}
}

func (r *loanFlowHistoryRepository) Create(history *model.LoanFlowHistory) error {
	return r.db.Create(history).Error
}

func (r *loanFlowHistoryRepository) GetByLoanID(loanID uint) ([]model.LoanFlowHistory, error) {
	var histories []model.LoanFlowHistory

	err := r.db.
		Where("loan_id = ?", loanID).
		Order("created_at ASC").
		Find(&histories).Error

	return histories, err
}
