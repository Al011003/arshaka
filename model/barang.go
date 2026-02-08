// model/barang.go
package model

import (
	"time"

	"gorm.io/gorm"
)

// ============= CONSTANTS =============

// Barang Status
const (
	BarangStatusAktif    = "aktif"
	BarangStatusNonAktif = "nonaktif"
)

// Unit Status
const (
	UnitStatusSiap                = "siap"
	UnitStatusDipinjam            = "dipinjam"
	UnitStatusSiapTerbatas        = "siap_terbatas"
	UnitStatusMaintenanceRequired = "maintenance_required"
	UnitStatusNonAktif            = "nonaktif"
	UnitStatusReserved = "reserved"
)

// Unit Kondisi
const (
	UnitKondisiBaik  = "baik"
	UnitKondisiLecet = "lecet"
	UnitKondisiRusak = "rusak"
)

// Komponen Status
const (
	KomponenStatusLengkap = "lengkap"
	KomponenStatusKurang  = "kurang"
	KomponenStatusRusak   = "rusak"
)

// Fix Priority
const (
	FixPriorityNormal = 0
	FixPriorityLow    = 1
	FixPriorityMedium = 2
	FixPriorityHigh   = 3
)

// ============= MODELS =============

type Barang struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Kode      string         `json:"kode" gorm:"unique;not null;index"`
	Nama      string         `json:"nama" gorm:"not null"`
	Merk      string         `json:"merk"`
	Deskripsi string         `json:"deskripsi" gorm:"type:text"`
	Kategori  string         `json:"kategori" gorm:"not null;index"`
	CoverURL  string         `json:"cover_url"`
	Status    string         `json:"status" gorm:"default:'aktif'"`
	Harga     float64        `json:"harga" gorm:"type:decimal(15,2)"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Units []BarangUnit `gorm:"foreignKey:BarangID"`
}

func (Barang) TableName() string {
	return "barang"
}

// GetStockInfo - dapatkan info stock barang
func (b *Barang) GetStockInfo() map[string]int {
	total := len(b.Units)
	siap := 0
	dipinjam := 0
	maintenance := 0
	nonaktif := 0
	
	for _, unit := range b.Units {
		switch unit.Status {
		case UnitStatusSiap, UnitStatusSiapTerbatas:
			siap++
		case UnitStatusDipinjam:
			dipinjam++
		case UnitStatusMaintenanceRequired:
			maintenance++
		case UnitStatusNonAktif:
			nonaktif++
		}
	}
	
	return map[string]int{
		"total":       total,
		"siap":        siap,
		"dipinjam":    dipinjam,
		"maintenance": maintenance,
		"nonaktif":    nonaktif,
	}
}

// IsHabis - cek apakah semua unit tidak tersedia
func (b *Barang) IsHabis() bool {
	stock := b.GetStockInfo()
	return stock["siap"] == 0
}

// GetStockStatus - dapatkan status ketersediaan stock
func (b *Barang) GetStockStatus() string {
	stock := b.GetStockInfo()
	
	if stock["siap"] == 0 {
		return "habis"
	}
	if stock["siap"] <= 2 {
		return "terbatas"
	}
	return "tersedia"
}

type BarangUnit struct {
	ID uint `json:"id" gorm:"primaryKey"`

	BarangID uint   `json:"barang_id" gorm:"not null;index"`
	KodeUnit string `json:"kode_unit" gorm:"unique;not null;index"`
	Kondisi  string `json:"kondisi"`
	Status   string `json:"status" gorm:"default:'siap'"`

	// SiapTerbatasCount int `json:"siap_terbatas_count" gorm:"default:0"`
	FixPriority       int `json:"fix_priority" gorm:"default:0"`

	TahunPerolehan  int    `json:"tahun_perolehan"`
	LokasiPenyimpan string `json:"lokasi_penyimpan"`
	Catatan         string `json:"catatan" gorm:"type:text"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Barang   Barang           `gorm:"foreignKey:BarangID"`
	Komponen []BarangKomponen `gorm:"foreignKey:BarangUnitID"`

	LoanItems []LoanItem `gorm:"many2many:loan_item_units; joinForeignKey:UnitID; joinReferences:LoanItemID"`
}

func (BarangUnit) TableName() string {
	return "barang_unit"
}

// HitungStatus - compute status berdasarkan kondisi unit & komponen
func (u *BarangUnit) HitungStatus() string {
	// 1. Maintenance override
	if u.FixPriority > FixPriorityNormal {
		return UnitStatusMaintenanceRequired
	}
	
	// 2. Cek komponen
	if len(u.Komponen) > 0 {
		for _, k := range u.Komponen {
			if k.Status == KomponenStatusKurang || k.Status == KomponenStatusRusak {
				return UnitStatusSiapTerbatas
			}
		}
	}
	
	// 3. Default siap
	return UnitStatusSiap
}

// BeforeSave - auto-update status jika bukan manual override
func (u *BarangUnit) BeforeSave(tx *gorm.DB) error {
	// Jangan override kalau manually set ke nonaktif atau dipinjam
	if u.Status != UnitStatusNonAktif && u.Status != UnitStatusDipinjam {
		u.Status = u.HitungStatus()
	}
	return nil
}

type BarangKomponen struct {
	ID uint `json:"id" gorm:"primaryKey"`

	BarangUnitID uint   `json:"barang_unit_id" gorm:"not null;index:idx_unit_komponen"`
	Nama         string `json:"nama" gorm:"index:idx_unit_komponen"` // ← composite unique
	JumlahWajib  int    `json:"jumlah_wajib"`
	JumlahAktual int    `json:"jumlah_aktual"`
	Status       string `json:"status"`
	HargaSatuan  int64  `json:"harga_satuan" gorm:"default:0"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (BarangKomponen) TableName() string {
	return "barang_komponen"
}

// HitungStatus - auto-compute status komponen
func (k *BarangKomponen) HitungStatus() string {
	if k.JumlahAktual == 0 {
		return KomponenStatusRusak
	}
	if k.JumlahAktual < k.JumlahWajib {
		return KomponenStatusKurang
	}
	return KomponenStatusLengkap
}

// BeforeSave - auto-update status komponen
func (k *BarangKomponen) BeforeSave(tx *gorm.DB) error {
	k.Status = k.HitungStatus()
	return nil
}

// HitungBiayaPenggantian - hitung biaya penggantian komponen yang kurang/rusak
func (k *BarangKomponen) HitungBiayaPenggantian() int64 {
	kurang := k.JumlahWajib - k.JumlahAktual
	if kurang <= 0 {
		return 0
	}
	return int64(kurang) * k.HargaSatuan
}