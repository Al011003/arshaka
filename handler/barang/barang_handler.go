// handler/barang_handler.go
package handler

import (
	"fmt"
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

// Create godoc
// @Summary      Create new barang
// @Description  Membuat barang baru (admin only)
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        request  body      barang.CreateBarangRequest  true  "Barang data"
// @Success      201      {object}  map[string]interface{}  "Barang berhasil dibuat"
// @Failure      400      {object}  map[string]interface{}  "Bad request - validation error"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang [post]
// @Security     ApiKeyAuth
func (h *BarangHandler) Create(c *gin.Context) {
	var req barang.CreateBarangRequest

	// Bind & Validate
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Custom validation
	if err := req.CustomValidate(); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Call usecase
	result, err := h.barangUsecase.Create(req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "barang berhasil dibuat",
		"data":    result,
	})
}

// GetAll godoc
// @Summary      Get all barang
// @Description  Mengambil semua data barang dengan filter dan pagination
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        status      query     string  false  "Filter by status (aktif/tidak_aktif)"
// @Param        kategori    query     string  false  "Filter by kategori"
// @Param        search      query     string  false  "Search by nama or kode"
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        limit       query     int     false  "Items per page (default: 10)"
// @Param        sort_by     query     string  false  "Sort field (default: created_at)"
// @Param        sort_order  query     string  false  "Sort order: asc/desc (default: desc)"
// @Success      200         {object}  map[string]interface{}  "Data barang berhasil diambil"
// @Failure      400         {object}  map[string]interface{}  "Bad request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      500         {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang [get]
// @Security     ApiKeyAuth
func (h *BarangHandler) GetAll(c *gin.Context) {
    // Cek role
    role := c.GetString("role")
    fmt.Println("🔥 ROLE YANG MASUK:", role)

    // Parse filter
    var filter barang.BarangFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    fmt.Println("🔥 FILTER.STATUS:", filter.Status)

	// Set default values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// 🔥 Panggil usecase dengan role
	data, pagination, err := h.barangUsecase.GetAll(filter, role)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"status":     "success",
		"message":    "data barang berhasil diambil",
		"data":       data,
		"pagination": pagination,
	})
}



// GetByID godoc
// @Summary      Get barang by ID
// @Description  Mengambil detail barang berdasarkan ID
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Barang ID"
// @Success      200  {object}  map[string]interface{}  "Data barang berhasil diambil"
// @Failure      400  {object}  map[string]interface{}  "Bad request - ID tidak valid"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      404  {object}  map[string]interface{}  "Barang tidak ditemukan"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang/{id} [get]
// @Security     ApiKeyAuth
func (h *BarangHandler) GetByID(c *gin.Context) {
	// Ambil ID barang dari URL
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	// Ambil role dari context (diinjeksi dari middleware JWT)
	roleValue, exists := c.Get("role")
	if !exists {
		utils.Unauthorized(c, "role tidak ditemukan dalam token")
		return
	}

	role, ok := roleValue.(string)
	if !ok {
		utils.InternalError(c, "format role tidak valid")
		return
	}

	// Panggil usecase dengan role
	result, err := h.barangUsecase.GetByID(uint(barangID), role)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "data barang berhasil diambil")
}

// Update godoc
// @Summary      Update barang
// @Description  Update data barang berdasarkan ID (admin only)
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        id       path      string                      true  "Barang ID"
// @Param        request  body      barang.UpdateBarangRequest  true  "Updated barang data"
// @Success      200      {object}  map[string]interface{}  "Barang berhasil diupdate"
// @Failure      400      {object}  map[string]interface{}  "Bad request - validation error"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Failure      404      {object}  map[string]interface{}  "Barang tidak ditemukan"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang/{id} [put]
// @Security     ApiKeyAuth
func (h *BarangHandler) Update(c *gin.Context) {
	// Get ID from URL param
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	var req barang.UpdateBarangRequest

	// Bind & Validate
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Cek apakah ada field yang diupdate
	if !req.HasUpdates() {
		utils.BadRequest(c, "tidak ada data yang diupdate")
		return
	}

	// Call usecase
	result, err := h.barangUsecase.Update(uint(barangID), req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "barang berhasil diupdate")
}

// Delete godoc
// @Summary      Delete barang
// @Description  Menghapus barang berdasarkan ID (admin only)
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Barang ID"
// @Success      200  {object}  map[string]interface{}  "Barang berhasil dihapus"
// @Failure      400  {object}  map[string]interface{}  "Bad request - ID tidak valid"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      404  {object}  map[string]interface{}  "Barang tidak ditemukan"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/admin/barang/{id} [delete]
// @Security     ApiKeyAuth
func (h *BarangHandler) Delete(c *gin.Context) {
	// Get ID from URL param
	id := c.Param("id")
	barangID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID barang tidak valid")
		return
	}

	// Call usecase
	err = h.barangUsecase.Delete(uint(barangID))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang berhasil dihapus")
}