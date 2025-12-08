package handler

import (
	usecase "backend/usecase/super_admin/getprofile"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type SuperAdminProfileHandler struct {
	profileUC usecase.SuperAdminProfileUsecase
}

func NewSuperAdminProfileHandler(profileUC usecase.SuperAdminProfileUsecase) *SuperAdminProfileHandler {
	return &SuperAdminProfileHandler{
		profileUC: profileUC,
	}
}

func (h *SuperAdminProfileHandler) GetProfile(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	superAdminID := idRaw.(uint)

	profile, err := h.profileUC.GetProfile(superAdminID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, profile, "superadmin profile retrieved")
}
