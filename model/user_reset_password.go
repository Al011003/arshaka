package model

import "time"

type PasswordReset struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    OTP       string    `gorm:"size:10;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    Used      bool      `gorm:"default:false"`
    CreatedAt time.Time
}
