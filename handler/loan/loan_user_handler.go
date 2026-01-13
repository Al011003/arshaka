package handler

import (
	req "backend/dto/request/loan"
	usecase "backend/usecase/loan"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type LoanUserHandler struct {
	loanUsecase usecase.LoanUserUsecase
}

func NewLoanUserHandler(loanUC usecase.LoanUserUsecase) *LoanUserHandler {
	return &LoanUserHandler{
		loanUsecase: loanUC,
	}
}

// POST /api/user/loans
func (h *LoanUserHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	var reqBody req.CreateLoanRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.loanUsecase.Create(reqBody, userID.(uint))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, result, "peminjaman berhasil dibuat")
}

// PUT /api/user/loans/:loan_id
func (h *LoanUserHandler) UpdateHeader(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}


	var reqBody req.UpdateLoanRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.loanUsecase.UpdateHeader(loanCode, userID, reqBody); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "peminjaman berhasil diperbarui")
}

// PUT /api/user/loans/:loan_id/items
func (h *LoanUserHandler) UpdateItems(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	var reqBody req.UpdateLoanItemsRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.loanUsecase.UpdateItems(loanCode, userID, reqBody); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang peminjaman berhasil diperbarui")
}

// GET /api/user/loans
func (h *LoanUserHandler) GetMyLoans(c *gin.Context) {
	userID := c.GetUint("user_id")

	loans, err := h.loanUsecase.GetMyLoans(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, loans, "daftar peminjaman")
}

// DELETE /api/user/loans/:loan_id
func (h *LoanUserHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	if err := h.loanUsecase.Cancel(loanCode, userID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "peminjaman berhasil dibatalkan")
}

// GET /api/user/loans/:loan_code
func (h *LoanUserHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	result, err := h.loanUsecase.GetDetail(loanCode, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, result, "detail peminjaman")
}
