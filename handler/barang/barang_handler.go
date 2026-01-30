package handler

import (
	barangReq "backend/dto/request/barang"
	barangUC "backend/usecase/barang"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type BarangHandler struct {
	uc barangUC.BarangUseCase
}

func NewBarangHandler(uc barangUC.BarangUseCase) *BarangHandler {
	return &BarangHandler{uc: uc}
}

//
// ==========================
// CREATE BARANG (MASTER)
// ==========================
//

// CreateBarang godoc
// @Summary      Create barang
// @Description  Membuat barang master (admin)
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        request body barangReq.CreateBarangRequest true "barang master"
// @Success      201 {object} map[string]interface{}
// @Router       /api/admin/barang [post]
// @Security     ApiKeyAuth
func (h *BarangHandler) Create(c *gin.Context) {
	var body barangReq.CreateBarangRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.Create(body)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result, "barang berhasil dibuat")
}

//
// ==========================
// GET ALL
// ==========================
//

// GetAllBarang godoc
// @Summary      Get list barang
// @Tags         barang
// @Produce      json
// @Param        keyword query string false "keyword"
// @Param        kategori query string false "kategori"
// @Param        status query string false "status"
// @Param        page query int false "page"
// @Param        limit query int false "limit"
// @Router       /api/barang [get]
func (h *BarangHandler) GetAll(c *gin.Context) {
	role := c.GetString("role")

	var filter barangReq.BarangFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	filter.Page, filter.Limit = utils.ValidatePagination(filter.Page, filter.Limit)

	data, pagination, err := h.uc.GetAll(filter, role)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.PaginatedSuccess(c, data, pagination, "list barang")
}

//
// ==========================
// GET BY KODE
// ==========================
//

// GetBarangByKode godoc
// @Summary      Get barang detail
// @Tags         barang
// @Produce      json
// @Param        kode path string true "barang kode"
// @Router       /api/barang/{kode} [get]
func (h *BarangHandler) GetByKode(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode barang wajib diisi")
		return
	}

	role := c.GetString("role")
	result, err := h.uc.GetByKode(kode, role)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "detail barang")
}

//
// ==========================
// UPDATE BY KODE
// ==========================
//

// UpdateBarang godoc
// @Summary      Update barang
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        kode path string true "barang kode"
// @Param        request body barangReq.UpdateBarangRequest true "update barang"
// @Router       /api/admin/barang/{kode} [put]
// @Security     ApiKeyAuth
func (h *BarangHandler) Update(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode barang wajib diisi")
		return
	}

	var body barangReq.UpdateBarangRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.UpdateByKode(kode, body)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "barang berhasil diupdate")
}

//
// ==========================
// DELETE BY KODE
// ==========================
//

// DeleteBarang godoc
// @Summary      Delete barang
// @Tags         barang
// @Produce      json
// @Param        kode path string true "barang kode"
// @Router       /api/admin/barang/{kode} [delete]
// @Security     ApiKeyAuth
func (h *BarangHandler) Delete(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode barang wajib diisi")
		return
	}

	if err := h.uc.DeleteByKode(kode); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang berhasil dihapus")
}


func (h *BarangHandler) SetNonAktif(c *gin.Context) {
	kode := c.Param("kode")

	if err := h.uc.SetNonAktif(kode); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Barang dan semua unitnya berhasil dinonaktifkan",
	})
}

func (h *BarangHandler) SetAktif(c *gin.Context) {
	kode := c.Param("kode")

	if err := h.uc.SetAktif(kode); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Barang dan semua unitnya berhasil diaktifkan kembali",
	})
}


func (h *BarangHandler) CheckStatus(c *gin.Context) {
	kode := c.Param("kode")
	
	result, err := h.uc.CheckAndSuggestBarangStatus(kode)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "status check")
}