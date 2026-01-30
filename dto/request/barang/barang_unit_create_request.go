package barang

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateUnitRequest struct {
	KodeBarang      string `json:"kode_barang"`
	Jumlah          int    `json:"jumlah"`
	TahunPerolehan  int    `json:"tahun_perolehan"`
	LokasiPenyimpan string `json:"lokasi_penyimpan"`
	Catatan         string `json:"catatan"`
}

func (r *CreateUnitRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	currentYear := time.Now().Year()

	return validation.ValidateStruct(r,
		validation.Field(&r.KodeBarang,
			validation.Required.Error("kode barang wajib diisi"),
		),
		validation.Field(&r.Jumlah,
			validation.Min(1).Error("jumlah minimal 1"),
			validation.Max(100).Error("jumlah maksimal 100 unit sekaligus"),
		),
		validation.Field(&r.TahunPerolehan,
			validation.Min(1900).Error("tahun perolehan tidak valid"),
			validation.Max(currentYear).Error("tahun perolehan tidak boleh melebihi tahun sekarang"),
		),
		validation.Field(&r.LokasiPenyimpan,
			validation.Length(0, 100).Error("lokasi penyimpan maksimal 100 karakter"),
		),
		validation.Field(&r.Catatan,
			validation.Length(0, 500).Error("catatan maksimal 500 karakter"),
		),
	)
}