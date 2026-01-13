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
    
    // Method baru untuk token-based reset
    UpdateResetToken(ctx context.Context, id uint, token string, expiresAt time.Time) error
    GetByResetToken(ctx context.Context, token string) (*model.PasswordReset, error)
    MarkTokenUsed(ctx context.Context, id uint) error
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

func (r *passwordResetRepository) ExpireAllOTPByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.PasswordReset{}).
		Where("user_id = ? AND used = ?", userID, false).
		Update("used", true).Error
}

// Update reset token setelah OTP diverifikasi
func (r *passwordResetRepository) UpdateResetToken(ctx context.Context, id uint, token string, expiresAt time.Time) error {
    return r.db.WithContext(ctx).
        Model(&model.PasswordReset{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "reset_token":       token,
            "token_expires_at":  expiresAt,
            "token_used":        false,
        }).Error
}

// Get password reset record by reset token
func (r *passwordResetRepository) GetByResetToken(ctx context.Context, token string) (*model.PasswordReset, error) {
    var pr model.PasswordReset
    err := r.db.WithContext(ctx).
        Preload("User").
        Where("reset_token = ? AND token_used = ? AND token_expires_at IS NOT NULL AND token_expires_at > ?", 
              token, false, time.Now()).
        First(&pr).Error

    if err != nil {
        return nil, err
    }
    return &pr, nil
}

// Mark reset token as used
func (r *passwordResetRepository) MarkTokenUsed(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).
        Model(&model.PasswordReset{}).
        Where("id = ?", id).
        Update("token_used", true).Error
}