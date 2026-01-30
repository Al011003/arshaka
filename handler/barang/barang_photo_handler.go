package handler

import (
	barangReq "backend/dto/request/barang"
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

// ==========================
// UPDATE PHOTO (BY KODE)
// ==========================

// UpdatePhoto godoc
// @Summary      Upload / Update foto barang
// @Description  Upload atau mengganti foto utama barang (admin only)
// @Tags         barang-photo
// @Accept       multipart/form-data
// @Produce      json
// @Param        kode path string true "Kode Barang"
// @Param        file formData file true "Foto barang"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/admin/barang/{kode}/photo [post]
// @Security     ApiKeyAuth
func (h *BarangPhotoHandler) UpdatePhoto(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode barang tidak valid")
		return
	}

	var req barangReq.BarangPhotoRequest
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	url, err := h.barangPhotoUsecase.UpdatePhotoByKode(kode, req.File)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"cover_url": url,
	}, "foto barang berhasil diupload")
}

// ==========================
// DELETE PHOTO (BY KODE)
// ==========================

// DeletePhoto godoc
// @Summary      Hapus foto barang
// @Description  Menghapus foto utama barang (admin only)
// @Tags         barang-photo
// @Produce      json
// @Param        kode path string true "Kode Barang"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/admin/barang/{kode}/photo [delete]
// @Security     ApiKeyAuth
func (h *BarangPhotoHandler) DeletePhoto(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode barang tidak valid")
		return
	}

	if err := h.barangPhotoUsecase.DeletePhotoByKode(kode); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "foto barang berhasil dihapus")
}
