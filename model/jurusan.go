package model

import "time"

type Jurusan struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Nama        string    `gorm:"type:varchar(100);not null" json:"nama"`
    FakultasID  uint      `gorm:"not null" json:"fakultas_id"`
    Fakultas    Fakultas  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
