package handler

import (
	"github.com/gin-gonic/gin"

	usecase "backend/usecase/admin/profile"
	"backend/utils"
)

type AdminProfileHandler struct {
	uc usecase.AdminProfileUsecase
}

func NewAdminProfileHandler(uc usecase.AdminProfileUsecase) *AdminProfileHandler {
	return &AdminProfileHandler{
		uc: uc,
	}
}

// GetProfile godoc
// @Summary      Get admin profile
// @Description  Mengambil data profile admin yang sedang login
// @Tags         admin-profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Admin profile"
// @Failure      401  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/admin/profile [get]
// @Security     ApiKeyAuth
func (h *AdminProfileHandler) GetProfile(c *gin.Context) {
	// Ambil admin ID dari JWT
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	adminID := idRaw.(uint)

	// Call usecase
	profile, err := h.uc.GetProfile(adminID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, profile, "admin profile retrieved")
}
