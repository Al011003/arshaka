package model

import "time"

type DeviceToken struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    UserID      uint      `json:"user_id" gorm:"not null"`
    DeviceToken string    `json:"device_token" gorm:"type:text;not null"`
    DeviceType  string    `json:"device_type" gorm:"type:varchar(50);not null"` // android / ios
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
