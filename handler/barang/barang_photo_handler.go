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

// UpdatePhoto - POST /admin/barang/:id/photo
func (h *BarangPhotoHandler) UpdatePhoto(c *gin.Context) {
	// Get ID from URL param
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	// Validate request
	var req barang.BarangPhotoRequest
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Upload foto
	url, err := h.barangPhotoUsecase.UpdatePhoto(uint(barangID), req.File)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"cover_url": url,
	}, "foto barang berhasil diupload")
}

// DeletePhoto - DELETE /admin/barang/:id/photo
func (h *BarangPhotoHandler) DeletePhoto(c *gin.Context) {
	// Get ID from URL param
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	// Delete foto
	err = h.barangPhotoUsecase.DeletePhoto(uint(barangID))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "foto barang berhasil dihapus")
}