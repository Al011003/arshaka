package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	usecase "backend/usecase/admin/profile"
	"backend/utils"
)

type AdminGetUserHandler struct {
	uc usecase.AdminUserDetailUsecase
}

func NewAdminGetUserHandler(uc usecase.AdminUserDetailUsecase) *AdminGetUserHandler {
	return &AdminGetUserHandler{
		uc: uc,
	}
}

// GetDetailUser godoc
// @Summary      Get user detail
// @Description  Admin mengambil detail user berdasarkan ID
// @Tags         admin-user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{} "User detail"
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/admin/user/{id} [get]
// @Security     ApiKeyAuth
func (h *AdminGetUserHandler) GetDetailUser(c *gin.Context) {
	// Parse user ID
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	// Call usecase
	userDetail, err := h.uc.GetDetail(uint(userID))
	if err != nil {
		// asumsi usecase return error kalau data tidak ditemukan
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, userDetail, "user detail retrieved successfully")
}
