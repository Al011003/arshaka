package handler

import (
	request "backend/dto/request/user"
	usecase "backend/usecase/super_admin/getprofile"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type SuperAdminUserHandler struct {
	uc usecase.SuperAdminGetUserUsecase
}

func NewSuperAdminUserHandler(uc usecase.SuperAdminGetUserUsecase) *SuperAdminUserHandler {
	return &SuperAdminUserHandler{uc: uc}
}

//
// ==========================
// GET ALL USERS
// ==========================

// GetUsers godoc
// @Summary      Get all users
// @Description  Mengambil semua data user dengan filter dan pagination (super admin only)
// @Tags         super-admin-user
// @Produce      json
// @Param        search   query  string  false  "Search by name/email/nim"
// @Param        role     query  string  false  "Filter by role"
// @Param        page     query  int     false  "Page number"
// @Param        limit    query  int     false  "Items per page"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Router       /api/super-admin/user [get]
// @Security     ApiKeyAuth
func (h *SuperAdminUserHandler) GetUsers(c *gin.Context) {
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	superAdminID := idRaw.(uint)

	var filter request.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	filter.Page, filter.Limit = utils.ValidatePagination(filter.Page, filter.Limit)

	users, total, err := h.uc.GetUsers(superAdminID, filter)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit

	pagination := map[string]interface{}{
		"page":        filter.Page,
		"limit":       filter.Limit,
		"total_rows":  total,
		"total_pages": totalPages,
	}

	utils.PaginatedSuccess(c, users, pagination, "Users retrieved successfully")
}