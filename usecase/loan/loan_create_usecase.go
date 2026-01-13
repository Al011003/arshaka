package usecase

import (
	req "backend/dto/request/loan"
	res "backend/dto/response/loan"
	"backend/model"
	repository "backend/repo"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LoanUsecase interface {
	CreateLoanFromCart(userID uint, req req.CreateLoanRequest) (*res.LoanResponse, error)
}

type loanUsecase struct {
	loanRepo              repository.LoanRepository
	cartRepo              repository.CartRepository
	loanStatusHistoryRepo repository.LoanStatusHistoryRepository
	loanFlowHistoryRepo   repository.LoanFlowHistoryRepository
}

func NewLoanUsecase(
	loanRepo repository.LoanRepository,
	cartRepo repository.CartRepository,
	loanStatusHistoryRepo repository.LoanStatusHistoryRepository,
	loanFlowHistoryRepo repository.LoanFlowHistoryRepository,
) LoanUsecase {
	return &loanUsecase{
		loanRepo:              loanRepo,
		cartRepo:              cartRepo,
		loanStatusHistoryRepo: loanStatusHistoryRepo,
		loanFlowHistoryRepo:   loanFlowHistoryRepo,
	}
}

// =====================================================
// CREATE LOAN FROM CART
// =====================================================
func (u *loanUsecase) CreateLoanFromCart(
	userID uint,
	request req.CreateLoanRequest,
) (*res.LoanResponse, error) {

	// 1. Ambil cart + items
	cart, err := u.cartRepo.GetCartWithItems(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart tidak ditemukan")
		}
		return nil, err
	}

	// 2. Validasi cart tidak kosong
	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart kosong, silakan tambahkan barang")
	}

	startDate := request.ParsedStartDate
	endDate := request.ParsedEndDate

	// 3. Validasi ketersediaan barang
	for _, ci := range cart.CartItems {
		if ci.Barang == nil {
			return nil, errors.New("data barang tidak lengkap")
		}

		if ci.Barang.Status != "tersedia" {
			return nil, fmt.Errorf("barang %s tidak aktif", ci.Barang.Nama)
		}

		borrowedQty, err := u.loanRepo.GetTotalBorrowedQtyByDate(
			ci.BarangID,
			startDate,
			endDate,
		)
		if err != nil {
			return nil, err
		}

		available := ci.Barang.StokTotal - int(borrowedQty)
		if ci.Quantity > available {
			return nil, fmt.Errorf(
				"stok %s tidak cukup (tersedia %d, diminta %d)",
				ci.Barang.Nama,
				available,
				ci.Quantity,
			)
		}
	}

	// 4. Create loan header
	loanCode := u.generateLoanCode()

// 5. Create loan header
	loan := &model.Loan{
		UserID:         userID,
		LoanCode:       loanCode,
		StartDate:      startDate,
		EndDate:        endDate,
		Reason:         request.Reason,
		LoanStatus:     "PENDING",
		LoanFlowStatus: "REQUESTED",
	}

	// 5. Create loan items
	items := make([]model.LoanItem, 0, len(cart.CartItems))
	for _, ci := range cart.CartItems {
		items = append(items, model.LoanItem{
			BarangID: ci.BarangID,
			Quantity: ci.Quantity,
			Status:   "REQUESTED",
		})
	}

	// 6. Transaction (loan + items + clear cart)
	if err := u.loanRepo.CreateLoanFromCartTx(loan, items, cart.ID); err != nil {

		// retry sekali kalau loan_code duplicate
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			loan.LoanCode = u.generateLoanCode()
			if err := u.loanRepo.CreateLoanFromCartTx(loan, items, cart.ID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}



	// =====================================================
	// 7. CREATE HISTORY (SETELAH TX SUKSES)
	// =====================================================

	// Loan Flow History
	toFlow := "REQUESTED"
	note := "User mengajukan peminjaman"

	if err := u.loanFlowHistoryRepo.Create(&model.LoanFlowHistory{
		LoanID:         loan.ID,
		FromFlow:       nil,
		ToFlow:         toFlow,
		ChangedBy:      userID,
		ChangedByRole: "SYSTEM",
		Note:           &note,
	}); err != nil {
		return nil, err
	}

	// Loan Status History
	noteL := "Loan dibuat oleh user"

	if err := u.loanStatusHistoryRepo.Create(&model.LoanStatusHistory{
		LoanID:         loan.ID,
		FromStatus:     nil,
		ToStatus:       "PENDING",
		ChangedBy:      userID,
		ChangedByRole: "SYSTEM",
		Note:           &noteL,
	}); err != nil {
		return nil, err
	}

	// 8. Ambil loan lengkap
	createdLoan, err := u.loanRepo.GetLoanByIDWithItems(loan.ID)
	if err != nil {
		return nil, err
	}

	return toLoanResponse(createdLoan), nil
}

// =====================================================
// HELPERS
// =====================================================
func toLoanResponse(loan *model.Loan) *res.LoanResponse {
	items := make([]res.LoanItemResponse, 0, len(loan.LoanItems))

	for _, item := range loan.LoanItems {
		items = append(items, *toLoanItemResponse(&item))
	}

	return &res.LoanResponse{
		ID:             loan.ID,
		UserID:         loan.UserID,
		StartDate:      loan.StartDate,
		EndDate:        loan.EndDate,
		Reason:         loan.Reason,
		LoanStatus:     loan.LoanStatus,
		LoanFlowStatus: loan.LoanFlowStatus,
		TotalItems:     len(items),
		LoanItems:      items,
		CreatedAt:      loan.CreatedAt,
		UpdatedAt:      loan.UpdatedAt,
		LoanCode: loan.LoanCode,

	}
}

func toLoanItemResponse(item *model.LoanItem) *res.LoanItemResponse {
	var barangResp *res.BarangInLoanResponse

	if item.Barang != nil {
		barangResp = &res.BarangInLoanResponse{
			ID:       item.Barang.ID,
			Kode:     item.Barang.Kode,
			Nama:     item.Barang.Nama,
			Merk:     item.Barang.Merk,
			Kategori: item.Barang.Kategori,
			CoverURL: item.Barang.CoverURL,
		}
	}

	return &res.LoanItemResponse{
		ID:       item.ID,
		LoanID:   item.LoanID,
		BarangID: item.BarangID,
		Barang:   barangResp,
		Quantity: item.Quantity,
		Status:   item.Status,
	}
}

func (u *loanUsecase) generateLoanCode() string {
	now := time.Now()

	timePart := now.Format("1504")   // HHMM
	datePart := now.Format("020106") // DDMMYY

	// SEMENTARA: pakai unix nano mod 1000
	// NOTE: ini aman untuk sekarang, nanti bisa diganti counter DB
	seq := now.UnixNano() % 1000

	return fmt.Sprintf("AB%s%s%03d", timePart, datePart, seq)
}