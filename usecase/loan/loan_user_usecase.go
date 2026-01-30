// usecase/loan_user_usecase.go
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
	UpdateHeader(loanCode string, userID uint, req req.UpdateLoanRequest) error
	UpdateItems(loanCode string, userID uint, req req.UpdateLoanItemsRequest) error
	Cancel(loanCode string, userID uint) error
	GetMyLoans(userID uint) ([]res.LoanListResponse, error)
	GetDetail(loanCode string, userID uint) (*res.LoanDetailResponse, error)
}

type loanUserUsecase struct {
	loanRepo       repository.LoanRepository
	cartRepo       repository.CartRepository
	barangRepo     repository.BarangRepository
	unitRepo       repository.BarangUnitRepository
	statusHistRepo repository.LoanStatusHistoryRepository
	flowHistRepo   repository.LoanFlowHistoryRepository
}

func NewLoanUserUsecase(
	loanRepo repository.LoanRepository,
	cartRepo repository.CartRepository,
	barangRepo repository.BarangRepository,
	unitRepo repository.BarangUnitRepository,
	statusHistRepo repository.LoanStatusHistoryRepository,
	flowHistRepo repository.LoanFlowHistoryRepository,
) LoanUserUsecase {
	return &loanUserUsecase{
		loanRepo:       loanRepo,
		cartRepo:       cartRepo,
		barangRepo:     barangRepo,
		unitRepo:       unitRepo,
		statusHistRepo: statusHistRepo,
		flowHistRepo:   flowHistRepo,
	}
}

//
// =====================================================
// CREATE LOAN (CHECKOUT CART)
// =====================================================
//
func (u *loanUserUsecase) Create(
	req req.CreateLoanRequest,
	userID uint,
) (*res.LoanResponse, error) {

	// 1. Get cart with items
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

	// 2. ✅ VALIDASI AVAILABILITY untuk SEMUA items
	unavailableItems := []string{}

	for _, ci := range cart.CartItems {
		// Cek barang aktif
		if ci.Barang.Status == "nonaktif" {
			unavailableItems = append(unavailableItems,
				fmt.Sprintf("%s: barang tidak aktif", ci.Barang.Nama))
			continue
		}

		// ✅ Cek available units di tanggal tersebut
		availableUnits, err := u.getAvailableUnitsInDateRange(
			ci.BarangID,
			req.ParsedStartDate,
			req.ParsedEndDate,
			nil, // excludeLoanID = nil (create baru)
		)

		if err != nil {
			return nil, err
		}

		if len(availableUnits) < ci.Quantity {
			unavailableItems = append(unavailableItems,
				fmt.Sprintf("%s: hanya %d unit tersedia (diminta %d) di tanggal %s - %s",
					ci.Barang.Nama,
					len(availableUnits),
					ci.Quantity,
					req.ParsedStartDate.Format("02 Jan"),
					req.ParsedEndDate.Format("02 Jan"),
				))
		}
	}

	// Kalau ada yang gak available, return error
	if len(unavailableItems) > 0 {
		return nil, errors.New("Beberapa barang tidak tersedia:\n" +
			fmt.Sprintf("%v", unavailableItems))
	}

	// 3. Create loan header
	loan := &model.Loan{
		UserID:         userID,
		LoanCode:       generateLoanCode(),
		StartDate:      req.ParsedStartDate,
		EndDate:        req.ParsedEndDate,
		Reason:         req.Reason,
		LoanStatus:     "PENDING",
		LoanFlowStatus: "REQUESTED",
	}

	// 4. Create loan items (belum assign units)
	items := make([]model.LoanItem, 0, len(cart.CartItems))
	for _, ci := range cart.CartItems {
		items = append(items, model.LoanItem{
			BarangID: ci.BarangID,
			Status:   "REQUESTED",
		})
	}

	// 5. Create loan + clear cart (transaction)
	if err := u.loanRepo.CreateLoanFromCartTx(loan, items, cart.ID); err != nil {
		return nil, err
	}

	// 6. ✅ AUTO-ASSIGN UNITS ke loan items
	for i, ci := range cart.CartItems {
		availableUnits, _ := u.getAvailableUnitsInDateRange(
			ci.BarangID,
			req.ParsedStartDate,
			req.ParsedEndDate,
			nil,
		)

		// Ambil sejumlah quantity yang diminta
		unitIDs := make([]uint, 0, ci.Quantity)
		for j := 0; j < ci.Quantity && j < len(availableUnits); j++ {
			unitIDs = append(unitIDs, availableUnits[j].ID)
		}

		// Assign units ke loan item
		if err := u.loanRepo.AssignUnitsToLoanItem(items[i].ID, unitIDs); err != nil {
			return nil, err
		}
	}

	// 7. ✅ RECORD HISTORY: Status awal
	if err := u.statusHistRepo.Create(&model.LoanStatusHistory{
		LoanID:        loan.ID,
		FromStatus:    nil,
		ToStatus:      "PENDING",
		ChangedBy:     userID,
		ChangedByRole: "USER",
		Note:          stringPtr("Loan dibuat oleh user"),
	}); err != nil {
		return nil, err
	}

	// 8. ✅ RECORD HISTORY: Flow awal
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:        loan.ID,
		FromFlow:      nil,
		ToFlow:        "REQUESTED",
		ChangedBy:     userID,
		ChangedByRole: "USER",
		Note:          stringPtr("Loan diajukan oleh user"),
	}); err != nil {
		return nil, err
	}

	// 9. Get created loan with items & units
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

	// ✅ Ganti validasi editable
	if !isEditable(loan.LoanFlowStatus) {
		return errors.New("loan tidak bisa diedit pada status saat ini")
	}

	// Store old values untuk history notes
	oldStartDate := loan.StartDate.Format("02 Jan 2006")
	oldEndDate := loan.EndDate.Format("02 Jan 2006")
	newStartDate := req.ParsedStartDate.Format("02 Jan 2006")
	newEndDate := req.ParsedEndDate.Format("02 Jan 2006")
	oldFlow := loan.LoanFlowStatus

	// ✅ Validasi availability untuk tanggal baru
	for _, item := range loan.LoanItems {
		availableUnits, err := u.getAvailableUnitsInDateRange(
			item.BarangID,
			req.ParsedStartDate,
			req.ParsedEndDate,
			&loan.ID, // exclude loan ini
		)

		if err != nil {
			return err
		}

		currentUnitCount := len(item.AssignedUnits)
		if len(availableUnits) < currentUnitCount {
			return fmt.Errorf(
				"%s: hanya %d unit tersedia di tanggal baru (saat ini: %d unit)",
				item.Barang.Nama,
				len(availableUnits),
				currentUnitCount,
			)
		}
	}

	// Update loan header
	loan.StartDate = req.ParsedStartDate
	loan.EndDate = req.ParsedEndDate
	loan.Reason = req.Reason

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	// ✅ RECORD HISTORY: Update header
	note := fmt.Sprintf("User mengubah tanggal dari %s - %s menjadi %s - %s",
		oldStartDate, oldEndDate, newStartDate, newEndDate)

	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:        loan.ID,
		FromFlow:      &oldFlow,
		ToFlow:        loan.LoanFlowStatus,
		ChangedBy:     userID,
		ChangedByRole: "USER",
		Note:          stringPtr(note),
	}); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// UPDATE ITEMS (ADD / UPDATE / DELETE)
