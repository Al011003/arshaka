package model

import "time"

type Fakultas struct {
    ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
    Nama      string     `gorm:"type:varchar(100);not null;unique" json:"nama"`
    Jurusan   []Jurusan  `gorm:"foreignKey:FakultasID" json:"jurusan,omitempty"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}