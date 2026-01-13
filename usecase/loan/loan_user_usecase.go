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

type LoanUserUsecase interface {
	Create(req req.CreateLoanRequest, userID uint) (*res.LoanResponse, error)

	UpdateHeader(
		loanCode string,
		userID uint,
		req req.UpdateLoanRequest,
	) error

	UpdateItems(
		loanCode string,
		userID uint,
		req req.UpdateLoanItemsRequest,
	) error

	Cancel(
	loanCode string,
	userID uint,
) error
	GetMyLoans(userID uint) ([]res.LoanListResponse, error)
	 GetDetail(
	loanCode string,
	userID uint,
) (*res.LoanDetailResponse, error)

}

type loanUserUsecase struct {
	loanRepo              repository.LoanRepository
	cartRepo              repository.CartRepository
	barangRepo 			repository.BarangRepository
	loanStatusHistoryRepo repository.LoanStatusHistoryRepository
	loanFlowHistoryRepo   repository.LoanFlowHistoryRepository
}

func NewLoanUserUsecase(
	loanRepo repository.LoanRepository,
	cartRepo repository.CartRepository,
	barangRepo 			repository.BarangRepository,
	loanStatusHistoryRepo repository.LoanStatusHistoryRepository,
	loanFlowHistoryRepo repository.LoanFlowHistoryRepository,
) LoanUserUsecase {
	return &loanUserUsecase{
		loanRepo:              loanRepo,
		cartRepo:              cartRepo,
		barangRepo: 			barangRepo,
		loanStatusHistoryRepo: loanStatusHistoryRepo,
		loanFlowHistoryRepo:   loanFlowHistoryRepo,
	}
}

//
// =====================================================
// CREATE LOAN (USER)
// =====================================================
//
func (u *loanUserUsecase) Create(
	req req.CreateLoanRequest,
	userID uint,
) (*res.LoanResponse, error) {

	cart, err := u.cartRepo.GetCartWithItems(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart tidak ditemukan")
		}
		return nil, err
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart kosong")
	}

	for _, ci := range cart.CartItems {

	err := u.checkOverlapAndStock(
		ci.BarangID,
		req.ParsedStartDate,
		req.ParsedEndDate,
		ci.Quantity,
		nil, // excludeLoanID = nil (karena create)
	)

	if err != nil {
		return nil, fmt.Errorf(
			"barang %d: %s",
			ci.BarangID,
			err.Error(),
		)
	}
}


	loan := &model.Loan{
		UserID:         userID,
		LoanCode:       generateLoanCode(),
		StartDate:      req.ParsedStartDate,
		EndDate:        req.ParsedEndDate,
		Reason:         req.Reason,
		LoanStatus:     "PENDING",
		LoanFlowStatus: "REQUESTED",
	}

	items := make([]model.LoanItem, 0, len(cart.CartItems))
	for _, ci := range cart.CartItems {
		items = append(items, model.LoanItem{
			BarangID: ci.BarangID,
			Quantity: ci.Quantity,
			Status:   "REQUESTED",
		})
	}

	if err := u.loanRepo.CreateLoanFromCartTx(loan, items, cart.ID); err != nil {
		return nil, err
	}

	_ = u.loanFlowHistoryRepo.Create(&model.LoanFlowHistory{
		LoanID:         loan.ID,
		FromFlow:       nil,
		ToFlow:         "REQUESTED",
		ChangedBy:      userID,
		ChangedByRole: "SYSTEM",
	})

	_ = u.loanStatusHistoryRepo.Create(&model.LoanStatusHistory{
		LoanID:         loan.ID,
		FromStatus:     nil,
		ToStatus:       "PENDING",
		ChangedBy:      userID,
		ChangedByRole: "SYSTEM",
	})

	createdLoan, err := u.loanRepo.GetLoanByIDWithItems(loan.ID)
	if err != nil {
		return nil, err
	}

	return toLoanResponse(createdLoan), nil
}

//
// =====================================================
// UPDATE HEADER (TANGGAL + REASON)
// =====================================================
//
func (u *loanUserUsecase) UpdateHeader(
	loanCode string,
	userID uint,
	req req.UpdateLoanRequest,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.UserID != userID {
		return errors.New("akses ditolak")
	}

	if !isEditableByUser(loan) {
		return errors.New("loan tidak bisa diedit pada status saat ini")
	}

	for _, item := range loan.LoanItems {

	err := u.checkOverlapAndStock(
		item.BarangID,
		req.ParsedStartDate,
		req.ParsedEndDate,
		item.Quantity,
		&loan.ID, // exclude dirinya sendiri
	)

	if err != nil {
		return err
	}
}


	loan.StartDate = req.ParsedStartDate
	loan.EndDate = req.ParsedEndDate
	loan.Reason = req.Reason

	return u.loanRepo.UpdateLoan(loan)
}

