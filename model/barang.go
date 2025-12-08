// models/barang.go - LEAN MODEL
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
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Barang) TableName() string {
    return "barang"
}

// Simple validation di model aja
func (b *Barang) BeforeCreate(tx *gorm.DB) error {
    if b.StokSisa > b.StokTotal {
        return errors.New("stok sisa tidak boleh lebih dari stok total")
    }
    return nil
}

func (b *Barang) BeforeUpdate(tx *gorm.DB) error {
    if b.StokSisa > b.StokTotal || b.StokSisa < 0 {
        return errors.New("stok tidak valid")
    }
    return nil
}