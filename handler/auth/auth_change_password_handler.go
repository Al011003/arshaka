package handler

import (
	req "backend/dto/request/auth"
	usecase "backend/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdatePasswordHandler struct {
    uc usecase.UpdatePasswordUsecase
}

func NewUpdatePasswordHandler(uc usecase.UpdatePasswordUsecase) *UpdatePasswordHandler {
    return &UpdatePasswordHandler{
        uc: uc,
    }
}

// UpdatePassword godoc
// @Summary      Update user password
// @Description  Update password untuk user yang sedang login (memerlukan old password)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.UpdatePasswordRequest  true  "Old password and new password"
// @Success      200      {object}  map[string]interface{}  "Password berhasil diperbarui"
// @Failure      400      {object}  map[string]interface{}  "Bad request - validation error atau old password salah"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized - user tidak ditemukan dalam token"
// @Router       /api/user/password [post]
// @Security ApiKeyAuth
func (h *UpdatePasswordHandler) UpdatePassword(c *gin.Context) {
    var request req.UpdatePasswordRequest

    // Bind + validate request
    if err := request.BindAndValidate(c); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Ambil userID dari JWT (contoh: sudah di-set di middleware)
    userIDValue, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "user tidak ditemukan dalam token",
        })
        return
    }
    userID := userIDValue.(uint)

    // Call usecase
    if err := h.uc.UpdatePassword(c.Request.Context(), userID, request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Success response
    c.JSON(http.StatusOK, gin.H{
        "message": "password berhasil diperbarui",
    })
}