//
// =====================================================
// UPDATE ITEMS (ADD / UPDATE / DELETE)
// =====================================================
func (u *loanUserUsecase) UpdateItems(
	loanCode string,
	userID uint,
	req req.UpdateLoanItemsRequest,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	// ownership
	if loan.UserID != userID {
		return errors.New("akses ditolak")
	}

	// status editable
	if !isEditableByUser(loan) {
		return errors.New("loan tidak bisa diedit pada status saat ini")
	}

	// map existing items
	existing := map[uint]*model.LoanItem{}
	for i := range loan.LoanItems {
		existing[loan.LoanItems[i].BarangID] = &loan.LoanItems[i]
	}

	// loop request items
	for _, item := range req.Items {

		switch item.Action {

		// ================= ADD =================
		case "ADD":
			if _, ok := existing[item.BarangID]; ok {
				return errors.New("barang sudah ada di loan")
			}

			// 🔥 VALIDASI BARU
			if err := u.checkOverlapAndStock(
				item.BarangID,
				loan.StartDate,
				loan.EndDate,
				item.Quantity,
				&loan.ID,
			); err != nil {
				return err
			}

			if err := u.loanRepo.CreateLoanItem(&model.LoanItem{
				LoanID:   loan.ID,
				BarangID: item.BarangID,
				Quantity: item.Quantity,
				Status:   "REQUESTED",
			}); err != nil {
				return err
			}

		// ================= UPDATE =================
		case "UPDATE":
			li, ok := existing[item.BarangID]
			if !ok {
				return errors.New("barang tidak ditemukan di loan")
			}

			// 🔥 VALIDASI BARU
			if err := u.checkOverlapAndStock(
				item.BarangID,
				loan.StartDate,
				loan.EndDate,
				item.Quantity,
				&loan.ID,
			); err != nil {
				return err
			}

			li.Quantity = item.Quantity
			if err := u.loanRepo.UpdateLoanItem(li); err != nil {
				return err
			}

		// ================= DELETE =================
		case "DELETE":
			if _, ok := existing[item.BarangID]; !ok {
				return errors.New("barang tidak ditemukan di loan")
			}

			// ❌ DELETE TIDAK PERLU OVERLAP CHECK
			if err := u.loanRepo.DeleteLoanItem(loan.ID, item.BarangID); err != nil {
				return err
			}

		default:
			return errors.New("action tidak valid")
		}
	}

	return nil
}


//
// =====================================================
// CANCEL LOAN (USER)
// =====================================================
//
func (u *loanUserUsecase) Cancel(
	loanCode string,
	userID uint,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.UserID != userID {
		return errors.New("akses ditolak")
	}

	if loan.LoanStatus != "PENDING" {
		return errors.New("loan tidak bisa dibatalkan")
	}

	if err := u.loanRepo.SoftDelete(loan.ID); err != nil {
		return err
	}

	note := "Dibatalkan oleh user"
	_ = u.loanStatusHistoryRepo.Create(&model.LoanStatusHistory{
		LoanID:         loan.ID,
		FromStatus:     &loan.LoanStatus,
		ToStatus:       "CANCELLED",
		ChangedBy:      userID,
		ChangedByRole: "SYSTEM",
		Note:           &note,
	})

	return nil
}


//
// =====================================================
// GET MY LOANS
// =====================================================
//
func (u *loanUserUsecase) GetMyLoans(
	userID uint,
) ([]res.LoanListResponse, error) {

	loans, err := u.loanRepo.GetLoanByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := make([]res.LoanListResponse, 0, len(loans))
	for _, loan := range loans {
		resp = append(resp, res.LoanListResponse{
			LoanCode:       loan.LoanCode,
			LoanStatus:     loan.LoanStatus,
			LoanFlowStatus: loan.LoanFlowStatus,
			StartDate:      loan.StartDate,
			EndDate:        loan.EndDate,
			TotalItems:     len(loan.LoanItems),
			CreatedAt:      loan.CreatedAt,
		})
	}

	return resp, nil
}

func (u *loanUserUsecase) GetDetail(
	loanCode string,
	userID uint,
) (*res.LoanDetailResponse, error) {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return nil, err
	}

	if loan.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	items := make([]res.LoanItemDetailResponse, 0)
	totalItems := 0

	for _, item := range loan.LoanItems {
		totalItems += item.Quantity

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
		TotalItems:     totalItems,
		Items:          items,
		CreatedAt:      loan.CreatedAt,
	}, nil
}

//
// =====================================================
// HELPERS
// =====================================================
//
func isEditableByUser(loan *model.Loan) bool {
	if loan.LoanStatus == "PENDING" {
		return true
	}

	if loan.LoanFlowStatus == "REVISION" {
		return true
	}

	return false
}





func generateLoanCode() string {
	now := time.Now()

	timePart := now.Format("1504")
	datePart := now.Format("020106")
	seq := now.UnixNano() % 1000

	return fmt.Sprintf("AB%s%s%03d", timePart, datePart, seq)
}

func (u *loanUserUsecase) checkOverlapAndStock(
	barangID uint,
	start time.Time,
	end time.Time,
	requestQty int,
	excludeLoanID *uint,
) error {

	// 1. barang aktif?
	active, err := u.barangRepo.IsActive(barangID)
	if err != nil {
		return err
	}
	if !active {
		return errors.New("barang tidak aktif")
	}

	// 2. stok total
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil {
		return err
	}

	// 3. overlap qty
	usedQty, err := u.loanRepo.GetOverlappingLoanItems(
		barangID,
		start,
		end,
		excludeLoanID,
	)
	if err != nil {
		return err
	}

	if usedQty+requestQty > barang.StokTotal {
		return errors.New("stok tidak mencukupi pada tanggal tersebut")
	}

	return nil
}
