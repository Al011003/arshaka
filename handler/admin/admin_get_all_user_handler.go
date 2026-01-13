package handler

import (
	request "backend/dto/request/user"
	response "backend/dto/response/common"
	usecase "backend/usecase/admin/profile"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	UC usecase.AdminGetUserUsecase
}

func NewAdminUserHandler(uc usecase.AdminGetUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{UC: uc}
}

// GetUsers godoc
// @Summary      Get list of users
// @Description  Admin mengambil list user
// @Tags         admin-user
// @Accept       json
// @Produce      json
// @Param        page   query  int  false  "Page number"
// @Param        limit  query  int  false  "Limit per page"
// @Success      200    {object}  response.PaginatedResponse
// @Failure      401    {object}  response.BaseResponse
// @Failure      403    {object}  response.BaseResponse
// @Router       /api/admin/user [get]
// @Security     ApiKeyAuth
func (h *AdminUserHandler) GetUsers(c *gin.Context) {

	// === Auth Context ===
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}
	adminID := idRaw.(uint)

	// === Query Filter ===
	var filter request.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, "invalid query parameters")
		return
	}

	// === Pagination Normalization ===
	filter.Page, filter.Limit = utils.ValidatePagination(filter.Page, filter.Limit)

	// === Usecase ===
	users, total, err := h.UC.GetUsers(adminID, filter)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	// === Response ===
	totalPages := utils.CalcTotalPages(total, filter.Limit)

	utils.PaginatedSuccess(
			c,
			users,
			response.Pagination{
				Page:       filter.Page,
				Limit:      filter.Limit,
				TotalRows:  total,
				TotalPages: totalPages,
			},
			"users retrieved successfully",
		)
}
