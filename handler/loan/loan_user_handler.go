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

// CreateLoan godoc
// @Summary      Create loan
// @Description  User membuat pengajuan peminjaman barang
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateLoanRequest  true  "Create loan request"
// @Success      201      {object}  map[string]interface{}  "Loan berhasil dibuat"
// @Failure      400      {object}  map[string]interface{}  "Bad request"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/user/loan [post]
// @Security     ApiKeyAuth
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

// UpdateLoanHeader godoc
// @Summary      Update loan header
// @Description  Update data header peminjaman (tanggal, catatan, dll)
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        loan_code  path      string                  true  "Loan code"
// @Param        request    body      request.UpdateLoanRequest   true  "Update loan header"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /api/user/loan/{loan_code} [put]
// @Security     ApiKeyAuth
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

// UpdateLoanItems godoc
// @Summary      Update loan items
// @Description  Update daftar barang dalam peminjaman
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        loan_code  path      string                        true  "Loan code"
// @Param        request    body      request.UpdateLoanItemsRequest    true  "Update loan items"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /api/user/loan/{loan_code}/items [put]
// @Security     ApiKeyAuth
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

// GetMyLoans godoc
// @Summary      Get my loans
// @Description  Mengambil daftar peminjaman milik user
// @Tags         user-loan
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/user/loan [get]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) GetMyLoans(c *gin.Context) {
	userID := c.GetUint("user_id")

	loans, err := h.loanUsecase.GetMyLoans(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, loans, "daftar peminjaman")
}

// CancelLoan godoc
// @Summary      Cancel loan
// @Description  Membatalkan peminjaman (selama masih pending)
// @Tags         user-loan
// @Produce      json
// @Param        loan_code  path      string  true  "Loan code"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /api/user/loan/{loan_code} [delete]
// @Security     ApiKeyAuth
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

// GetLoanDetail godoc
// @Summary      Get loan detail
// @Description  Mengambil detail peminjaman berdasarkan loan code
// @Tags         user-loan
// @Produce      json
// @Param        loan_code  path      string  true  "Loan code"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /api/user/loan/{loan_code} [get]
// @Security     ApiKeyAuth
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
