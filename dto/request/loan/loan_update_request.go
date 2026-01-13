package request

import (
	"backend/utils"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateLoanRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string `json:"reason"`

	ParsedStartDate time.Time `json:"-"`
	ParsedEndDate   time.Time `json:"-"`
}

type UpdateLoanItemsRequest struct {
	Items []LoanItemUpdate `json:"items"`
}

type LoanItemUpdate struct {
	BarangID uint `json:"barang_id"`
	Quantity int  `json:"quantity"`
	Action   string `json:"action"` 
	// ADD | UPDATE | DELETE
}

func (r *UpdateLoanItemsRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("payload tidak valid")
	}

	if len(r.Items) == 0 {
		return errors.New("items tidak boleh kosong")
	}

	for i, item := range r.Items {
		if item.BarangID == 0 {
			return errors.New("barang_id wajib diisi")
		}

		switch item.Action {
		case "ADD", "UPDATE":
			if item.Quantity <= 0 {
				return errors.New("quantity harus lebih dari 0 untuk action " + item.Action)
			}
		case "DELETE":
			// quantity boleh 0
		default:
			return errors.New("action tidak valid di item ke-" + fmt.Sprint(i+1))
		}
	}

	return nil
}

func (r *UpdateLoanRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("payload tidak valid")
	}

	// minimal 1 field harus diisi
	if r.StartDate == "" && r.EndDate == "" && r.Reason == "" {
		return errors.New("minimal satu field harus diisi")
	}

	now := time.Now().Truncate(24 * time.Hour)

	// validate start_date (optional)
	if r.StartDate != "" {
		startDate, err := utils.ParseDate(r.StartDate)
		if err != nil {
			return errors.New("start_date: " + err.Error())
		}
		if startDate.Before(now) {
			return errors.New("tanggal mulai tidak boleh di masa lalu")
		}
		r.ParsedStartDate = startDate
	}

	// validate end_date (optional)
	if r.EndDate != "" {
		endDate, err := utils.ParseDate(r.EndDate)
		if err != nil {
			return errors.New("end_date: " + err.Error())
		}
		r.ParsedEndDate = endDate
	}

	// kalau dua-duanya ada → cek relasi
	if !r.ParsedStartDate.IsZero() && !r.ParsedEndDate.IsZero() {
		if !r.ParsedEndDate.After(r.ParsedStartDate) {
			return errors.New("tanggal selesai harus setelah tanggal mulai")
		}
	}

	// validate reason (optional)
	if r.Reason != "" {
		if len(r.Reason) < 10 {
			return errors.New("alasan minimal 10 karakter")
		}
	}

	return nil
}
