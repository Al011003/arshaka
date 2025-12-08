package repo

import (
	"context"
	"time"

	"backend/model"

	"gorm.io/gorm"
)

type PasswordResetRepository interface {
    Create(ctx context.Context, pr *model.PasswordReset) error
    GetValidOTP(ctx context.Context, userID uint, otp string) (*model.PasswordReset, error)
    MarkUsed(ctx context.Context, id uint) error

	GetLatestOTPByUserID(ctx context.Context, userID uint) (*model.PasswordReset, error)
	ExpireAllOTPByUserID(ctx context.Context, userID uint) error
}

type passwordResetRepository struct {
    db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
    return &passwordResetRepository{db}
}

func (r *passwordResetRepository) Create(ctx context.Context, pr *model.PasswordReset) error {
    return r.db.WithContext(ctx).Create(pr).Error
}

func (r *passwordResetRepository) GetValidOTP(ctx context.Context, userID uint, otp string) (*model.PasswordReset, error) {
    var pr model.PasswordReset
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND otp = ? AND used = ? AND expires_at > ?", userID, otp, false, time.Now()).
        First(&pr).Error

    if err != nil {
        return nil, err
    }
    return &pr, nil
}

func (r *passwordResetRepository) MarkUsed(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).
        Model(&model.PasswordReset{}).
        Where("id = ?", id).
        Update("used", true).Error
}

func (r *passwordResetRepository) GetLatestOTPByUserID(ctx context.Context, userID uint) (*model.PasswordReset, error) {
	var pr model.PasswordReset
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&pr).Error
	
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// METHOD BARU: Expire semua OTP lama milik user
func (r *passwordResetRepository) ExpireAllOTPByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.PasswordReset{}).
		Where("user_id = ? AND used = ?", userID, false).
		Update("used", true).Error
}