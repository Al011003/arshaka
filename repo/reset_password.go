package repo

import (
	"backend/model"
	"time"

	"gorm.io/gorm"
)

type PasswordResetRepo interface {
	Create(userID uint) error
	GetPendingByUserID(userID uint) (*model.PasswordResetRequest, error)
	CancelRequest(requestID uint, adminID uint, canceledAt time.Time) error
	ApproveRequest(requestID uint, adminID uint, approvedAt time.Time) error
	GetAllFiltered(status string) ([]model.PasswordResetRequest, error)
	GetByID(id uint) (*model.PasswordResetRequest, error)
}

type passwordResetRepo struct {
	DB *gorm.DB
}

func NewPasswordResetRepo(db *gorm.DB) PasswordResetRepo {
	return &passwordResetRepo{DB: db}
}

func (r *passwordResetRepo) Create(userID uint) error {
	req := model.PasswordResetRequest{
		UserID:      userID,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	return r.DB.Create(&req).Error
}

func (r *passwordResetRepo) GetPendingByUserID(userID uint) (*model.PasswordResetRequest, error) {
	var req model.PasswordResetRequest
	err := r.DB.Where("user_id = ? AND status = 'pending'", userID).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *passwordResetRepo) CancelRequest(requestID uint, adminID uint, canceledAt time.Time) error {
	return r.DB.Model(&model.PasswordResetRequest{}).
		Where("id = ?", requestID).
		Updates(map[string]interface{}{
			"status":      "canceled",
			"canceled_at": canceledAt,
		}).Error
}

func (r *passwordResetRepo) ApproveRequest(requestID uint, adminID uint, approvedAt time.Time) error {
	return r.DB.Model(&model.PasswordResetRequest{}).
		Where("id = ?", requestID).
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_at": approvedAt,
		}).Error
}

func (r *passwordResetRepo) GetAllFiltered(status string) ([]model.PasswordResetRequest, error) {
	var reqs []model.PasswordResetRequest
	db := r.DB

	if status != "" && status != "all" {
		db = db.Where("status = ?", status)
	}

	err := db.Order("requested_at DESC").Find(&reqs).Error
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (r *passwordResetRepo) GetByID(id uint) (*model.PasswordResetRequest, error) {
	var req model.PasswordResetRequest
	err := r.DB.First(&req, id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

