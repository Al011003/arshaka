package barang

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)



type UpdateKomponenRequest struct {
	Nama         *string `json:"nama"`
	JumlahWajib  *int    `json:"jumlah_wajib"`
	JumlahAktual *int    `json:"jumlah_aktual"`
	HargaSatuan  *int64  `json:"harga_satuan"`
}

func (r *UpdateKomponenRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	// Validate nama jika ada
	if r.Nama != nil {
		if err := validation.Validate(r.Nama,
			validation.Required.Error("nama komponen wajib diisi"),
			validation.Length(1, 100).Error("nama komponen maksimal 100 karakter"),
		); err != nil {
			return err
		}
	}

	// Validate jumlah wajib jika ada
	if r.JumlahWajib != nil {
		if err := validation.Validate(r.JumlahWajib,
			validation.Min(1).Error("jumlah wajib minimal 1"),
			validation.Max(1000).Error("jumlah wajib maksimal 1000"),
		); err != nil {
			return err
		}
	}

	// Validate jumlah aktual jika ada
	if r.JumlahAktual != nil {
		if err := validation.Validate(r.JumlahAktual,
			validation.Min(0).Error("jumlah aktual tidak boleh negatif"),
			validation.Max(1000).Error("jumlah aktual maksimal 1000"),
		); err != nil {
			return err
		}
	}

	// Validate harga satuan jika ada
	if r.HargaSatuan != nil {
		if err := validation.Validate(r.HargaSatuan,
			validation.Min(int64(0)).Error("harga satuan tidak boleh negatif"),
		); err != nil {
			return err
		}
	}

	return nil
}