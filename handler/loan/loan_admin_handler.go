// handler/loan_admin_handler.go
package handler

import (
	req "backend/dto/request/loan"
	usecase "backend/usecase/loan"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type LoanAdminHandler struct {
	loanAdminUC usecase.LoanAdminUsecase
}

func NewLoanAdminHandler(uc usecase.LoanAdminUsecase) *LoanAdminHandler {
	return &LoanAdminHandler{
		loanAdminUC: uc,
	}
}

// ApproveLoan - Approve loan
func (h *LoanAdminHandler) ApproveLoan(c *gin.Context) {
	adminID := c.GetUint("user_id")
	loanCode := c.Param("loan_code")

	if err := h.loanAdminUC.ApproveLoan(loanCode, adminID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "loan berhasil disetujui")
}

// RejectLoan - Reject loan
func (h *LoanAdminHandler) RejectLoan(c *gin.Context) {
	adminID := c.GetUint("user_id")
	loanCode := c.Param("loan_code")

	var reqBody req.RejectLoanRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.loanAdminUC.RejectLoan(loanCode, adminID, reqBody.Reason); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "loan berhasil ditolak")
}

// RequestRevision - Minta user revisi loan
func (h *LoanAdminHandler) RequestRevision(c *gin.Context) {
	adminID := c.GetUint("user_id")
	loanCode := c.Param("loan_code")

	var reqBody req.RequestRevisionRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.loanAdminUC.RequestRevision(loanCode, adminID, reqBody.Note); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "revisi berhasil diminta")
}

// GetAllLoans - Get all loans (admin view)
func (h *LoanAdminHandler) GetAllLoans(c *gin.Context) {
	var filter req.AdminLoanFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	loans, err := h.loanAdminUC.GetAllLoans(filter)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, loans, "daftar loan")
}

// GetLoanDetail - Get loan detail (admin view)
func (h *LoanAdminHandler) GetLoanDetail(c *gin.Context) {
	loanCode := c.Param("loan_code")

	result, err := h.loanAdminUC.GetLoanDetail(loanCode)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, result, "detail loan")
}

// MarkAsReady - Tandai loan siap diambil
func (h *LoanAdminHandler) MarkAsReady(c *gin.Context) {
	adminID := c.GetUint("user_id")
	loanCode := c.Param("loan_code")

	if err := h.loanAdminUC.MarkAsReady(loanCode, adminID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "loan ditandai siap diambil")
}

// MarkAsTaken - Tandai barang sudah diambil
func (h *LoanAdminHandler) MarkAsTaken(c *gin.Context) {
	adminID := c.GetUint("user_id")
	loanCode := c.Param("loan_code")

	if err := h.loanAdminUC.MarkAsTaken(loanCode, adminID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang ditandai sudah diambil")
}

// MarkAsReturned - Tandai barang sudah dikembalikan
func (h *LoanAdminHandler) MarkAsReturned(c *gin.Context) {
	adminID := c.GetUint("user_id")
	loanCode := c.Param("loan_code")

	if err := h.loanAdminUC.MarkAsReturned(loanCode, adminID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang ditandai sudah dikembalikan")
}