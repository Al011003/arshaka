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

// ApproveReset godoc
// @Summary      Approve reset password request
// @Description  Menyetujui request reset password dan mereset password user (super admin only)
// @Tags         super-admin-reset-password
// @Produce      json
// @Param        resetID  path  int  true  "Reset Request ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/super-admin/reset-password/approve/{resetID} [post]
// @Security     ApiKeyAuth
func (h *SuperAdminResetPasswordHandler) ApproveReset(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	AdminID := idRaw.(uint)
	resetIDParam := c.Param("resetID")
	resetIDUint64, err := strconv.ParseUint(resetIDParam, 10, 64)
	if err != nil {
		utils.BadRequest(c, "resetID tidak valid")
		return
	}

	resetID := uint(resetIDUint64)
	err = h.uc.Approve(resetID, AdminID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, nil, "Reset password berhasil disetujui. Password user telah direset.")
}


// CancelReset godoc
// @Summary      Cancel reset password request
// @Description  Membatalkan request reset password (super admin only)
// @Tags         super-admin-reset-password
// @Produce      json
// @Param        resetID  path  int  true  "Reset Request ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/super-admin/reset-password/cancel/{resetID} [post]
// @Security     ApiKeyAuth
func (h *SuperAdminResetPasswordHandler)CancelReset(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	AdminID := idRaw.(uint)
	resetIDParam := c.Param("resetID")
	resetIDUint64, err := strconv.ParseUint(resetIDParam, 10, 64)
	if err != nil {
		utils.BadRequest(c, "resetID tidak valid")
		return
	}

	resetID := uint(resetIDUint64)
	err = h.uc.CancelResetByIdentity(resetID, AdminID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, nil, "Reset password tidak disetujui. Password user tidak berubah.")
}

// GetAllRequests godoc
// @Summary      Get all reset password requests
// @Description  Mengambil semua request reset password dengan filter status (super admin only)
// @Tags         super-admin-reset-password
// @Produce      json
// @Param        status  query  string  false  "Filter status (pending/approved/cancelled)"
// @Success      200     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /api/super-admin/reset-password [get]
// @Security     ApiKeyAuth
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

