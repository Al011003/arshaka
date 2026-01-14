// handler/barang_photo_handler.go
package handler

import (
	"strconv"

	barang "backend/dto/request/barang"
	usecase "backend/usecase/barang"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type BarangPhotoHandler struct {
	barangPhotoUsecase usecase.BarangPhotoUsecase
}

func NewBarangPhotoHandler(barangPhotoUC usecase.BarangPhotoUsecase) *BarangPhotoHandler {
	return &BarangPhotoHandler{
		barangPhotoUsecase: barangPhotoUC,
	}
}

// UpdatePhoto godoc
// @Summary      Upload / Update foto barang
// @Description  Upload atau mengganti foto utama barang (admin only)
// @Tags         barang-photo
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      string  true  "ID Barang"
// @Param        file  formData  file    true  "Foto barang"
// @Success      200   {object}  map[string]interface{}  "Foto barang berhasil diupload"
// @Failure      400   {object}  map[string]interface{}  "Bad request - ID tidak valid / file tidak valid"
// @Failure      401   {object}  map[string]interface{}  "Unauthorized"
// @Failure      404   {object}  map[string]interface{}  "Barang tidak ditemukan"
// @Failure      500   {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang/{id}/photo [post]
// @Security     ApiKeyAuth
func (h *BarangPhotoHandler) UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	var req barang.BarangPhotoRequest
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	url, err := h.barangPhotoUsecase.UpdatePhoto(uint(barangID), req.File)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"cover_url": url,
	}, "foto barang berhasil diupload")
}

// DeletePhoto godoc
// @Summary      Hapus foto barang
// @Description  Menghapus foto utama barang (admin only)
// @Tags         barang-photo
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID Barang"
// @Success      200  {object}  map[string]interface{}  "Foto barang berhasil dihapus"
// @Failure      400  {object}  map[string]interface{}  "Bad request - ID tidak valid"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      404  {object}  map[string]interface{}  "Barang tidak ditemukan"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang/{id}/photo [delete]
// @Security     ApiKeyAuth
func (h *BarangPhotoHandler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	err = h.barangPhotoUsecase.DeletePhoto(uint(barangID))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "foto barang berhasil dihapus")
}
