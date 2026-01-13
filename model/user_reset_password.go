package model

import "time"

type PasswordReset struct {
    ID         uint       `gorm:"primaryKey"`
    UserID     uint       `gorm:"not null"`
    User       User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    OTP        string     `gorm:"size:10;not null"`
    ResetToken string     `gorm:"size:255;index"` // Token untuk reset password
    ExpiresAt  time.Time  `gorm:"not null"`
    TokenExpiresAt *time.Time // Expiry untuk reset token (nullable)
    Used       bool       `gorm:"default:false"`
    TokenUsed  bool       `gorm:"default:false"` // Track apakah token sudah dipakai
    CreatedAt  time.Time
}