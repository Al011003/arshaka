package handler

import (
	"net/http"

	req "backend/dto/request/auth"
	usecase "backend/usecase/auth"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	authUsecase usecase.RegisterUsecase
}

func NewRegisterHandler(a usecase.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		authUsecase: a,
	}
}

// =================== Register User / Anggota ===================
func (h *RegisterHandler) RegisterUser(c *gin.Context) {
	var request req.RegisterUserRequest
	if err := request.BindandValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.authUsecase.RegisterUser(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	utils.Success(c, resp, "user Berhasil dibuat")

}

// =================== Register Admin ===================
func (h *RegisterHandler) RegisterAdmin(c *gin.Context) {
	var request req.RegisterAdminRequest
	// Validasi payload


	if err := request.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.authUsecase.RegisterAdmin(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	utils.Success(c, resp, "admin berhasil dibuat")
}
