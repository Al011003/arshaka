package model

import "time"

type PasswordResetRequest struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      
	Status      string    // pending, approved, canceled, rejected
	RequestedAt time.Time
	CanceledAt  *time.Time
	ApprovedAt  *time.Time
}
