package usecase

import (
	res "backend/dto/response/loan"
	"backend/model"
	repository "backend/repo"
	"errors"

	"gorm.io/gorm"
)

type LoanGetUsecase interface {
	GetMyLoans(userID uint) ([]res.LoanListResponse, error)
	GetLoanDetailByCode(userID uint, loanCode string) (*res.LoanDetailResponse, error)
}

type loanGetUsecase struct {
	loanRepo repository.LoanRepository
}

func NewLoanGetUsecase(
	loanRepo repository.LoanRepository,
) LoanGetUsecase {
	return &loanGetUsecase{
		loanRepo: loanRepo,
	}
}

func (u *loanGetUsecase) GetMyLoans(userID uint) ([]res.LoanListResponse, error) {
	loans, err := u.loanRepo.GetLoanByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []res.LoanListResponse{}, nil
		}
		return nil, err
	}

	result := make([]res.LoanListResponse, 0, len(loans))
	for _, loan := range loans {
		result = append(result, *toLoanListResponse(&loan))
	}

	return result, nil
}

func (u *loanGetUsecase) GetLoanDetailByCode(
	userID uint,
	loanCode string,
) (*res.LoanDetailResponse, error) {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan tidak ditemukan")
		}
		return nil, err
	}

	// 🔐 Ownership check
	if loan.UserID != userID {
		return nil, errors.New("tidak memiliki akses ke loan ini")
	}

	return toLoanDetailResponse(loan), nil
}

func toLoanListResponse(loan *model.Loan) *res.LoanListResponse {
	return &res.LoanListResponse{
		LoanCode:       loan.LoanCode,
		LoanStatus:     loan.LoanStatus,
		LoanFlowStatus: loan.LoanFlowStatus,
		StartDate:      loan.StartDate,
		EndDate:        loan.EndDate,
		TotalItems:     len(loan.LoanItems),
		CreatedAt:      loan.CreatedAt,
	}
}


func toLoanDetailResponse(loan *model.Loan) *res.LoanDetailResponse {
	items := make([]res.LoanItemDetailResponse, 0, len(loan.LoanItems))

	for _, item := range loan.LoanItems {
		if item.Barang == nil {
			continue
		}

		items = append(items, res.LoanItemDetailResponse{
			BarangID: item.BarangID,
			Nama:     item.Barang.Nama,
			Merk:     item.Barang.Merk,
			Kategori: item.Barang.Kategori,
			Quantity: item.Quantity,
			Status:   item.Status,
		})
	}

	return &res.LoanDetailResponse{
		LoanCode:       loan.LoanCode,
		UserID:         loan.UserID,
		LoanStatus:     loan.LoanStatus,
		LoanFlowStatus: loan.LoanFlowStatus,
		StartDate:      loan.StartDate,
		EndDate:        loan.EndDate,
		Reason:         loan.Reason,
		TotalItems:     len(items),
		Items:          items,
		CreatedAt:      loan.CreatedAt,
	}
}
