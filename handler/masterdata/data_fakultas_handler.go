package handler

import (
	req "backend/dto/request/masterdata"
	usecase "backend/usecase/masterdata"
	"backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FakultasHandler struct {
	fakultasUsecase usecase.FakultasUsecase
}

func NewFakultasHandler(u usecase.FakultasUsecase) *FakultasHandler {
	return &FakultasHandler{
		fakultasUsecase: u,
	}
}

// CreateFakultas godoc
// @Summary      Create Fakultas
// @Description  Menambahkan data fakultas baru
// @Tags         Fakultas
// @Accept       json
// @Produce      json
// @Param        request body data.FakultasRequest true "Payload fakultas"
// @Success      201 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security     BearerAuth
// @Router       /api/admin/fakultas [post]
func (h *FakultasHandler) CreateFakultas(c *gin.Context) {
	var reqBody req.FakultasRequest

	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.fakultasUsecase.Create(reqBody)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, resp, "berhasil membuat fakultas")
}

//
// ==========================
// GET ALL FAKULTAS
// ==========================

// GetAllFakultas godoc
// @Summary      Get All Fakultas
// @Description  Mengambil seluruh data fakultas
// @Tags         Fakultas
// @Produce      json
// @Success      200 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security     BearerAuth
// @Router       /api/admin/fakultas [get]
func (h *FakultasHandler) GetAllFakultas(c *gin.Context) {
	resp, err := h.fakultasUsecase.GetAll()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	// FE-friendly response
	if len(resp) == 0 {
		utils.Success(c, []interface{}{}, "data kosong")
		return
	}

	utils.Success(c, resp, "berhasil mengambil data fakultas")
}

//
// ==========================
// UPDATE FAKULTAS
// ==========================

// UpdateFakultas godoc
// @Summary      Update Fakultas
// @Description  Mengupdate data fakultas berdasarkan ID
// @Tags         Fakultas
// @Accept       json
// @Produce      json
// @Param        id path int true "Fakultas ID"
// @Param        request body data.FakultasRequest true "Payload fakultas"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Security     BearerAuth
// @Router       /api/admin/fakultas/{id} [put]
func (h *FakultasHandler) UpdateFakultas(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.BadRequest(c, "id tidak valid")
		return
	}

	var reqBody req.FakultasRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.fakultasUsecase.Update(uint(id), reqBody)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, resp, "berhasil update fakultas")
}

//
// ==========================
// DELETE FAKULTAS
// ==========================

// DeleteFakultas godoc
// @Summary      Delete Fakultas
// @Description  Menghapus data fakultas
// @Tags         Fakultas
// @Produce      json
// @Param        id path int true "Fakultas ID"
// @Success      200 {object} response.BaseResponse
// @Failure      400 {object} response.BaseResponse
// @Security     BearerAuth
// @Router       /api/admin/fakultas/{id} [delete]
func (h *FakultasHandler) DeleteFakultas(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.BadRequest(c, "id tidak valid")
		return
	}

	if err := h.fakultasUsecase.Delete(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "berhasil menghapus fakultas")
}
