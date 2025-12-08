package handler

import (
	req "backend/dto/request/angkatan_mapala"
	res "backend/dto/response/common"
	usecase "backend/usecase/angkatan_mapala"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AngkatanMapalaHandler struct {
	uc usecase.AngkatanMapalaUsecase
}

func NewAngkatanMapalaHandler(u usecase.AngkatanMapalaUsecase) *AngkatanMapalaHandler {
	return &AngkatanMapalaHandler{
		uc: u,
	}
}

// CREATE
func (h *AngkatanMapalaHandler) Create(c *gin.Context) {
	var request req.AngkatanMapalaRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	data, err := h.uc.Create(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res.BaseResponse{
		Status:  "success",
		Message: "berhasil membuat angkatan mapala",
		Data:    data,
	})
}

// GET ALL
func (h *AngkatanMapalaHandler) GetAll(c *gin.Context) {
	data, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Jika data kosong
	if len(data) == 0 {
		c.JSON(http.StatusOK, res.BaseResponse{
			Status:  "success",
			Message: "data kosong",
			Data:    []interface{}{}, // supaya FE tidak error
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "berhasil mengambil data",
		Data:    data,
	})
}
// UPDATE
func (h *AngkatanMapalaHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "id tidak valid",
		})
		return
	}

	var request req.AngkatanMapalaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	data, err := h.uc.Update(uint(id), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "berhasil update angkatan mapala",
		Data:    data,
	})
}

// DELETE
func (h *AngkatanMapalaHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "id tidak valid",
		})
		return
	}

	if err := h.uc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "berhasil dihapus",
	})
}
