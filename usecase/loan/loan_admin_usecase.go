// usecase/loan_admin_usecase.go
package usecase

import (
	req "backend/dto/request/loan"
	res "backend/dto/response/loan"
	"backend/model"
	repository "backend/repo"
	"errors"
	"fmt"
)

type LoanAdminUsecase interface {
	// Approval & Rejection
	ApproveLoan(loanCode string, adminID uint) error
	RejectLoan(loanCode string, adminID uint, reason string) error
	RequestRevision(loanCode string, adminID uint, note string) error

	// Get Loans (Admin View)
	GetAllLoans(filter req.AdminLoanFilter) ([]res.LoanListResponse, error)
	GetLoanDetail(loanCode string) (*res.LoanDetailResponse, error)

	// Flow Management
	MarkAsReady(loanCode string, adminID uint) error
	MarkAsTaken(loanCode string, adminID uint) error
	MarkAsReturned(loanCode string, adminID uint) error

	// Auto-Revision (dipanggil dari BarangUsecase)
	MoveLoansToRevision(barangID uint, adminID uint) error
}

type loanAdminUsecase struct {
	loanRepo       repository.LoanRepository
	barangRepo     repository.BarangRepository
	unitRepo       repository.BarangUnitRepository
	statusHistRepo repository.LoanStatusHistoryRepository
	flowHistRepo   repository.LoanFlowHistoryRepository
}

func NewLoanAdminUsecase(
	loanRepo repository.LoanRepository,
	barangRepo repository.BarangRepository,
	unitRepo repository.BarangUnitRepository,
	statusHistRepo repository.LoanStatusHistoryRepository,
	flowHistRepo repository.LoanFlowHistoryRepository,
) LoanAdminUsecase {
	return &loanAdminUsecase{
		loanRepo:       loanRepo,
		barangRepo:     barangRepo,
		unitRepo:       unitRepo,
		statusHistRepo: statusHistRepo,
		flowHistRepo:   flowHistRepo,
	}
}

//
// =====================================================
// APPROVE LOAN
// =====================================================
//
func (u *loanAdminUsecase) ApproveLoan(
	loanCode string,
	adminID uint,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.LoanStatus != model.LoanStatusPending {
		return errors.New("hanya loan dengan status PENDING yang bisa disetujui")
	}

	// 🔒 VALIDASI: pastikan semua unit masih available
	for _, item := range loan.LoanItems {
		for _, unit := range item.AssignedUnits {
			if unit.Status != model.UnitStatusSiap {
				return fmt.Errorf("unit %s sudah tidak tersedia", unit.KodeUnit)
			}
		}
	}

	oldStatus := loan.LoanStatus
	oldFlow := loan.LoanFlowStatus

	loan.LoanStatus = model.LoanStatusApproved
	loan.LoanFlowStatus = model.LoanFlowReady

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	// update item status
	for _, item := range loan.LoanItems {
		item.Status = model.LoanItemStatusReady
		_ = u.loanRepo.UpdateLoanItem(&item)
	}

	// Record history - Status
	if err := u.statusHistRepo.Create(&model.LoanStatusHistory{
		LoanID:        loan.ID,
		FromStatus:    &oldStatus,
		ToStatus:      model.LoanStatusApproved,
		ChangedBy:     adminID,
		ChangedByRole: model.HistoryRoleAdmin,
		Note:          stringPtr("Loan disetujui oleh admin"),
	}); err != nil {
		return err
	}

	// Record history - Flow
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          model.LoanFlowReady,
		TriggeredBy:     adminID,
		TriggeredByRole: model.HistoryRoleAdmin,
		Note:            "Loan disetujui, barang siap diambil",
	}); err != nil {
		return err
	}

	// Update loan items status
	for _, item := range loan.LoanItems {
		item.Status = model.LoanItemStatusReady
		_ = u.loanRepo.UpdateLoanItem(&item)
	}

	return nil
}

