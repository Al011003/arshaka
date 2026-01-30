// dto/request/barang/unit.go
package barang

import (
	"backend/model"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdateUnitRequest struct {
	Kondisi         *string `json:"kondisi"`
	TahunPerolehan  *int    `json:"tahun_perolehan"`
	LokasiPenyimpan *string `json:"lokasi_penyimpan"`
	Catatan         *string `json:"catatan"`
}

func (r *UpdateUnitRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	currentYear := time.Now().Year()

	// Validate kondisi jika ada
	if r.Kondisi != nil {
		if err := validation.Validate(r.Kondisi,
			validation.In(
				model.UnitKondisiBaik,
				model.UnitKondisiLecet,
				model.UnitKondisiRusak,
			).Error("kondisi harus: baik, lecet, atau rusak"),
		); err != nil {
			return err
		}
	}

	// Validate tahun perolehan jika ada
	if r.TahunPerolehan != nil {
		if err := validation.Validate(r.TahunPerolehan,
			validation.Min(1900).Error("tahun perolehan tidak valid"),
			validation.Max(currentYear).Error("tahun perolehan tidak boleh melebihi tahun sekarang"),
		); err != nil {
			return err
		}
	}

	// Validate lokasi jika ada
	if r.LokasiPenyimpan != nil {
		if err := validation.Validate(r.LokasiPenyimpan,
			validation.Length(0, 100).Error("lokasi penyimpan maksimal 100 karakter"),
		); err != nil {
			return err
		}
	}

	// Validate catatan jika ada
	if r.Catatan != nil {
		if err := validation.Validate(r.Catatan,
			validation.Length(0, 500).Error("catatan maksimal 500 karakter"),
		); err != nil {
			return err
		}
	}

	return nil
}

// dto/request/barang/unit.go

type SetMaintenanceRequest struct {
	Priority int    `json:"priority"`
	Alasan   string `json:"alasan"` // ← ganti dari Catatan jadi Alasan biar konsisten
}

func (r *SetMaintenanceRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Priority,
			validation.Required.Error("priority wajib diisi"),
			validation.In(
				model.FixPriorityLow,
				model.FixPriorityMedium,
				model.FixPriorityHigh,
			).Error("priority harus: 1 (low), 2 (medium), atau 3 (high)"),
		),
		validation.Field(&r.Alasan,
			validation.Required.Error("alasan wajib diisi"),
			validation.Length(5, 500).Error("alasan harus antara 5-500 karakter"),
		),
	)
}
type SetNonAktifRequest struct {
	Alasan string `json:"alasan"`
}

func (r *SetNonAktifRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Alasan,
			validation.Required.Error("alasan wajib diisi"),
			validation.Length(5, 500).Error("alasan harus antara 5-500 karakter"),
		),
	)
}