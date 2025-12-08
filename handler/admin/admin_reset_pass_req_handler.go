package handler

import (
	res "backend/dto/response/common" // kalau kamu simpan BaseResponse di utils
	usecase "backend/usecase/admin/reset_password"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminForgotPasswordHandler struct {
    uc usecase.AdminResetPasswordUsecase
}

func NewUserForgotPasswordHandler(uc usecase.AdminResetPasswordUsecase) *AdminForgotPasswordHandler {
    return &AdminForgotPasswordHandler{uc: uc}
}

func (h *AdminForgotPasswordHandler) RequestForgotPassword(c *gin.Context) {
    userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, res.BaseResponse{
			Status:  "error",
			Message: "unauthorized",
			Data:    nil,
		})
		return
	}

	userID := userIDValue.(uint)

    err := h.uc.RequestReset(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status: "error",
            Message: err.Error(),
            Data:    nil,
        })
        return
    }

    c.JSON(http.StatusOK, res.BaseResponse{
		Status: "success",
        Message: "Permintaan reset password berhasil dibuat",
        Data:    nil,
    })
}

// ======================================================
// 2. User Cancel Forgot Password
// POST /user/forgot-password/cancel
// ======================================================
func (h *AdminForgotPasswordHandler) CancelForgotPassword(c *gin.Context) {
    userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, res.BaseResponse{
			Status:  "error",
			Message: "unauthorized",
			Data:    nil,
		})
		return
	}

	userID := userIDValue.(uint)

    err := h.uc.CancelReset(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status: "error",
            Message: err.Error(),
            Data:    nil,
        })
        return
    }

    c.JSON(http.StatusOK, res.BaseResponse{
		Status: "success",
        Message: "Permintaan reset password berhasil dibatalkan",
        Data:    nil,
    })
}
