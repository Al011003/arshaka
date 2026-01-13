package request

import (
	"backend/utils"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateLoanRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string `json:"reason"`

	ParsedStartDate time.Time `json:"-"`
	ParsedEndDate   time.Time `json:"-"`
}
func (r *CreateLoanRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("payload tidak valid")
	}

	startDate, err := utils.ParseDate(r.StartDate)
	if err != nil {
		return errors.New("start_date: " + err.Error())
	}

	endDate, err := utils.ParseDate(r.EndDate)
	if err != nil {
		return errors.New("end_date: " + err.Error())
	}

	now := time.Now().Truncate(24 * time.Hour)

	r.ParsedStartDate = startDate
	r.ParsedEndDate = endDate
return validation.ValidateStruct(r,
	validation.Field(&r.StartDate,
		validation.Required.Error("tanggal mulai wajib diisi"),
		validation.By(func(_ interface{}) error {
			if startDate.Before(now) {
				return errors.New("tanggal mulai tidak boleh di masa lalu")
			}
			return nil
		}),
	),
	validation.Field(&r.EndDate,
		validation.Required.Error("tanggal selesai wajib diisi"),
		validation.By(func(_ interface{}) error {
			if !endDate.After(startDate) {
				return errors.New("tanggal selesai harus setelah tanggal mulai")
			}
			return nil
		}),
	),
	validation.Field(&r.Reason,
		validation.Required.Error("alasan peminjaman wajib diisi"),
		validation.Length(10, 500).Error("alasan harus 10–500 karakter"),
	),
)
}





