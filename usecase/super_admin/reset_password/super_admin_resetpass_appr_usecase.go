package usecase

import (
	"backend/model"
	"backend/repo"
	"backend/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SuperAdminAccResetPasswordUsecase interface {
	Approve(resetID uint, adminID uint) error
	GetAllRequestsFiltered(status string) ([]model.PasswordResetRequest, error)
	CancelResetByIdentity(resetID uint, adminID uint) error
}

type superAdminAccResetPasswordUsecase struct {
	resetRepo repo.PasswordResetRepo
	userRepo  repo.UserRepo
}

func NewSuperAdminAccResetPasswordUsecase(
	resetRepo repo.PasswordResetRepo,
	userRepo repo.UserRepo,
) SuperAdminAccResetPasswordUsecase {
	return &superAdminAccResetPasswordUsecase{
		resetRepo: resetRepo,
		userRepo:  userRepo,
	}
}

// ======================================================================
// APPROVE REQUEST RESET PASSWORD
// ======================================================================
func (uc *superAdminAccResetPasswordUsecase) Approve(resetID uint, adminID uint) error {

	// 1. Ambil request berdasarkan resetID
	req, err := uc.resetRepo.GetByID(resetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("request tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil data request: %w", err)
	}

	// 2. Pastikan status masih pending
	if req.Status != "pending" {
		return errors.New("request ini sudah tidak pending")
	}

	// 3. Reset password user ke password default
	hashed, err := utils.HashPassword("adminarshaka123")
	if err != nil {
		return fmt.Errorf("gagal meng-hash password default: %w", err)
	}

	// Update password dan set is_default_password = true
	if err := uc.userRepo.ResetPassword(req.UserID, hashed, true); err != nil {
		return fmt.Errorf("gagal memperbarui password user: %w", err)
	}

	// 4. Tandai request sebagai approved
	if err := uc.resetRepo.ApproveRequest(resetID, adminID, time.Now()); err != nil {
		return fmt.Errorf("gagal meng-approve request: %w", err)
	}

	return nil
}

// ======================================================================
// GET ALL RESET REQUESTS (OPTIONAL FILTER BY STATUS)
// ======================================================================
func (uc *superAdminAccResetPasswordUsecase) GetAllRequestsFiltered(status string) ([]model.PasswordResetRequest, error) {
	reqs, err := uc.resetRepo.GetAllFiltered(status)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data reset request: %w", err)
	}
	return reqs, nil
}

func (uc *superAdminAccResetPasswordUsecase) CancelResetByIdentity(resetID uint, adminID uint) error {
	req, err := uc.resetRepo.GetByID(resetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("request tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil data request: %w", err)
	}

	if req.Status != "pending" {
		return errors.New("request ini sudah tidak pending")
	}
		if err := uc.resetRepo.CancelRequest(resetID, adminID, time.Now()); err != nil {
		return fmt.Errorf("gagal meng-cancelrequest: %w", err)
	}

	return nil
}
