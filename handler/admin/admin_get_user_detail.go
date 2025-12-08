package handler

import (
	res "backend/dto/response/common"
	usecase "backend/usecase/admin/profile"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminGetUserHandler struct {
	adminUserDetailUC usecase.AdminUserDetailUsecase
}

func NewAdminGetUserHandler(adminUserDetailUC usecase.AdminUserDetailUsecase) *AdminGetUserHandler {
	return &AdminGetUserHandler{
		adminUserDetailUC: adminUserDetailUC,
	}
}

func (h *AdminGetUserHandler) GetDetailUser(c *gin.Context) {
	// ambil param id user yang mau diliat
	idStr := c.Param("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "invalid user id",
		})
		return
	}

	// panggil usecase
	userDetail, err := h.adminUserDetailUC.GetDetail(uint(targetID))
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "user detail retrieved successfully",
		Data:    userDetail,
	})
}
