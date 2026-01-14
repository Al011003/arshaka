package handler

import (
	req "backend/dto/request/user"
	usecase "backend/usecase/user/profile"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type UserPhotoHandler struct {
	photoUC usecase.UserPhotoUsecase
}

func NewUserPhotoHandler(photoUC usecase.UserPhotoUsecase) *UserPhotoHandler {
	return &UserPhotoHandler{photoUC}
}

// ==================== UPDATE / SET PHOTO ====================

// UpdatePhoto godoc
// @Summary      Update user profile photo
// @Description  Upload atau update foto profil user
// @Tags         user-profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Profile photo file"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/user/update/photo/upload [post]
// @Security     ApiKeyAuth
func (h *UserPhotoHandler) UpdatePhoto(c *gin.Context) {
	var request req.UserPhotoRequest

	if err := request.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uint)

	url, err := h.photoUC.UpdatePhoto(userID, request.File)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"foto_url": url}, "foto profil berhasil diperbarui")
}

// ==================== DELETE PHOTO ====================

// DeletePhoto godoc
// @Summary      Delete user profile photo
// @Description  Menghapus foto profil user
// @Tags         user-profile
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/user/update/photo/delete [delete]
// @Security     ApiKeyAuth
func (h *UserPhotoHandler) DeletePhoto(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	if err := h.photoUC.DeletePhoto(userID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "foto profil berhasil dihapus")
}
