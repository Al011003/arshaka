package usecase

import (
	"backend/repo"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AdminResetPasswordUsecase interface {
	RequestReset(userID uint) error
	CancelReset(userID uint) error
}

type adminResetUC struct {
	resetRepo repo.PasswordResetRepo
	userRepo  repo.UserRepo
}

func NewAdminResetPasswordUsecase(r repo.PasswordResetRepo, u repo.UserRepo) AdminResetPasswordUsecase {
	return &adminResetUC{
		resetRepo: r,
		userRepo:  u,
	}
}

func (uc *adminResetUC) RequestReset(userID uint) error {
	existing, err := uc.resetRepo.GetPendingByUserID(userID)

	if err == nil && existing != nil {
		return errors.New("masih ada permintaan reset password yang pending")
	}

	// Error lain selain not found → error beneran
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("gagal mengecek permintaan reset: %w", err)
	}

	return uc.resetRepo.Create(userID)
}

func (uc *adminResetUC) CancelReset(userID uint) error {
	req, err := uc.resetRepo.GetPendingByUserID(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tidak ada permintaan reset yang pending")
		}
		return fmt.Errorf("gagal mengambil data reset: %w", err)
	}

	return uc.resetRepo.CancelRequest(req.ID, time.Now())
}
