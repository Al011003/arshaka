package handler

import (
	"net/http"
	"strconv"

	req "backend/dto/request/masterdata"
	"backend/dto/response"
	usecase "backend/usecase/masterdata"

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

// Create
func (h *JurusanHandler) CreateJurusan(c *gin.Context) {
	var request req.JurusanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.jurusanUsecase.Create(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  "success",
		Message: "jurusan berhasil dibuat",
		Data:    resp,
	})
}

// Get All
func (h *JurusanHandler) GetAllJurusan(c *gin.Context) {
	resp, err := h.jurusanUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// jika data kosong
	if len(resp) == 0 {
		c.JSON(http.StatusOK, response.BaseResponse{
			Status:  "success",
			Message: "data jurusan kosong",
			Data:    []interface{}{}, // array kosong
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  "success",
		Message: "berhasil mengambil data jurusan",
		Data:    resp,
	})
}

// Update
func (h *JurusanHandler) UpdateJurusan(c *gin.Context) {
	var request req.JurusanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  "error",
			Message: "invalid ID",
		})
		return
	}

	resp, err := h.jurusanUsecase.Update(uint(id), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  "success",
		Message: "jurusan berhasil diupdate",
		Data:    resp,
	})
}

// Delete
func (h *JurusanHandler) DeleteJurusan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  "error",
			Message: "invalid ID",
		})
		return
	}

	if err := h.jurusanUsecase.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  "success",
		Message: "jurusan berhasil dihapus",
	})
}

// Get by Fakultas ID
func (h *JurusanHandler) GetJurusanByFakultas(c *gin.Context) {
	fakultasIDParam := c.Param("fakultas_id")
	fakultasID, err := strconv.Atoi(fakultasIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  "error",
			Message: "invalid fakultas ID",
		})
		return
	}

	resp, err := h.jurusanUsecase.GetByFakultas(uint(fakultasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if len(resp) == 0 {
		c.JSON(http.StatusOK, response.BaseResponse{
			Status:  "success",
			Message: "jurusan tidak ditemukan untuk fakultas ini",
			Data:    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  "success",
		Message: "berhasil mengambil data jurusan",
		Data:    resp,
	})
}
