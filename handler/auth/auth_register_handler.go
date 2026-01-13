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

// RegisterUser godoc
// @Summary      Register new user/anggota
// @Description  Mendaftarkan user/anggota baru ke sistem
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.RegisterUserRequest  true  "User registration data"
// @Success      200      {object}  map[string]interface{}  "User berhasil dibuat"
// @Failure      400      {object}  map[string]interface{}  "Bad request - validation error"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized - token tidak valid atau tidak ada"
// @Router       /api/super-admin/register-user [post]
// @Security     ApiKeyAuth
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

// RegisterAdmin godoc
// @Summary      Register new admin
// @Description  Mendaftarkan admin baru (hanya bisa dilakukan oleh super admin)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.RegisterAdminRequest  true  "Admin registration data"
// @Success      200      {object}  map[string]interface{}  "Admin berhasil dibuat"
// @Failure      400      {object}  map[string]interface{}  "Bad request - validation error"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized - token tidak valid atau tidak ada"
// @Router       /api/super-admin/register-admin [post]
// @Security ApiKeyAuth
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
