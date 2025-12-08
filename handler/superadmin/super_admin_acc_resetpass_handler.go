package handler

import (
	"strconv"

	usecase "backend/usecase/super_admin/reset_password"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type SuperAdminResetPasswordHandler struct {
	uc usecase.SuperAdminAccResetPasswordUsecase
}

func NewAdminResetPasswordHandler(uc usecase.SuperAdminAccResetPasswordUsecase) *SuperAdminResetPasswordHandler {
	return &SuperAdminResetPasswordHandler{uc: uc}
}

// ======================================================
// POST /admin/reset-password/approve/:userID
// ======================================================
func (h *SuperAdminResetPasswordHandler) ApproveReset(c *gin.Context) {
	// --------------------------------------------------
	// 1. Ambil adminID dari JWT
	// --------------------------------------------------
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	AdminID := idRaw.(uint)

	// --------------------------------------------------
	// 2. Ambil userID dari param
	// --------------------------------------------------
	resetIDParam := c.Param("resetID")
	resetIDUint64, err := strconv.ParseUint(resetIDParam, 10, 64)
	if err != nil {
		utils.BadRequest(c, "resetID tidak valid")
		return
	}

	resetID := uint(resetIDUint64)

	// --------------------------------------------------
	// 3. Panggil usecase
	// --------------------------------------------------
	err = h.uc.Approve(resetID, AdminID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// --------------------------------------------------
	// 4. Response sukses
	// --------------------------------------------------
	utils.Success(c, nil, "Reset password berhasil disetujui. Password user telah direset.")
}

func (h *SuperAdminResetPasswordHandler) GetAllRequests(c *gin.Context) {
	// AdminID check (opsional jika sudah dicek middleware)
	status := c.Query("status")
	// Ambil list request
	reqs, err := h.uc.GetAllRequestsFiltered(status)
	if err != nil {
		utils.InternalError(c, "gagal mengambil data request reset")
		return
	}

	utils.Success(c, reqs, "berhasil mengambil data request reset")
}