//
// =====================================================
// REJECT LOAN
// =====================================================
//
func (u *loanAdminUsecase) RejectLoan(
	loanCode string,
	adminID uint,
	reason string,
) error {

	if reason == "" {
		return errors.New("alasan penolakan harus diisi")
	}

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	// Validasi: hanya PENDING yang bisa di-reject
	if loan.LoanStatus != model.LoanStatusPending {
		return errors.New("hanya loan dengan status PENDING yang bisa ditolak")
	}

	// Store old values
	oldStatus := loan.LoanStatus
	oldFlow := loan.LoanFlowStatus

	// Update status
	loan.LoanStatus = model.LoanStatusRejected
	loan.RejectReason = &reason

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	// Remove all unit assignments
	for _, item := range loan.LoanItems {
		_ = u.loanRepo.RemoveUnitsFromLoanItem(item.ID)
		
		// Update item status
		item.Status = model.LoanItemStatusRejected
		_ = u.loanRepo.UpdateLoanItem(&item)
	}

	// Record history - Status
	if err := u.statusHistRepo.Create(&model.LoanStatusHistory{
		LoanID:        loan.ID,
		FromStatus:    &oldStatus,
		ToStatus:      model.LoanStatusRejected,
		ChangedBy:     adminID,
		ChangedByRole: model.HistoryRoleAdmin,
		Note:          stringPtr(fmt.Sprintf("Loan ditolak: %s", reason)),
	}); err != nil {
		return err
	}

	// Record history - Flow
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          loan.LoanFlowStatus,
		TriggeredBy:     adminID,
		TriggeredByRole: model.HistoryRoleAdmin,
		Note:            fmt.Sprintf("Loan ditolak: %s", reason),
	}); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// REQUEST REVISION
// =====================================================
//
func (u *loanAdminUsecase) RequestRevision(
	loanCode string,
	adminID uint,
	note string,
) error {

	if note == "" {
		return errors.New("catatan revisi harus diisi")
	}

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	// Validasi: hanya PENDING yang bisa diminta revisi
	if loan.LoanStatus != model.LoanStatusPending {
		return errors.New("hanya loan dengan status PENDING yang bisa diminta revisi")
	}

	// Store old flow
	oldFlow := loan.LoanFlowStatus

	// Update flow status
	loan.LoanFlowStatus = model.LoanFlowRevision

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	// Record history - Flow
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          model.LoanFlowRevision,
		TriggeredBy:     adminID,
		TriggeredByRole: model.HistoryRoleAdmin,
		Note:            fmt.Sprintf("Admin meminta revisi: %s", note),
	}); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// GET ALL LOANS (ADMIN)
// =====================================================
//
func (u *loanAdminUsecase) GetAllLoans(
    filter req.AdminLoanFilter,
) ([]res.LoanListResponse, error) {

    loans, err := u.loanRepo.GetAllLoansWithFilter(filter)
    if err != nil {
        return nil, err
    }

    resp := make([]res.LoanListResponse, 0, len(loans))

    for _, loan := range loans {
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
// GET LOAN DETAIL (ADMIN)
// =====================================================
//
func (u *loanAdminUsecase) GetLoanDetail(
	loanCode string,
) (*res.LoanDetailResponse, error) {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return nil, err
	}

	items := make([]res.LoanItemDetailResponse, 0)
	totalUnits := 0

	for _, item := range loan.LoanItems {
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
		RejectReason:   loan.RejectReason,
		TotalItems:     len(loan.LoanItems),
		TotalUnits:     totalUnits,
		Items:          items,
		CreatedAt:      loan.CreatedAt,
	}, nil
}

//
// =====================================================
// MARK AS READY
// =====================================================
//
func (u *loanAdminUsecase) MarkAsReady(
	loanCode string,
	adminID uint,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.LoanStatus != model.LoanStatusApproved {
		return errors.New("loan harus dalam status APPROVED")
	}

	if loan.LoanFlowStatus != model.LoanFlowRequested {
		return errors.New("loan sudah dalam status " + loan.LoanFlowStatus)
	}

	oldFlow := loan.LoanFlowStatus
	loan.LoanFlowStatus = model.LoanFlowReady

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	// Update items status
	for _, item := range loan.LoanItems {
		item.Status = model.LoanItemStatusReady
		_ = u.loanRepo.UpdateLoanItem(&item)
	}

	// Record history
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          model.LoanFlowReady,
		TriggeredBy:     adminID,
		TriggeredByRole: model.HistoryRoleAdmin,
		Note:            "Barang siap diambil",
	}); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// MARK AS TAKEN
