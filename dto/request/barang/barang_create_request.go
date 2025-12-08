// dto/request/barang/create.go
package barang

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateBarangRequest struct {
	Kode      string  `json:"kode"`
	Nama      string  `json:"nama"`
	Merk      string  `json:"merk"`
	Deskripsi string  `json:"deskripsi"`
	Kategori  string  `json:"kategori"`
	StokTotal int     `json:"stok_total"`
	StokSisa  int     `json:"stok_sisa"`
	TahunBeli int     `json:"tahun_beli"`
	HargaBeli float64 `json:"harga_beli"`
}

func (r *CreateBarangRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON ke struct
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	// Validasi dengan ozzo-validation
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
		validation.Field(&r.StokTotal, 
			validation.Required.Error("stok total wajib diisi"),
			validation.Min(0).Error("stok total tidak boleh negatif"),
		),
		validation.Field(&r.StokSisa, 
			validation.Required.Error("stok sisa wajib diisi"),
			validation.Min(0).Error("stok sisa tidak boleh negatif"),
		),
		validation.Field(&r.TahunBeli, 
			validation.Min(0).Error("tahun beli tidak valid"),
		),
		validation.Field(&r.HargaBeli, 
			validation.Min(0.0).Error("harga beli tidak boleh negatif"),
		),
	)
}

// CustomValidate - validasi tambahan yang butuh logic
func (r *CreateBarangRequest) CustomValidate() error {
	// Validasi stok sisa tidak boleh lebih dari stok total
	if r.StokSisa > r.StokTotal {
		return errors.New("stok sisa tidak boleh lebih dari stok total")
	}
	return nil
}