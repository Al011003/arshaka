package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Nama           string         `json:"nama"`
	NamaLengkap    string         `json:"nama_lengkap"`
	NRA            string         `json:"nra"`                 // Nomor Anggota Mapala
	AngkatanMapala string         `json:"angkatan_mapala"`     // Angkatan Mapala
	AngkatanKampus string         `json:"angkatan_kampus"`     // 22/23/24 dst
	NIM            string         `json:"nim"`
	Jurusan        string         `json:"jurusan"`
	Fakultas       string         `json:"fakultas"`
	NoHP           string         `json:"no_hp"`
	Role           string         `json:"role"`                // "superadmin", "admin", "user"
	Status string `gorm:"type:enum('aktif','lulus','keluar');default:'aktif'" json:"status"`


	Email        *string `gorm:"unique;default:null" json:"email"` // nullable + unique
	IsEmailNull  bool    `gorm:"default:true" json:"is_email_null"`


	Password       string         `json:"-"`   
	IsDefaultPassword bool 			`json:"-"`               
	FotoURL        string         `json:"foto_url"`            

	RefreshToken   string         `json:"-"`                   

	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
