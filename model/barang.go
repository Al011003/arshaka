// model/barang.go
package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Barang struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Kode        string    `json:"kode" gorm:"unique;not null;index"`
	Nama        string    `json:"nama" gorm:"not null"`
	Merk        string    `json:"merk"`
	Deskripsi   string    `json:"deskripsi" gorm:"type:text"`
	Kategori    string    `json:"kategori" gorm:"not null;index"`
	StokTotal   int       `json:"stok_total" gorm:"not null;default:0"`
	StokSisa    int       `json:"stok_sisa" gorm:"not null;default:0"`
	TahunBeli   int       `json:"tahun_beli"`
	HargaBeli   float64   `json:"harga_beli" gorm:"type:decimal(15,2)"`
	CoverURL    string    `json:"cover_url"`
	Status      string    `json:"status" gorm:"default:'tersedia'"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Barang) TableName() string {
	return "barang"
}

// BeforeCreate - validasi & set status sebelum create
func (b *Barang) BeforeCreate(tx *gorm.DB) error {
	// Validasi stok
	if b.StokSisa > b.StokTotal {
		return errors.New("stok sisa tidak boleh lebih dari stok total")
	}
	if b.StokSisa < 0 {
		return errors.New("stok sisa tidak boleh negatif")
	}
	if b.StokTotal < 0 {
		return errors.New("stok total tidak boleh negatif")
	}

	// Auto set status berdasarkan stok
	if b.StokSisa == 0 {
		b.Status = "habis"
	} else {
		b.Status = "tersedia"
	}

	return nil
}

// BeforeUpdate - validasi & set status sebelum update
func (b *Barang) BeforeUpdate(tx *gorm.DB) error {
    if b.StokSisa < 0 || b.StokSisa > b.StokTotal {
        return errors.New("stok tidak valid")
    }

    // Kalau status kosong atau bukan "nonaktif", auto-set
    if b.Status != "nonaktif" {
        if b.StokSisa > 0 {
            b.Status = "tersedia"
        } else {
            b.Status = "habis"
        }
    }

    return nil
}
