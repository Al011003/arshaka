package handler

import (
	usecase "backend/usecase/user/profile"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	profileUC usecase.UserProfileUsecase
}

func NewUserProfileHandler(profileUC usecase.UserProfileUsecase) *UserProfileHandler {
	return &UserProfileHandler{
		profileUC: profileUC,
	}
}

func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	userID := idRaw.(uint)

	profile, err := h.profileUC.GetProfile(userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}


	utils.Success(c, profile, "user profile retrieved")
}
