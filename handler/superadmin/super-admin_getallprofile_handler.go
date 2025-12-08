package handler

import (
	request "backend/dto/request/user"
	response "backend/dto/response/common"
	usecase "backend/usecase/super_admin/getprofile"
	"backend/utils"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin" // Pake Gin!
)

type SuperAdminUserHandler struct {
	UC usecase.SuperAdminGetUserUsecase
}

func NewSuperAdminUserHandler(uc usecase.SuperAdminGetUserUsecase) *SuperAdminUserHandler {
	return &SuperAdminUserHandler{UC: uc}
}

func (h *SuperAdminUserHandler) GetUsers(c *gin.Context) {
    // Get current user dari middleware/JWT
    idRaw, exists := c.Get("user_id")
    if !exists {
       utils.Unauthorized(c, "unauthorized")
        return
    }
    superAdminID := idRaw.(uint)
    
    // Parse filter dari query params
    var filter request.UserFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

	fmt.Printf("Filter parsed: Search=%s, Page=%d, Limit=%d\n", filter.Search, filter.Page, filter.Limit)

    // SET DEFAULT VALUES SEBELUM DIPAKE! ⚠️
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10 // Default limit
    }

    // Call usecase
    users, total, err := h.UC.GetUsers(superAdminID, filter)
    if err != nil {
        c.JSON(http.StatusForbidden, response.BaseResponse{
            Status:  "error",
            Message: err.Error(),
        })
        return
    }

    // Calculate total pages (sekarang aman, limit pasti > 0)
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