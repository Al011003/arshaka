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

// Create - POST /admin/barang
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

// GetAll - GET /admin/barang
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



// GetByID - GET /admin/barang/:id
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

// Update - PUT /admin/barang/:id
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

// Delete - DELETE /admin/barang/:id
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