package handler

import (
	"encoding/json"
	"net/http"

	req "backend/dto/request/user"
	res "backend/dto/response/common"
	apperrors "backend/errors"
	usecase "backend/usecase/super_admin/update"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type SuperAdminSelfUpdateHandler struct {
	selfUC usecase.SuperAdminSelfUpdateUsecase
}

func NewSuperAdminSelfUpdateHandler(selfUC usecase.SuperAdminSelfUpdateUsecase) *SuperAdminSelfUpdateHandler {
	return &SuperAdminSelfUpdateHandler{
		selfUC: selfUC,
	}
}

// ===========================================
// SUPERADMIN UPDATE DIRI SENDIRI
// ===========================================
func (h *SuperAdminSelfUpdateHandler) SuperAdminUpdateSelf(c *gin.Context) {
	// Ambil superadminID dari JWT
	idRaw, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	superAdminID := idRaw.(uint)

	var reqBody req.SuperAdminSelfUpdateRequest

	// 🔥 Decoder strict untuk mencegah field tambahan
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "invalid request body: " + err.Error(),
		})
		return
	}

	// Validasi request
	if err := reqBody.Validate(); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Panggil usecase
	if err := h.selfUC.SuperAdminUpdateSelf(superAdminID, reqBody); err != nil {
		if apperrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, res.BaseResponse{
				Status:  "error",
				Message: "superadmin tidak ditemukan",
			})
			return
		}

		utils.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "profile superadmin updated successfully",
	})
}
