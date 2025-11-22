package model

import "time"

type AngkatanMapala struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama      string    `gorm:"type:varchar(100);not null;unique" json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}