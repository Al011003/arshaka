package handler

import (
	"net/http"

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

// POST /auth/request-otp
func (h *PasswordResetHandler) RequestOTP(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required,email"`
    }

    // Ambil email dari body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "email tidak valid",
        })
        return
    }

    // Kirim OTP
    if err := h.authUC.SendOTP(c, req.Email); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "OTP telah dikirim ke email",
    })
}

func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
    var req struct {
        Email       string `json:"email" binding:"required,email"`
        OTP         string `json:"otp" binding:"required,len=6"`
        NewPassword string `json:"new_password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "input tidak valid",
        })
        return
    }

    // Reset password
    if err := h.authUC.ResetPassword(c, req.Email, req.OTP, req.NewPassword); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": err.Error(),
        })
        return
    }

    utils.Success(c, nil, "Password berhasil diganti")
}