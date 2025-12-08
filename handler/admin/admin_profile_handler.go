package handler

import (
	"net/http"

	res "backend/dto/response/common"
	usecase "backend/usecase/admin/profile"

	"github.com/gin-gonic/gin"
)

type AdminProfileHandler struct {
	profileUC usecase.AdminProfileUsecase
}

func NewAdminProfileHandler(profileUC usecase.AdminProfileUsecase) *AdminProfileHandler {
	return &AdminProfileHandler{
		profileUC: profileUC,
	}
}

func (h *AdminProfileHandler) GetProfile(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, res.BaseResponse{
			Status:  "error",
			Message: "unauthorized",
		})
		return
	}

	AdminID := idRaw.(uint)

	profile, err := h.profileUC.GetProfile(AdminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "admin profile retrieved",
		Data:    profile,
	})
}
