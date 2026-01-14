package handler

import (
	"strconv"

	req "backend/dto/request/masterdata"
	usecase "backend/usecase/masterdata"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type JurusanHandler struct {
	jurusanUsecase usecase.JurusanUsecase
}

func NewJurusanHandler(u usecase.JurusanUsecase) *JurusanHandler {
	return &JurusanHandler{
		jurusanUsecase: u,
	}
}

//
// ==========================
// CREATE
// ==========================

// CreateJurusan godoc
// @Summary      Create new jurusan
// @Description  Membuat jurusan baru (admin only)
// @Tags         jurusan
// @Accept       json
// @Produce      json
// @Param        request  body      data.JurusanRequest  true  "Jurusan data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /api/admin/jurusan [post]
// @Security     ApiKeyAuth
func (h *JurusanHandler) CreateJurusan(c *gin.Context) {
	var body req.JurusanRequest

	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.jurusanUsecase.Create(body)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result, "jurusan berhasil dibuat")
}

//
// ==========================
// GET ALL
// ==========================

// GetAllJurusan godoc
// @Summary      Get all jurusan
// @Description  Mengambil semua data jurusan
// @Tags         jurusan
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/admin/jurusan [get]
// @Security     ApiKeyAuth
func (h *JurusanHandler) GetAllJurusan(c *gin.Context) {
	data, err := h.jurusanUsecase.GetAll()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, data, "data jurusan berhasil diambil")
}

//
// ==========================
// GET BY FAKULTAS ID
// ==========================

// GetJurusanByFakultas godoc
// @Summary      Get jurusan by fakultas ID
// @Description  Mengambil data jurusan berdasarkan fakultas ID
// @Tags         jurusan
// @Produce      json
// @Param        fakultas_id  path  int  true  "Fakultas ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/admin/jurusan/fakultas/{fakultas_id} [get]
// @Security     ApiKeyAuth
func (h *JurusanHandler) GetJurusanByFakultas(c *gin.Context) {
	fakultasID, err := strconv.ParseUint(c.Param("fakultas_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID fakultas tidak valid")
		return
	}

	data, err := h.jurusanUsecase.GetByFakultas(uint(fakultasID))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, data, "data jurusan berhasil diambil")
}

//
// ==========================
// UPDATE
// ==========================

// UpdateJurusan godoc
// @Summary      Update jurusan
// @Description  Update data jurusan berdasarkan ID
// @Tags         jurusan
// @Accept       json
// @Produce      json
// @Param        id       path  int                 true  "Jurusan ID"
// @Param        request  body  data.JurusanRequest  true  "Jurusan data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /api/admin/jurusan/{id} [put]
// @Security     ApiKeyAuth
func (h *JurusanHandler) UpdateJurusan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID jurusan tidak valid")
		return
	}

	var body req.JurusanRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.jurusanUsecase.Update(uint(id), body)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "jurusan berhasil diupdate")
}

//
// ==========================
// DELETE
// ==========================

// DeleteJurusan godoc
// @Summary      Delete jurusan
// @Description  Menghapus data jurusan berdasarkan ID
// @Tags         jurusan
// @Produce      json
// @Param        id   path  int  true  "Jurusan ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/admin/jurusan/{id} [delete]
// @Security     ApiKeyAuth
func (h *JurusanHandler) DeleteJurusan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID jurusan tidak valid")
		return
	}

	if err := h.jurusanUsecase.Delete(uint(id)); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "jurusan berhasil dihapus")
}