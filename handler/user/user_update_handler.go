package handler

import (
	req "backend/dto/request/user"
	usecase "backend/usecase/user/update"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type UserUpdateHandler struct {
	selfUC usecase.UserSelfUsecase
}

func NewUserUpdateHandler(u usecase.UserSelfUsecase) *UserUpdateHandler {
	return &UserUpdateHandler{
		selfUC: u,
	}
}

func (h *UserUpdateHandler) UpdateSelf(c *gin.Context) {

	// Ambil userID dari JWT
	userIDAny, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	userID := userIDAny.(uint)

	// Bind + validate
	var reqData req.UserUpdateRequest
	if err := reqData.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	// Usecase
	if err := h.selfUC.UpdateSelf(userID, reqData); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "profil berhasil diperbarui")
}
