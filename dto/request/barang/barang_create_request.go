package barang

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateBarangRequest struct {
	Kode      string `json:"kode"`
	Nama      string `json:"nama"`
	Merk      string `json:"merk"`
	Deskripsi string `json:"deskripsi"`
	Kategori  string `json:"kategori"`
	Harga float64 `json:"harga"`
}

func (r *CreateBarangRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Kode,
			validation.Required.Error("kode wajib diisi"),
			validation.Length(1, 50).Error("kode maksimal 50 karakter"),
		),
		validation.Field(&r.Nama,
			validation.Required.Error("nama wajib diisi"),
			validation.Length(1, 200).Error("nama maksimal 200 karakter"),
		),
		validation.Field(&r.Merk,
			validation.Length(0, 100).Error("merk maksimal 100 karakter"),
		),
		validation.Field(&r.Kategori,
			validation.Required.Error("kategori wajib diisi"),
			validation.Length(1, 100).Error("kategori maksimal 100 karakter"),
		),
		validation.Field(&r.Harga,
			validation.Required.Error("harga wajib diisi"),
		),
	)
}
