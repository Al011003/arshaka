package model

import "time"

type PasswordResetRequest struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null;index"`
	Status      string     `gorm:"type:varchar(20);not null"` // pending, approved, canceled, rejected
	RequestedAt time.Time  `gorm:"autoCreateTime"`
	ApprovedAt  *time.Time
	CanceledAt  *time.Time
}

