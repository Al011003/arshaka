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
	
	// 🔥 NEW: Check availability with gap rule & suggestions
	CheckAvailabilityWithSuggestions(KodeBarang string, quantity int, startDate, endDate time.Time) (*res.AvailabilityCheckResponse, error)
	
	// 🔥 NEW: Auto-revision handler
	MoveLoansToRevision(barangID uint, adminID uint) error
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

	// 2. ✅ VALIDASI AVAILABILITY dengan GAP RULE + SUGGESTIONS
	unavailableItems := []string{}
	suggestions := []string{}

	for _, ci := range cart.CartItems {
		// Cek barang aktif
		if ci.Barang.Status == "nonaktif" {
			unavailableItems = append(unavailableItems,
				fmt.Sprintf("%s: barang tidak aktif", ci.Barang.Nama))
			
			// 🔥 NEW: Cari rekomendasi barang serupa
			similarBarangs, _ := u.barangRepo.FindSimilarBarang(
				ci.Barang.Kategori,
				ci.Barang.Merk,
				ci.BarangID,
				3, // max 3 rekomendasi
				ci.Barang.Nama, 
			)
			
			if len(similarBarangs) > 0 {
				suggestionList := "Rekomendasi barang serupa: "
				for i, sb := range similarBarangs {
					if i > 0 {
						suggestionList += ", "
					}
					suggestionList += fmt.Sprintf("%s (%s)", sb.Nama, sb.Merk)
				}
				suggestions = append(suggestions, suggestionList)
			}
			
			continue
		}

		// ✅ Cek available units dengan GAP RULE
		availCheck, err := u.checkAvailabilityWithGapRule(
			ci.BarangID,
			ci.Quantity,
			req.ParsedStartDate,
			req.ParsedEndDate,
			nil, // excludeLoanID = nil (create baru)
			
		)

		if err != nil {
			return nil, err
		}

		if !availCheck.IsAvailable {
			msg := fmt.Sprintf("%s: %s", ci.Barang.Nama, availCheck.Message)
			unavailableItems = append(unavailableItems, msg)
			
			// Tambahkan suggestion kalau ada
			if availCheck.SuggestedDate != nil {
				suggestions = append(suggestions, 
					fmt.Sprintf("💡 %s: Tersedia %d unit mulai %s",
						ci.Barang.Nama,
						ci.Quantity,
						availCheck.SuggestedDate.Format("02 Jan 2006"),
					))
			}
			
			// 🔥 NEW: Cari barang alternatif
			if len(availCheck.AlternativeBarangs) > 0 {
				altList := fmt.Sprintf("Alternatif untuk %s: ", ci.Barang.Nama)
				for i, alt := range availCheck.AlternativeBarangs {
					if i > 0 {
						altList += ", "
					}
					altList += fmt.Sprintf("%s (%s)", alt.BarangNama, alt.Merk)
				}
				suggestions = append(suggestions, altList)
			}
		}
	}

	// Kalau ada yang gak available, return error dengan suggestions
	if len(unavailableItems) > 0 {
		errorMsg := "Beberapa barang tidak tersedia:\n"
		for _, item := range unavailableItems {
			errorMsg += "- " + item + "\n"
		}
		
		if len(suggestions) > 0 {
			errorMsg += "\n📋 Saran:\n"
			for _, sug := range suggestions {
				errorMsg += "- " + sug + "\n"
			}
		}
		
		return nil, errors.New(errorMsg)
	}

	// 3. Create loan header
	loan := &model.Loan{
		UserID:         userID,
		LoanCode:       generateLoanCode(),
		StartDate:      req.ParsedStartDate,
		EndDate:        req.ParsedEndDate,
		Reason:         req.Reason,
		LoanStatus:     model.LoanStatusPending,
		LoanFlowStatus: model.LoanFlowRequested,
	}

	// 4. Create loan items (belum assign units)
	items := make([]model.LoanItem, 0, len(cart.CartItems))
	for _, ci := range cart.CartItems {
		items = append(items, model.LoanItem{
			BarangID: ci.BarangID,
			Status:   model.LoanItemStatusRequested,
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
		ToStatus:      model.LoanStatusPending,
		ChangedBy:     userID,
		ChangedByRole: model.HistoryRoleUser,
		Note:          stringPtr("Loan dibuat oleh user"),
	}); err != nil {
		return nil, err
	}

	// 8. ✅ RECORD HISTORY: Flow awal
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        nil,
		ToFlow:          model.LoanFlowRequested,
		TriggeredBy:     userID,
		TriggeredByRole: model.HistoryRoleUser,
		Note:            "Loan diajukan oleh user",
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

	// ✅ Validasi availability dengan GAP RULE
	for _, item := range loan.LoanItems {
		currentQty := len(item.AssignedUnits)
		
		availCheck, err := u.checkAvailabilityWithGapRule(
			item.BarangID,
			currentQty,
			req.ParsedStartDate,
			req.ParsedEndDate,
			&loan.ID, // exclude loan ini
		)

		if err != nil {
			return err
		}

		if !availCheck.IsAvailable {
			errorMsg := fmt.Sprintf("%s: %s", item.Barang.Nama, availCheck.Message)
			
			if availCheck.SuggestedDate != nil {
				errorMsg += fmt.Sprintf(" (tersedia mulai %s)", 
					availCheck.SuggestedDate.Format("02 Jan 2006"))
			}
			
			return errors.New(errorMsg)
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
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          loan.LoanFlowStatus,
		TriggeredBy:     userID,
		TriggeredByRole: model.HistoryRoleUser,
		Note:            note,
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

			// ✅ Cek availability dengan GAP RULE
			availCheck, err := u.checkAvailabilityWithGapRule(
				item.BarangID,
				item.Quantity,
				loan.StartDate,
				loan.EndDate,
				&loan.ID,
			)

			if err != nil {
				return err
			}

			if !availCheck.IsAvailable {
				errorMsg := availCheck.Message
				
				if availCheck.SuggestedDate != nil {
					errorMsg += fmt.Sprintf(" (tersedia mulai %s)", 
						availCheck.SuggestedDate.Format("02 Jan 2006"))
				}
				
				return errors.New(errorMsg)
			}

			// Get available units
			availableUnits, _ := u.getAvailableUnitsInDateRange(
				item.BarangID,
				loan.StartDate,
				loan.EndDate,
				&loan.ID,
			)

			// Create loan item
			newItem := &model.LoanItem{
				LoanID:   loan.ID,
				BarangID: item.BarangID,
				Status:   model.LoanItemStatusRequested,
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
				LoanID:          loan.ID,
				FromFlow:        &oldFlow,
				ToFlow:          loan.LoanFlowStatus,
				TriggeredBy:     userID,
				TriggeredByRole: model.HistoryRoleUser,
				Note:            note,
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

			// ✅ Cek availability dengan GAP RULE
			availCheck, err := u.checkAvailabilityWithGapRule(
				item.BarangID,
				item.Quantity,
				loan.StartDate,
				loan.EndDate,
				&loan.ID,
			)

			if err != nil {
				return err
			}

			if !availCheck.IsAvailable {
				errorMsg := availCheck.Message
				
				if availCheck.SuggestedDate != nil {
					errorMsg += fmt.Sprintf(" (tersedia mulai %s)", 
						availCheck.SuggestedDate.Format("02 Jan 2006"))
				}
				
				return errors.New(errorMsg)
			}

			// Get available units
			availableUnits, _ := u.getAvailableUnitsInDateRange(
				item.BarangID,
				loan.StartDate,
				loan.EndDate,
				&loan.ID,
			)

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
				LoanID:          loan.ID,
				FromFlow:        &oldFlow,
				ToFlow:          loan.LoanFlowStatus,
				TriggeredBy:     userID,
				TriggeredByRole: model.HistoryRoleUser,
				Note:            note,
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
				LoanID:          loan.ID,
				FromFlow:        &oldFlow,
				ToFlow:          loan.LoanFlowStatus,
				TriggeredBy:     userID,
				TriggeredByRole: model.HistoryRoleUser,
				Note:            note,
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

	if loan.LoanStatus != model.LoanStatusPending {
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
		ToStatus:      model.LoanStatusRejected,
		ChangedBy:     userID,
		ChangedByRole: model.HistoryRoleUser,
		Note:          stringPtr("Loan dibatalkan oleh user"),
	}); err != nil {
		return err
	}

	// ✅ RECORD HISTORY: Cancel loan - Flow
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          model.LoanFlowRequested, // tetap REQUESTED karena cancelled
		TriggeredBy:     userID,
		TriggeredByRole: model.HistoryRoleUser,
		Note:            "Loan dibatalkan oleh user",
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
// 🔥 NEW: CHECK AVAILABILITY WITH SUGGESTIONS
// =====================================================
//
func (u *loanUserUsecase) CheckAvailabilityWithSuggestions(
	kodeBarang string,
	quantity int,
	startDate, endDate time.Time,
) (*res.AvailabilityCheckResponse, error) {

	// 1. Lookup barang ID by kode
	barangID, err := u.barangRepo.GetIDByKode(kodeBarang)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	// 2. Check availability
	return u.checkAvailabilityWithGapRule(barangID, quantity, startDate, endDate, nil)
}

//
// =====================================================
// HELPER FUNCTIONS
// =====================================================
//

// ✅ Get available units in date range
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

// 🔥 NEW: Check availability dengan GAP 1 HARI RULE
func (u *loanUserUsecase) checkAvailabilityWithGapRule(
	barangID uint,
	requestedQty int,
	startDate time.Time,
	endDate time.Time,
	excludeLoanID *uint,
) (*res.AvailabilityCheckResponse, error) {

	// 1. Get available units NOW
	availableNow, err := u.getAvailableUnitsInDateRange(
		barangID,
		startDate,
		endDate,
		excludeLoanID,
	)

	if err != nil {
		return nil, err
	}

	currentlyAvailable := len(availableNow)

	// 2. Kalau available cukup → langsung OK
	if currentlyAvailable >= requestedQty {
		return &res.AvailabilityCheckResponse{
			IsAvailable:    true,
			AvailableCount: currentlyAvailable,
			RequestedCount: requestedQty,
			Message:        fmt.Sprintf("Tersedia %d unit", currentlyAvailable),
		}, nil
	}

	// 3. Kalau kurang → cari upcoming returns (gap 1 hari rule)
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil {
		return nil, err
	}

	// Get upcoming returns setelah startDate
	upcomingReturns, err := u.loanRepo.GetUpcomingReturns(barangID, startDate)
	if err != nil {
		return nil, err
	}

	// Cari tanggal earliest yang bisa kasih enough units
	needed := requestedQty - currentlyAvailable
	accumulated := 0
	var suggestedDate *time.Time

	for _, returnInfo := range upcomingReturns {
		accumulated += returnInfo.UnitCount
		
		if accumulated >= needed {
			// 🔥 GAP 1 HARI: endDate + 2 hari (1 hari gap)
			suggested := returnInfo.EndDate.AddDate(0, 0, 2)
			suggestedDate = &suggested
			break
		}
	}

	// 4. Build response dengan suggestions
	response := &res.AvailabilityCheckResponse{
		IsAvailable:    false,
		AvailableCount: currentlyAvailable,
		RequestedCount: requestedQty,
		Message: fmt.Sprintf(
			"Hanya %d unit tersedia (diminta %d) di tanggal %s - %s",
			currentlyAvailable,
			requestedQty,
			startDate.Format("02 Jan"),
			endDate.Format("02 Jan"),
		),
		SuggestedDate: suggestedDate,
	}

	// 5. 🔥 NEW: Cari barang alternatif (kategori sama, merk beda)
	similarBarangs, _ := u.barangRepo.FindSimilarBarang(
		barang.Kategori,
		barang.Merk,
		barangID,
		3, // max 3 alternatif
		barang.Nama, 
	)

	if len(similarBarangs) > 0 {
		alternatives := make([]res.AlternativeBarang, 0)
		
		for _, sb := range similarBarangs {
			// Check availability untuk barang alternatif ini
			altAvail, _ := u.getAvailableUnitsInDateRange(
				sb.ID,
				startDate,
				endDate,
				nil,
			)
			
			// Kalau ada yang available
			if len(altAvail) > 0 {
				alternatives = append(alternatives, res.AlternativeBarang{
					BarangID:       sb.ID,
					BarangKode:     sb.Kode,
					BarangNama:     sb.Nama,
					Merk:           sb.Merk,
					Kategori:       sb.Kategori,
					AvailableCount: len(altAvail),
				})
			}
		}
		
		response.AlternativeBarangs = alternatives
	}

	return response, nil
}

// ✅ Check if loan editable by user
func isEditable(flow string) bool {
	return flow == model.LoanFlowRequested || flow == model.LoanFlowRevision
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

//
// =====================================================
// 🔥 NEW: AUTO-REVISION HANDLER (dipanggil dari BarangUsecase)
// =====================================================
//

// MoveLoansToRevision - dipanggil ketika barang di-nonaktifkan
func (u *loanUserUsecase) MoveLoansToRevision(barangID uint, adminID uint) error {
	// 1. Get all active loans yang punya barang ini
	activeLoans, err := u.loanRepo.GetActiveLoansByBarangID(barangID)
	if err != nil {
		return err
	}

	if len(activeLoans) == 0 {
		return nil // Gak ada loan aktif
	}

	// 2. Get barang info untuk history notes
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil {
		return err
	}

	// 3. Collect loan IDs
	loanIDs := make([]uint, 0, len(activeLoans))
	for _, loan := range activeLoans {
		loanIDs = append(loanIDs, loan.ID)
	}

	// 4. Bulk update flow status ke REVISION
	if err := u.loanRepo.BulkUpdateLoanFlow(loanIDs, model.LoanFlowRevision); err != nil {
		return err
	}

	// 5. Create flow history untuk semua loans
	histories := make([]model.LoanFlowHistory, 0, len(activeLoans))
	
	for _, loan := range activeLoans {
		oldFlow := loan.LoanFlowStatus
		
		histories = append(histories, model.LoanFlowHistory{
			LoanID:          loan.ID,
			FromFlow:        &oldFlow,
			ToFlow:          model.LoanFlowRevision,
			TriggeredBy:     adminID,
			TriggeredByRole: model.HistoryRoleAdmin,
			Note:            fmt.Sprintf(
				"Loan masuk revisi karena barang '%s' dinonaktifkan oleh admin",
				barang.Nama,
			),
		})
	}

	// Bulk insert histories
	if err := u.flowHistRepo.BulkCreate(histories); err != nil {
		return err
	}

	return nil
}