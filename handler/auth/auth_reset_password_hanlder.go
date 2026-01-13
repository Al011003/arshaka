package handler

import (
	"net/http"

	req "backend/dto/request/auth"
	res "backend/dto/response/common"
	usecase "backend/usecase/auth"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	authUC usecase.PasswordResetUsecase
}

func NewPasswordResetHandler(authUC usecase.PasswordResetUsecase) *PasswordResetHandler {
	return &PasswordResetHandler{
		authUC: authUC,
	}
}

// RequestOTP godoc
// @Summary      Request OTP reset password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.RequestOTPRequest  true  "Email user"
// @Success      200      {object}  response.BaseResponse
// @Failure      400      {object}  response.BaseResponse
// @Router       /auth/request-otp [post]
func (h *PasswordResetHandler) RequestOTP(c *gin.Context) {
	var request req.RequestOTPRequest

	if err := request.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if err := h.authUC.SendOTP(c, request.Email); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	utils.Success(c, nil, "OTP berhasil dikirim ke email")
}

// ======================================================

// VerifyResetOTP godoc
// @Summary      Verifikasi OTP reset password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.VerifyResetOTPRequest  true  "Email dan OTP"
// @Success      200      {object}  response.BaseResponse
// @Failure      400      {object}  response.BaseResponse
// @Router       /auth/verify-reset-otp [post]
func (h *PasswordResetHandler) VerifyResetOTP(c *gin.Context) {
	var request req.VerifyResetOTPRequest

	if err := request.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	resetToken, expiresIn, err := h.authUC.VerifyOTPAndGenerateToken(
		c,
		request.Email,
		request.OTP,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	utils.Success(c, gin.H{
		"reset_token": resetToken,
		"expires_in":  expiresIn,
	}, "OTP berhasil diverifikasi")
}

// ======================================================

// ResetPassword godoc
// @Summary      Reset password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.ResetPasswordRequest  true  "Reset token dan password baru"
// @Success      200      {object}  response.BaseResponse
// @Failure      400      {object}  response.BaseResponse
// @Router       /auth/reset-password [post]
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var request req.ResetPasswordRequest

	if err := request.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if err := h.authUC.ResetPasswordWithToken(
		c,
		request.ResetToken,
		request.NewPassword,
	); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	utils.Success(c, nil, "Password berhasil diganti")
}
