package handler

import (
	req "backend/dto/request/admin"
	usecase "backend/usecase/admin/reset_password"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type AdminForgotPasswordHandler struct {
	uc usecase.AdminResetPasswordUsecase
}

func NewAdminForgotPasswordHandler(uc usecase.AdminResetPasswordUsecase) *AdminForgotPasswordHandler {
	return &AdminForgotPasswordHandler{uc: uc}
}
func (h *AdminForgotPasswordHandler) RequestForgotPassword(c *gin.Context) {
	var body req.ForgotPasswordRequest

	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.uc.RequestResetByIdentity(body.NRA, body.Name); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "Permintaan reset password berhasil dibuat")
}
