package handler

import (
	request "backend/dto/request/user"
	response "backend/dto/response/common"
	usecase "backend/usecase/admin/profile" // Sesuaikan path-nya

	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	UC usecase.AdminGetUserUsecase
}

func NewAdminUserHandler(uc usecase.AdminGetUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{UC: uc}
}

func (h *AdminUserHandler) GetUsers(c *gin.Context) {
	// Get admin ID dari middleware/JWT
	idRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Status:  "error",
			Message: "unauthorized",
		})
		return
	}
	adminID := idRaw.(uint)

	// Parse filter dari query params
	var filter request.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  "error",
			Message: "Invalid query parameters",
		})
		return
	}

	// Set default values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	// Call usecase
	users, total, err := h.UC.GetUsers(adminID, filter)
	if err != nil {
		c.JSON(http.StatusForbidden, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := (total + filter.Limit - 1) / filter.Limit

	// Return dengan PaginatedResponse
	c.JSON(http.StatusOK, response.PaginatedResponse{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
		Pagination: response.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalRows:  total,
			TotalPages: totalPages,
		},
	})
}