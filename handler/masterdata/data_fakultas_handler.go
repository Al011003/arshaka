package handler

import (
	req "backend/dto/request/masterdata"
	res "backend/dto/response/common"
	usecase "backend/usecase/masterdata"
	"net/http"
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

// CREATE
func (h *FakultasHandler) CreateFakultas(c *gin.Context) {
	var request req.FakultasRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.fakultasUsecase.Create(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res.BaseResponse{
		Status:  "success",
		Message: "berhasil membuat fakultas",
		Data:    resp,
	})
}

// GET ALL
func (h *FakultasHandler) GetAllFakultas(c *gin.Context) {
	resp, err := h.fakultasUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// jika data kosong
	if len(resp) == 0 {
		c.JSON(http.StatusOK, res.BaseResponse{
			Status:  "success",
			Message: "data kosong",
			Data:    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "berhasil mengambil data",
		Data:    resp,
	})
}

// UPDATE
func (h *FakultasHandler) UpdateFakultas(c *gin.Context) {
	var request req.FakultasRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "invalid ID",
		})
		return
	}

	resp, err := h.fakultasUsecase.Update(uint(id), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "berhasil update fakultas",
		Data:    resp,
	})
}

// DELETE
func (h *FakultasHandler) DeleteFakultas(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "invalid ID",
		})
		return
	}

	if err := h.fakultasUsecase.Delete(uint(id)); err != nil {
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
