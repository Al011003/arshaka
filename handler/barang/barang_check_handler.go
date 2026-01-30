// handler/barang_check_avail_handler.go
package handler

import (
	"time"

	usecase "backend/usecase/barang"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type BarangAvailabilityHandler struct {
	availabilityUC usecase.AvailabilityUseCase
}

func NewBarangAvailabilityHandler(availUC usecase.AvailabilityUseCase) *BarangAvailabilityHandler {
	return &BarangAvailabilityHandler{availabilityUC: availUC}
}

// handler/availability_handler.go

// GetCalendar godoc
// @Summary      Get availability calendar
// @Description  Mendapatkan kalender ketersediaan barang untuk bulan tertentu
// @Tags         availability
// @Produce      json
// @Param        kode path string true "kode barang"
// @Param        month query string false "month (YYYY-MM)" default(current month)
// @Success      200 {object} map[string]interface{}
// @Router       /api/availability/{kode}/calendar [get]
func (h *BarangAvailabilityHandler) GetCalendar(c *gin.Context) {
	kode := c.Param("kode")
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	result, err := h.availabilityUC.GetCalendarSummary(kode, month)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, result, "availability calendar")
}

// GetDateDetail godoc
// @Summary      Get date detail
// @Description  Mendapatkan detail ketersediaan barang pada tanggal tertentu
// @Tags         availability
// @Produce      json
// @Param        kode path string true "kode barang"
// @Param        date query string true "date (YYYY-MM-DD)"
// @Success      200 {object} map[string]interface{}
// @Router       /api/availability/{kode}/date [get]
func (h *BarangAvailabilityHandler) GetDateDetail(c *gin.Context) {
	kode := c.Param("kode")
	date := c.Query("date")

	if date == "" {
		utils.BadRequest(c, "parameter date wajib diisi")
		return
	}

	result, err := h.availabilityUC.GetDateDetail(kode, date)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, result, "date detail")
}