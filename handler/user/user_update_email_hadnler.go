package handler

import (
	req "backend/dto/request/user"
	usecase "backend/usecase/user/update"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type UpdateEmailHandler struct {
	emailUC usecase.UpdateEmailUsecase
}

func NewUpdateEmailHandler(u usecase.UpdateEmailUsecase) *UpdateEmailHandler {
	return &UpdateEmailHandler{
		emailUC: u,
	}
}

func (h *UpdateEmailHandler) UpdateEmail(c *gin.Context) {
	var request req.UpdateEmailRequest

	// Bind + Validasi DTO
	if err := request.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Ambil userID dari JWT
	userIDAny, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	userID := userIDAny.(uint)

	// Eksekusi usecase
	if err := h.emailUC.UpdateEmail(c, userID, request); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Success
	utils.Success(c, nil, "email berhasil diperbarui")
}
