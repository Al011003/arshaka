package repo

import (
	"backend/model"

	"gorm.io/gorm"
)

type DeviceTokenRepo interface {
    SaveOrUpdate(token *model.DeviceToken) error
    GetByUserID(userID uint) ([]model.DeviceToken, error)
}

type deviceTokenRepo struct {
    db *gorm.DB
}

func NewDeviceTokenRepo(db *gorm.DB) DeviceTokenRepo {
    return &deviceTokenRepo{
        db: db,
    }
}

// SaveOrUpdate: update kalau token sudah pernah ada, insert kalau baru
func (r *deviceTokenRepo) SaveOrUpdate(dt *model.DeviceToken) error {

    var existing model.DeviceToken

    err := r.db.Where("user_id = ? AND device_token = ?", dt.UserID, dt.DeviceToken).
        First(&existing).Error

    // If found → update type + timestamps
    if err == nil {
        existing.DeviceType = dt.DeviceType
        return r.db.Save(&existing).Error
    }

    // If not found → create
    if err == gorm.ErrRecordNotFound {
        return r.db.Create(dt).Error
    }

    return err
}

// Get all tokens for a user
func (r *deviceTokenRepo) GetByUserID(userID uint) ([]model.DeviceToken, error) {
    var tokens []model.DeviceToken
    err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
    return tokens, err
}