// =====================================================
//
func (u *loanUserUsecase) UpdateItems(
	loanCode string,
	userID uint,
	req req.UpdateLoanItemsRequest,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.UserID != userID {
		return errors.New("akses ditolak")
	}

	// ✅ Ganti validasi editable
	if !isEditable(loan.LoanFlowStatus) {
		return errors.New("loan tidak bisa diedit pada status saat ini")
	}

	// Map existing items
	existing := map[uint]*model.LoanItem{}
	for i := range loan.LoanItems {
		existing[loan.LoanItems[i].BarangID] = &loan.LoanItems[i]
	}

	// Loop request items
	for _, item := range req.Items {

		switch item.Action {

		// ================= ADD =================
		case "ADD":
			if _, ok := existing[item.BarangID]; ok {
				return errors.New("barang sudah ada di loan")
			}

			// ✅ Cek availability
			availableUnits, err := u.getAvailableUnitsInDateRange(
				item.BarangID,
				loan.StartDate,
				loan.EndDate,
				&loan.ID,
			)

			if err != nil {
				return err
			}

			if len(availableUnits) < item.Quantity {
				return fmt.Errorf(
					"hanya %d unit tersedia (diminta %d)",
					len(availableUnits),
					item.Quantity,
				)
			}

			// Create loan item
			newItem := &model.LoanItem{
				LoanID:   loan.ID,
				BarangID: item.BarangID,
				Status:   "REQUESTED",
			}

			if err := u.loanRepo.CreateLoanItem(newItem); err != nil {
				return err
			}

			// ✅ Assign units
			unitIDs := make([]uint, 0, item.Quantity)
			for j := 0; j < item.Quantity && j < len(availableUnits); j++ {
				unitIDs = append(unitIDs, availableUnits[j].ID)
			}

			if err := u.loanRepo.AssignUnitsToLoanItem(newItem.ID, unitIDs); err != nil {
				return err
			}

			// ✅ RECORD HISTORY: Add item
			barang, _ := u.barangRepo.FindByID(item.BarangID)
			oldFlow := loan.LoanFlowStatus
			note := fmt.Sprintf("User menambah barang: %s (%d unit)", barang.Nama, item.Quantity)

			if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
				LoanID:        loan.ID,
				FromFlow:      &oldFlow,
				ToFlow:        loan.LoanFlowStatus,
				ChangedBy:     userID,
				ChangedByRole: "USER",
				Note:          stringPtr(note),
			}); err != nil {
				return err
			}

		// ================= UPDATE =================
		case "UPDATE":
			li, ok := existing[item.BarangID]
			if !ok {
				return errors.New("barang tidak ditemukan di loan")
			}

			oldQuantity := len(li.AssignedUnits)
			oldFlow := loan.LoanFlowStatus

			// ✅ Cek availability
			availableUnits, err := u.getAvailableUnitsInDateRange(
				item.BarangID,
				loan.StartDate,
				loan.EndDate,
				&loan.ID,
			)

			if err != nil {
				return err
			}

			if len(availableUnits) < item.Quantity {
				return fmt.Errorf(
					"hanya %d unit tersedia (diminta %d)",
					len(availableUnits),
					item.Quantity,
				)
			}

			// ✅ Re-assign units
			unitIDs := make([]uint, 0, item.Quantity)
			for j := 0; j < item.Quantity && j < len(availableUnits); j++ {
				unitIDs = append(unitIDs, availableUnits[j].ID)
			}

			if err := u.loanRepo.AssignUnitsToLoanItem(li.ID, unitIDs); err != nil {
				return err
			}

			// ✅ RECORD HISTORY: Update item
			note := fmt.Sprintf("User mengubah quantity %s dari %d menjadi %d unit",
				li.Barang.Nama, oldQuantity, item.Quantity)

			if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
				LoanID:        loan.ID,
				FromFlow:      &oldFlow,
				ToFlow:        loan.LoanFlowStatus,
				ChangedBy:     userID,
				ChangedByRole: "USER",
				Note:          stringPtr(note),
			}); err != nil {
				return err
			}

		// ================= DELETE =================
		case "DELETE":
			li, ok := existing[item.BarangID]
			if !ok {
				return errors.New("barang tidak ditemukan di loan")
			}

			oldQuantity := len(li.AssignedUnits)
			oldFlow := loan.LoanFlowStatus

			// ✅ Remove units assignment
			if err := u.loanRepo.RemoveUnitsFromLoanItem(li.ID); err != nil {
				return err
			}

			// Delete loan item
			if err := u.loanRepo.DeleteLoanItem(loan.ID, item.BarangID); err != nil {
				return err
			}

			// ✅ RECORD HISTORY: Delete item
			note := fmt.Sprintf("User menghapus barang: %s (%d unit)",
				li.Barang.Nama, oldQuantity)

			if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
				LoanID:        loan.ID,
				FromFlow:      &oldFlow,
				ToFlow:        loan.LoanFlowStatus,
				ChangedBy:     userID,
				ChangedByRole: "USER",
				Note:          stringPtr(note),
			}); err != nil {
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

	// Store old values
	oldStatus := loan.LoanStatus
	oldFlow := loan.LoanFlowStatus

	// ✅ Remove all unit assignments sebelum delete
	for _, item := range loan.LoanItems {
		_ = u.loanRepo.RemoveUnitsFromLoanItem(item.ID)
	}

	// ✅ RECORD HISTORY: Cancel loan - Status
	if err := u.statusHistRepo.Create(&model.LoanStatusHistory{
		LoanID:        loan.ID,
		FromStatus:    &oldStatus,
		ToStatus:      "REJECTED",
		ChangedBy:     userID,
		ChangedByRole: "USER",
		Note:          stringPtr("Loan dibatalkan oleh user"),
	}); err != nil {
		return err
	}

	// ✅ RECORD HISTORY: Cancel loan - Flow
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:        loan.ID,
		FromFlow:      &oldFlow,
		ToFlow:        "REQUESTED", // tetap REQUESTED karena cancelled
		ChangedBy:     userID,
		ChangedByRole: "USER",
		Note:          stringPtr("Loan dibatalkan oleh user"),
	}); err != nil {
		return err
	}

	// Soft delete loan
	if err := u.loanRepo.SoftDelete(loan.ID); err != nil {
		return err
	}

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
		// ✅ Hitung total units (bukan total items)
		totalUnits := 0
		for _, item := range loan.LoanItems {
			totalUnits += len(item.AssignedUnits)
		}

		resp = append(resp, res.LoanListResponse{
			LoanCode:       loan.LoanCode,
			LoanStatus:     loan.LoanStatus,
			LoanFlowStatus: loan.LoanFlowStatus,
			StartDate:      loan.StartDate,
			EndDate:        loan.EndDate,
			TotalItems:     len(loan.LoanItems),
			TotalUnits:     totalUnits,
			CreatedAt:      loan.CreatedAt,
		})
	}

	return resp, nil
}

