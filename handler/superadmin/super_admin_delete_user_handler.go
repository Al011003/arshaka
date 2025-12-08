package handler

import (
	response "backend/dto/response/common"
	usecase "backend/usecase/super_admin/delete"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuperAdminDeleteUserHandler struct {
	UC usecase.SuperAdminDeleteUserUsecase
}

// Constructor
func NewSuperAdminDeleteUserHandler(uc usecase.SuperAdminDeleteUserUsecase) *SuperAdminDeleteUserHandler {
	return &SuperAdminDeleteUserHandler{UC: uc}
}

// DELETE /api/super-admin/user/:id
func (h *SuperAdminDeleteUserHandler) DeleteUser(c *gin.Context) {

	// Ambil superadmin ID dari JWT
	idRaw, ok := c.Get("user_id")
	if !ok {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	superAdminID := idRaw.(uint)

	// Ambil target user ID dari param
	paramID := c.Param("id")
	targetID64, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}
	targetID := uint(targetID64)

	// Eksekusi Usecase
	if err := h.UC.DeleteUser(superAdminID, targetID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}
