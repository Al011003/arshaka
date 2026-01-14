package handler

import (
	"strconv"

	barang "backend/dto/request/barang"
	usecase "backend/usecase/barang"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type BarangHandler struct {
	barangUsecase usecase.BarangUseCase
}

func NewBarangHandler(barangUC usecase.BarangUseCase) *BarangHandler {
	return &BarangHandler{
		barangUsecase: barangUC,
	}
}

//
// ==========================
// CREATE
// ==========================

// CreateBarang godoc
// @Summary      Create new barang
// @Description  Membuat barang baru (admin only)
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        request  body      barang.CreateBarangRequest  true  "Barang data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /api/admin/barang [post]
// @Security     ApiKeyAuth
func (h *BarangHandler) Create(c *gin.Context) {
	var body barang.CreateBarangRequest

	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := body.CustomValidate(); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.barangUsecase.Create(body)
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

// GetAllBarang godoc
// @Summary      Get all barang
// @Description  Mengambil semua data barang dengan filter dan pagination
// @Tags         barang
// @Produce      json
// @Param        status      query     string  false  "Filter status"
// @Param        kategori    query     string  false  "Filter kategori"
// @Param        search      query     string  false  "Search nama/kode"
// @Param        page        query     int     false  "Page"
// @Param        limit       query     int     false  "Limit"
// @Router       /api/admin/barang [get]
// @Security     ApiKeyAuth
func (h *BarangHandler) GetAll(c *gin.Context) {
	role := c.GetString("role")

	var filter barang.BarangFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	filter.Page, filter.Limit = utils.ValidatePagination(filter.Page, filter.Limit)

	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	data, pagination, err := h.barangUsecase.GetAll(filter, role)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.PaginatedSuccess(c, data, pagination, "data barang berhasil diambil")
}

//
// ==========================
// GET BY ID
// ==========================

// GetBarangByID godoc
// @Summary      Get barang by ID
// @Description  Mengambil detail barang berdasarkan ID
// @Tags         barang
// @Produce      json
// @Param        id   path  int  true  "Barang ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/admin/barang/{id} [get]
// @Security     ApiKeyAuth
func (h *BarangHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	role := c.GetString("role")
	if role == "" {
		utils.Unauthorized(c, "unauthorized")
		return
	}

	result, err := h.barangUsecase.GetByID(uint(id), role)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "data barang berhasil diambil")
}

//
// ==========================
// UPDATE
// ==========================

// UpdateBarang godoc
// @Summary      Update barang
// @Description  Update data barang berdasarkan ID
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        id       path  int                         true  "Barang ID"
// @Param        request  body  barang.UpdateBarangRequest  true  "Update data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /api/admin/barang/{id} [put]
// @Security     ApiKeyAuth
func (h *BarangHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	var body barang.UpdateBarangRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !body.HasUpdates() {
		utils.BadRequest(c, "tidak ada data yang diupdate")
		return
	}

	result, err := h.barangUsecase.Update(uint(id), body)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "barang berhasil diupdate")
}

//
// ==========================
// DELETE
// ==========================

// DeleteBarang godoc
// @Summary      Delete barang
// @Description  Menghapus barang berdasarkan ID
// @Tags         barang
// @Produce      json
// @Param        id   path  int  true  "Barang ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/admin/barang/{id} [delete]
// @Security     ApiKeyAuth
func (h *BarangHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	if err := h.barangUsecase.Delete(uint(id)); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang berhasil dihapus")
}