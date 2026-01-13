package handler

import (
	usecase "backend/usecase/loan"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type LoanUserGetHandler struct {
	loanGetUsecase usecase.LoanGetUsecase
}

func NewLoanUserGetHandler(loanGetUC usecase.LoanGetUsecase) *LoanUserGetHandler {
	return &LoanUserGetHandler{
		loanGetUsecase: loanGetUC,
	}
}

//
// =====================================================
// GET /api/loans
// List loan user (card / table)
// =====================================================
//
func (h *LoanUserGetHandler) GetMyLoans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	result, err := h.loanGetUsecase.GetMyLoans(userID.(uint))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "berhasil mengambil data peminjaman")
}

//
// =====================================================
// GET /api/loans/:loan_code
// Detail loan by code
// =====================================================
//
func (h *LoanUserGetHandler) GetMyLoanDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan code wajib diisi")
		return
	}

	result, err := h.loanGetUsecase.GetLoanDetailByCode(
		userID.(uint),
		loanCode,
	)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "berhasil mengambil detail peminjaman")
}
