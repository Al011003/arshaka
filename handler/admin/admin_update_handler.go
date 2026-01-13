package handler

import (
	"strconv"

	req "backend/dto/request/user"
	apperrors "backend/errors"
	usecase "backend/usecase/admin/update"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type AdminUpdateHandler struct {
	selfUC usecase.AdminSelfUpdateUsecase
	userUC usecase.AdminUpdateUserUsecase
}

func NewAdminUpdateHandler(
	selfUC usecase.AdminSelfUpdateUsecase,
	userUC usecase.AdminUpdateUserUsecase,
) *AdminUpdateHandler {
	return &AdminUpdateHandler{
		selfUC: selfUC,
		userUC: userUC,
	}
}

//
// ===========================================
// 1. ADMIN UPDATE DIRI SENDIRI
// ===========================================

// AdminUpdateSelf godoc
// @Summary      Update admin profile
// @Description  Admin mengupdate data dirinya sendiri
// @Tags         admin-profile
// @Accept       json
// @Produce      json
// @Param        request  body      user.AdminSelfUpdateRequest  true  "Admin self update payload"
// @Success      200      {object}  map[string]interface{}  "Profile updated successfully"
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Router       /api/admin/update [put]
// @Security     ApiKeyAuth
func (h *AdminUpdateHandler) AdminUpdateSelf(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	adminID := idRaw.(uint)

	var body req.AdminSelfUpdateRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.selfUC.AdminUpdateSelf(adminID, body); err != nil {
		if apperrors.IsNotFound(err) {
			utils.NotFound(c, "admin tidak ditemukan")
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "profile updated successfully")
}

//
// ===========================================
// 2. ADMIN UPDATE USER BIASA
// ===========================================

// AdminUpdateUser godoc
// @Summary      Update user by admin
// @Description  Admin mengupdate data user biasa (bukan admin / super admin)
// @Tags         admin-user
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "User ID"
// @Param        request  body      user.AdminUpdateUserRequest  true  "Update user payload"
// @Success      200      {object}  map[string]interface{}  "User updated successfully"
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Router       /api/admin/user/update/{id} [put]
// @Security     ApiKeyAuth
func (h *AdminUpdateHandler) AdminUpdateUser(c *gin.Context) {
	roleRaw, exists := c.Get("role")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	adminRole := roleRaw.(string)

	uidStr := c.Param("id")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	var body req.AdminUpdateUserRequest
	if err := body.BindandValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.userUC.AdminUpdateUser(uint(uid), body, adminRole)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, resp, "user updated successfully")
}
