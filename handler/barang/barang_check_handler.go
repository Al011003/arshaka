package handler

import (
	"strconv"
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

// GET /api/barang/availabilitycheck/:id?start_date=yyyy-mm-dd&end_date=yyyy-mm-dd
func (h *BarangAvailabilityHandler) GetCalendar(c *gin.Context) {
	idStr := c.Param("barang_id")
	barangID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "id barang tidak valid")
		return
	}

	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	result, err := h.availabilityUC.GetCalendarSummary(uint(barangID), month)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "berhasil ambil kalender ketersediaan")
}
