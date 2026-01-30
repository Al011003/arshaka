package barang

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdateBarangRequest struct {
	Kode      *string `json:"kode"`
	Nama      *string `json:"nama"`
	Merk      *string `json:"merk"`
	Deskripsi *string `json:"deskripsi"`
	Kategori  *string `json:"kategori"`
	Status    *string `json:"status"`
	Harga *float64 `json:"harga"`
}

func (r *UpdateBarangRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Kode,
			validation.Length(1, 50).Error("kode maksimal 50 karakter"),
		),
		validation.Field(&r.Nama,
			validation.Length(1, 200).Error("nama maksimal 200 karakter"),
		),
		validation.Field(&r.Merk,
			validation.Length(0, 100).Error("merk maksimal 100 karakter"),
		),
		validation.Field(&r.Kategori,
			validation.Length(1, 100).Error("kategori maksimal 100 karakter"),
		),
	
	)
}
