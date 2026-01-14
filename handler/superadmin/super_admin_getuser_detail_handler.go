package handler

import (
	res "backend/dto/response/common"
	usecase "backend/usecase/super_admin/getprofile"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuperAdminGetUserHandler struct {
	userDetailUC usecase.UserDetailUsecase
}

func NewSuperAdminGetUserHandler(userDetailUC usecase.UserDetailUsecase) *SuperAdminGetUserHandler {
	return &SuperAdminGetUserHandler{
		userDetailUC: userDetailUC,
	}
}

// GetDetailUser godoc
// @Summary      Get user detail by ID
// @Description  Mengambil detail user berdasarkan ID (super admin only)
// @Tags         super-admin-user
// @Produce      json
// @Param        id   path  int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/super-admin/user/{id} [get]
// @Security     ApiKeyAuth
func (h *SuperAdminGetUserHandler) GetDetailUser(c *gin.Context) {
	// ambil param id user yang mau diliat
	idStr := c.Param("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	// panggil usecase
	userDetail, err := h.userDetailUC.GetDetail(uint(targetID))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "user detail retrieved successfully",
		Data:    userDetail,
	})
}
