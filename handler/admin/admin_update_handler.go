package handler

import (
	"net/http"
	"strconv"

	req "backend/dto/request/user"
	res "backend/dto/response/common"
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
func (h *AdminUpdateHandler) AdminUpdateSelf(c *gin.Context) {
    // Ambil user ID dari JWT
    idRaw, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, res.BaseResponse{
            Status:  "error",
            Message: "unauthorized",
        })
        return
    }
    adminID := idRaw.(uint)

    var reqBody req.AdminSelfUpdateRequest

    // Bind + Validate dari DTO
    if err := reqBody.BindAndValidate(c); err != nil {
        c.JSON(http.StatusBadRequest, res.BaseResponse{
            Status:  "error",
            Message: err.Error(),
        })
        return
    }

    // Usecase update
    if err := h.selfUC.AdminUpdateSelf(adminID, reqBody); err != nil {
        // Jika tidak ditemukan
        if apperrors.IsNotFound(err) {
            c.JSON(http.StatusNotFound, res.BaseResponse{
                Status:  "error",
                Message: "admin tidak ditemukan",
            })
            return
        }

        c.JSON(http.StatusBadRequest, res.BaseResponse{
            Status:  "error",
            Message: err.Error(),
        })
        return
    }

    // SUCCESS pakai utils.Success (sesuai style lo)
    utils.Success(c, nil, "profile updated successfully")
}


//
// ===========================================
// 2. ADMIN UPDATE USER BIASA (BUKAN ADMIN/BKAN SUPERADMIN)
// ===========================================
func (h *AdminUpdateHandler) AdminUpdateUser(c *gin.Context) {
	roleRaw, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, res.BaseResponse{
			Status:  "error",
			Message: "unauthorized",
		})
		return
	}

	adminRole := roleRaw.(string)

	// Parse user ID
	uidStr := c.Param("id")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "invalid user id",
		})
		return
	}

	var req req.AdminUpdateUserRequest
	if err := req.BindandValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}


	resp, err := h.userUC.AdminUpdateUser(uint(uid), req, adminRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

utils.Success(c, resp, "user updated successfully")
}