//
// =====================================================
// GET DETAIL
// =====================================================
//
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
	totalUnits := 0

	for _, item := range loan.LoanItems {
		// ✅ Build unit details
		unitDetails := make([]res.UnitDetail, 0)
		for _, unit := range item.AssignedUnits {
			unitDetails = append(unitDetails, res.UnitDetail{
				UnitID:   unit.ID,
				KodeUnit: unit.KodeUnit,
				Kondisi:  unit.Kondisi,
				Status:   unit.Status,
			})
		}

		totalUnits += len(item.AssignedUnits)

		items = append(items, res.LoanItemDetailResponse{
			BarangID:   item.BarangID,
			BarangKode: item.Barang.Kode,
			BarangNama: item.Barang.Nama,
			Merk:       item.Barang.Merk,
			Kategori:   item.Barang.Kategori,
			Quantity:   len(item.AssignedUnits),
			Units:      unitDetails,
			Status:     item.Status,
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
		TotalItems:     len(loan.LoanItems),
		TotalUnits:     totalUnits,
		Items:          items,
		CreatedAt:      loan.CreatedAt,
	}, nil
}

//
// =====================================================
// HELPER FUNCTIONS
// =====================================================
//

// ✅ NEW: Get available units in date range
func (u *loanUserUsecase) getAvailableUnitsInDateRange(
	barangID uint,
	startDate time.Time,
	endDate time.Time,
	excludeLoanID *uint,
) ([]model.BarangUnit, error) {

	// 1. Get all units dari barang ini
	allUnits, err := u.unitRepo.FindAllByBarangID(barangID)
	if err != nil {
		return nil, err
	}

	// 2. Get overlapping unit IDs (yang udah dipinjam)
	bookedUnitIDs, err := u.loanRepo.GetOverlappingUnits(
		barangID,
		startDate,
		endDate,
		excludeLoanID,
	)
	if err != nil {
		return nil, err
	}

	// 3. Filter: ambil yang TIDAK ada di bookedUnitIDs
	bookedMap := make(map[uint]bool)
	for _, id := range bookedUnitIDs {
		bookedMap[id] = true
	}

	availableUnits := make([]model.BarangUnit, 0)
	for _, unit := range allUnits {
		// Skip kalau unit maintenance atau booked
		if unit.Status == "maintenance_required" {
			continue
		}

		if bookedMap[unit.ID] {
			continue
		}

		availableUnits = append(availableUnits, unit)
	}

	return availableUnits, nil
}

// ✅ NEW: Check if loan editable by user
func isEditable(flow string) bool {
	return flow == "REQUESTED" || flow == "REVISION"
}

// ✅ NEW: Helper untuk convert string ke pointer
func stringPtr(s string) *string {
	return &s
}

func generateLoanCode() string {
	now := time.Now()
	timePart := now.Format("1504")
	datePart := now.Format("020106")
	seq := now.UnixNano() % 1000
	return fmt.Sprintf("LN%s%s%03d", timePart, datePart, seq)
}

func toLoanResponse(loan *model.Loan) *res.LoanResponse {
	items := make([]res.LoanItemResponse, 0)
	totalUnits := 0

	for _, item := range loan.LoanItems {
		totalUnits += len(item.AssignedUnits)

		items = append(items, res.LoanItemResponse{
			BarangID:   item.BarangID,
			BarangKode: item.Barang.Kode,
			BarangNama: item.Barang.Nama,
			Quantity:   len(item.AssignedUnits),
			Status:     item.Status,
		})
	}

	return &res.LoanResponse{
		LoanCode:       loan.LoanCode,
		LoanStatus:     loan.LoanStatus,
		LoanFlowStatus: loan.LoanFlowStatus,
		StartDate:      loan.StartDate,
		EndDate:        loan.EndDate,
		Reason:         loan.Reason,
		TotalItems:     len(loan.LoanItems),
		TotalUnits:     totalUnits,
		Items:          items,
		CreatedAt:      loan.CreatedAt,
	}
}