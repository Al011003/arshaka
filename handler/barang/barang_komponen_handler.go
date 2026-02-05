// handler/barang_komponen_handler.go
package handler

import (
	request "backend/dto/request/barang"
	usecase "backend/usecase/barang"
	"backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BarangKomponenHandler struct {
	uc usecase.BarangKomponenUseCase
}

func NewBarangKomponenHandler(uc usecase.BarangKomponenUseCase) *BarangKomponenHandler {
	return &BarangKomponenHandler{uc: uc}
}

// AddKomponen godoc
// @Summary      Add komponen to unit
// @Description  Menambahkan komponen ke unit barang
// @Tags         barang-komponen
// @Accept       json
// @Produce      json
// @Param        kode_unit path string true "kode unit"
// @Param        request body barang.CreateKomponenRequest true "create komponen"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode}/komponen [post]
// @Security     ApiKeyAuth
func (h *BarangKomponenHandler) Add(c *gin.Context) {
	kodeUnit := c.Param("kode")
	if kodeUnit == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	var body request.CreateKomponenRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.uc.Add(kodeUnit, body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, nil, "komponen berhasil ditambahkan")
}

// UpdateKomponen godoc
// @Summary      Update komponen
// @Description  Update data komponen
// @Tags         barang-komponen
// @Accept       json
// @Produce      json
// @Param        id path int true "komponen id"
// @Param        request body barang.UpdateKomponenRequest true "update komponen"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/komponen/{id} [put]
// @Security     ApiKeyAuth
func (h *BarangKomponenHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "id komponen tidak valid")
		return
	}

	var body request.UpdateKomponenRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.uc.Update(uint(id), body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "komponen berhasil diperbarui")
}

// DeleteKomponen godoc
// @Summary      Delete komponen
// @Description  Hapus komponen
// @Tags         barang-komponen
// @Produce      json
// @Param        id path int true "komponen id"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/admin/barang/komponen/{id} [delete]
// @Security     ApiKeyAuth
func (h *BarangKomponenHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "id komponen tidak valid")
		return
	}

	if err := h.uc.Delete(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "komponen berhasil dihapus")
}

// GetKomponenByUnit godoc
// @Summary      Get komponen by unit
// @Description  Mendapatkan semua komponen dari unit tertentu
// @Tags         barang-komponen
// @Produce      json
// @Param        kode_unit path string true "kode unit"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/admin/barang/unit/{kode}/komponen [get]
// @Security     ApiKeyAuth
func (h *BarangKomponenHandler) GetByUnit(c *gin.Context) {
	kodeUnit := c.Param("kode")
	if kodeUnit == "" {
		utils.BadRequest(c, "kode unit tidak valid")
		return
	}

	result, err := h.uc.GetByUnitKode(kodeUnit)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "list komponen unit")
}