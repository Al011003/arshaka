// handler/loan_handler.go
package handler

import (
	loan "backend/dto/request/loan"
	usecase "backend/usecase/loan"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type LoanHandler struct {
	loanUsecase usecase.LoanUsecase
}

func NewLoanHandler(loanUC usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{
		loanUsecase: loanUC,
	}
}

// CreateLoanFromCart - POST /api/loans/create-from-cart
func (h *LoanHandler) CreateLoanFromCart(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Validate request
	var req loan.CreateLoanRequest
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Create loan from cart
	result, err := h.loanUsecase.CreateLoanFromCart(userID.(uint), req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result, "peminjaman berhasil dibuat")
}