package handler

import (
	request "backend/dto/request/barang"
	usecase "backend/usecase/barang"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type BarangUnitHandler struct {
	uc usecase.BarangUnitUseCase
}

func NewBarangUnitHandler(uc usecase.BarangUnitUseCase) *BarangUnitHandler {
	return &BarangUnitHandler{uc: uc}
}

//
// ==========================
// CREATE UNIT
// ==========================

// CreateUnit godoc
// @Summary      Create barang unit (bulk)
// @Description  Menambahkan unit ke barang secara bulk
// @Tags         barang-unit
// @Accept       json
// @Produce      json
// @Param        request body request.CreateUnitRequest true "create unit"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/unit [post]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) Create(c *gin.Context) {
	var body request.CreateUnitRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.Create(body)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, result, "unit berhasil ditambahkan")
}

//
// ==========================
// GET UNIT BY KODE
// ==========================

// GetUnitByKode godoc
// @Summary      Get unit detail
// @Description  Mendapatkan detail unit berdasarkan kode
// @Tags         barang-unit
// @Produce      json
// @Param        kode path string true "kode unit"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode} [get]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) GetByKode(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	result, err := h.uc.GetByKode(kode)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "detail unit")
}

//
// ==========================
// GET UNITS BY BARANG KODE
// ==========================

// GetUnitsByBarangKode godoc
// @Summary      Get units by barang kode
// @Description  Mendapatkan semua unit dari barang tertentu
// @Tags         barang-unit
// @Produce      json
// @Param        kode path string true "kode barang"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/admin/barang/{kode}/units [get]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) GetByBarangKode(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode barang tidak valid")
		return
	}

	result, err := h.uc.GetByBarangKode(kode)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "list unit barang")
}

//
// ==========================
// UPDATE UNIT METADATA
// ==========================

// UpdateUnit godoc
// @Summary      Update unit metadata
// @Description  Update kondisi, tahun perolehan, lokasi, dan catatan unit
// @Tags         barang-unit
// @Accept       json
// @Produce      json
// @Param        kode path string true "kode unit"
// @Param        request body request.UpdateUnitRequest true "update unit"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode} [put]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) Update(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	var body request.UpdateUnitRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.UpdateByKode(kode, body)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, result, "unit berhasil diperbarui")
}

//
// ==========================
// SET MAINTENANCE
// ==========================

// SetMaintenance godoc
// @Summary      Set unit to maintenance
// @Description  Set unit ke status maintenance dengan priority tertentu
// @Tags         barang-unit
// @Accept       json
// @Produce      json
// @Param        kode path string true "kode unit"
// @Param        request body request.SetMaintenanceRequest true "set maintenance"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode}/maintenance [post]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) SetMaintenance(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	var body request.SetMaintenanceRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.uc.SetMaintenance(kode, body.Priority, body.Alasan); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "unit berhasil di-set untuk maintenance")
}

//
// ==========================
// SET NON-AKTIF
// ==========================

// SetNonAktif godoc
// @Summary      Set unit to non-aktif
// @Description  Nonaktifkan unit dengan alasan tertentu
// @Tags         barang-unit
// @Accept       json
// @Produce      json
// @Param        kode path string true "kode unit"
// @Param        request body request.SetNonAktifRequest true "set non-aktif"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode}/nonaktif [post]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) SetNonAktif(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	var body request.SetNonAktifRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.uc.SetNonAktif(kode, body.Alasan); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "unit berhasil dinonaktifkan")
}

//
// ==========================
// SET AKTIF KEMBALI
// ==========================

// SetAktifKembali godoc
// @Summary      Reactivate unit
// @Description  Aktifkan kembali unit yang nonaktif atau maintenance
// @Tags         barang-unit
// @Produce      json
// @Param        kode path string true "kode unit"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode}/aktifkan [post]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) SetAktifKembali(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	if err := h.uc.SetAktifKembali(kode); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "unit berhasil diaktifkan kembali")
}

//
// ==========================
// DELETE UNIT
// ==========================

// DeleteUnit godoc
// @Summary      Delete unit
// @Description  Hapus unit (soft delete)
// @Tags         barang-unit
// @Produce      json
// @Param        kode path string true "kode unit"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode} [delete]
// @Security     ApiKeyAuth
func (h *BarangUnitHandler) Delete(c *gin.Context) {
	kode := c.Param("kode")
	if kode == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	if err := h.uc.DeleteByKode(kode); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "unit berhasil dihapus")
}