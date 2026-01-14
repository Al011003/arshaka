package handler

import (
	"strconv"

	req "backend/dto/request/angkatan_mapala"
	usecase "backend/usecase/angkatan_mapala"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type AngkatanMapalaHandler struct {
	uc usecase.AngkatanMapalaUsecase
}

func NewAngkatanMapalaHandler(u usecase.AngkatanMapalaUsecase) *AngkatanMapalaHandler {
	return &AngkatanMapalaHandler{uc: u}
}

//
// ==========================
// CREATE
// ==========================

// CreateAngkatanMapala godoc
// @Summary      Create angkatan mapala
// @Description  Admin membuat angkatan mapala baru
// @Tags         angkatan-mapala
// @Accept       json
// @Produce      json
// @Param        request  body      angkatan.AngkatanMapalaRequest  true  "Angkatan mapala payload"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /api/admin/angkatan-mapala [post]
// @Security     ApiKeyAuth
func (h *AngkatanMapalaHandler) Create(c *gin.Context) {
	var body req.AngkatanMapalaRequest

	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	data, err := h.uc.Create(body)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, data, "berhasil membuat angkatan mapala")
}

//
// ==========================
// GET ALL
// ==========================

// GetAllAngkatanMapala godoc
// @Summary      Get all angkatan mapala
// @Description  Admin mengambil semua data angkatan mapala
// @Tags         angkatan-mapala
// @Produce      json
// @Success      200      {object}  map[string]interface{}
// @Router       /api/admin/angkatan-mapala [get]
// @Security     ApiKeyAuth
func (h *AngkatanMapalaHandler) GetAll(c *gin.Context) {
	data, err := h.uc.GetAll()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	if len(data) == 0 {
		utils.Success(c, []interface{}{}, "data kosong")
		return
	}

	utils.Success(c, data, "berhasil mengambil data")
}

//
// ==========================
// UPDATE
// ==========================

// UpdateAngkatanMapala godoc
// @Summary      Update angkatan mapala
// @Description  Admin update data angkatan mapala
// @Tags         angkatan-mapala
// @Accept       json
// @Produce      json
// @Param        id       path      int                            true  "Angkatan ID"
// @Param        request  body      angkatan.AngkatanMapalaRequest  true  "Update payload"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /api/admin/angkatan-mapala/{id} [put]
// @Security     ApiKeyAuth
func (h *AngkatanMapalaHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "id tidak valid")
		return
	}

	var body req.AngkatanMapalaRequest
	if err := body.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	data, err := h.uc.Update(uint(id), body)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, data, "berhasil update angkatan mapala")
}

//
// ==========================
// DELETE
// ==========================

// DeleteAngkatanMapala godoc
// @Summary      Delete angkatan mapala
// @Description  Admin menghapus angkatan mapala
// @Tags         angkatan-mapala
// @Produce      json
// @Param        id   path  int  true  "Angkatan ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/admin/angkatan-mapala/{id} [delete]
// @Security     ApiKeyAuth
func (h *AngkatanMapalaHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "id tidak valid")
		return
	}

	if err := h.uc.Delete(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "berhasil dihapus")
}
