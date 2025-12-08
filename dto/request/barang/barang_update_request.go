// dto/request/barang/update.go
package barang

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdateBarangRequest struct {
	Kode      *string  `json:"kode,omitempty"`
	Nama      *string  `json:"nama,omitempty"`
	Merk      *string  `json:"merk,omitempty"`
	Deskripsi *string  `json:"deskripsi,omitempty"`
	Kategori  *string  `json:"kategori,omitempty"`
	StokTotal *int     `json:"stok_total,omitempty"`
	StokSisa  *int     `json:"stok_sisa,omitempty"`
	TahunBeli *int     `json:"tahun_beli,omitempty"`
	HargaBeli *float64 `json:"harga_beli,omitempty"`
}

func (r *UpdateBarangRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON ke struct
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("payload tidak valid")
	}

	// Validasi dengan ozzo-validation (hanya field yang ada)
	return validation.ValidateStruct(r,
		validation.Field(&r.Kode,
			validation.Length(1, 50).Error("kode harus 1-50 karakter"),
		),
		validation.Field(&r.Nama,
			validation.Length(1, 200).Error("nama harus 1-200 karakter"),
		),
		validation.Field(&r.Merk,
			validation.Length(0, 100).Error("merk maksimal 100 karakter"),
		),
		validation.Field(&r.Kategori,
			validation.Length(1, 100).Error("kategori harus 1-100 karakter"),
		),
		validation.Field(&r.StokTotal,
			validation.Min(0).Error("stok total tidak boleh negatif"),
		),
		validation.Field(&r.StokSisa,
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

// HasUpdates - cek apakah ada field yang diupdate
func (r *UpdateBarangRequest) HasUpdates() bool {
	return r.Kode != nil || r.Nama != nil || r.Merk != nil ||
		r.Deskripsi != nil || r.Kategori != nil ||
		r.StokTotal != nil || r.StokSisa != nil ||
		r.TahunBeli != nil || r.HargaBeli != nil
}