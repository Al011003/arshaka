package barang

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateKomponenRequest struct {
	Nama        string `json:"nama"`
	JumlahWajib int    `json:"jumlah_wajib"`
	JumlahAktual *int    `json:"jumlah_aktual"`
	HargaSatuan int64  `json:"harga_satuan"`
}

func (r *CreateKomponenRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Nama,
			validation.Required.Error("nama komponen wajib diisi"),
			validation.Length(1, 100).Error("nama komponen maksimal 100 karakter"),
		),
		validation.Field(&r.JumlahWajib,
			validation.Required.Error("jumlah wajib harus diisi"),
			validation.Min(1).Error("jumlah wajib minimal 1"),
			validation.Max(1000).Error("jumlah wajib maksimal 1000"),
		),
		validation.Field(&r.JumlahAktual,
			validation.Required.Error("jumlah wajib harus diisi"),
			validation.Min(1).Error("jumlah wajib minimal 1"),
			validation.Max(1000).Error("jumlah wajib maksimal 1000"),
		),
		validation.Field(&r.HargaSatuan,
			validation.Min(int64(0)).Error("harga satuan tidak boleh negatif"),
		),
	)
}