// =====================================================
//
func (u *loanAdminUsecase) MarkAsTaken(
	loanCode string,
	adminID uint,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.LoanFlowStatus != model.LoanFlowReady {
		return errors.New("loan harus dalam status READY")
	}

	oldFlow := loan.LoanFlowStatus
	loan.LoanFlowStatus = model.LoanFlowTaken

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	for _, item := range loan.LoanItems {
		item.Status = model.LoanItemStatusBorrowed
		_ = u.loanRepo.UpdateLoanItem(&item)

		// 🔥 INI PENTING
		for _, unit := range item.AssignedUnits {
			unit.Status = model.UnitStatusDipinjam
			_ = u.unitRepo.Update(&unit)
		}
	}

	// Record history
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          model.LoanFlowTaken,
		TriggeredBy:     adminID,
		TriggeredByRole: model.HistoryRoleAdmin,
		Note:            "Barang telah diambil oleh user",
	}); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// MARK AS RETURNED
// =====================================================
//
func (u *loanAdminUsecase) MarkAsReturned(
	loanCode string,
	adminID uint,
) error {

	loan, err := u.loanRepo.GetLoanByCodeWithItems(loanCode)
	if err != nil {
		return err
	}

	if loan.LoanFlowStatus != model.LoanFlowTaken {
		return errors.New("loan harus dalam status TAKEN")
	}

	oldFlow := loan.LoanFlowStatus
	loan.LoanFlowStatus = model.LoanFlowFinished

	if err := u.loanRepo.UpdateLoan(loan); err != nil {
		return err
	}

	for _, item := range loan.LoanItems {
		item.Status = model.LoanItemStatusReturned
		_ = u.loanRepo.UpdateLoanItem(&item)

		for _, unit := range item.AssignedUnits {
			unit.Status = model.UnitStatusSiap
			_ = u.unitRepo.Update(&unit)
		}

		_ = u.loanRepo.RemoveUnitsFromLoanItem(item.ID)
	}
	// Record history
	if err := u.flowHistRepo.Create(&model.LoanFlowHistory{
		LoanID:          loan.ID,
		FromFlow:        &oldFlow,
		ToFlow:          model.LoanFlowFinished,
		TriggeredBy:     adminID,
		TriggeredByRole: model.HistoryRoleAdmin,
		Note:            "Barang telah dikembalikan",
	}); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// AUTO-REVISION (dipanggil dari BarangUsecase)
// =====================================================
//
func (u *loanAdminUsecase) MoveLoansToRevision(
	barangID uint,
	adminID uint,
) error {

	// Get all active loans yang punya barang ini
	activeLoans, err := u.loanRepo.GetActiveLoansByBarangID(barangID)
	if err != nil {
		return err
	}

	if len(activeLoans) == 0 {
		return nil
	}

	// Get barang info
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil {
		return err
	}

	// Collect loan IDs
	loanIDs := make([]uint, 0, len(activeLoans))
	for _, loan := range activeLoans {
		loanIDs = append(loanIDs, loan.ID)
	}

	// Bulk update flow status
	if err := u.loanRepo.BulkUpdateLoanFlow(loanIDs, model.LoanFlowRevision); err != nil {
		return err
	}

	// Create flow history
	histories := make([]model.LoanFlowHistory, 0, len(activeLoans))

	for _, loan := range activeLoans {
		oldFlow := loan.LoanFlowStatus

		histories = append(histories, model.LoanFlowHistory{
			LoanID:          loan.ID,
			FromFlow:        &oldFlow,
			ToFlow:          model.LoanFlowRevision,
			TriggeredBy:     adminID,
			TriggeredByRole: model.HistoryRoleAdmin,
			Note: fmt.Sprintf(
				"Loan masuk revisi karena barang '%s' dinonaktifkan oleh admin",
				barang.Nama,
			),
		})
	}

	if err := u.flowHistRepo.BulkCreate(histories); err != nil {
		return err
	}

	return nil
}

//
// =====================================================
// HELPER FUNCTIONS
// =====================================================
//

func stringPtr(s string) *string {
	return &s
}