package usecase

import (
	"backend/repo"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AdminResetPasswordUsecase interface {
	RequestResetByIdentity(nra string, name string) error
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

func (uc *adminResetUC) RequestResetByIdentity(nra string, nama string) error {
	user, err := uc.userRepo.FindByNRAAndName(nra, nama)
	if err != nil {
		return errors.New("permintaan reset tidak dapat diproses")
	}

	// ❌ BLOCK SUPER ADMIN
	if user.Role == "super_admin" {
		return errors.New("permintaan reset tidak dapat diproses")
	}

	// ❌ BLOCK ROLE NON-ADMIN
	if user.Role != "admin" {
		return errors.New("permintaan reset tidak dapat diproses")
	}

	existing, err := uc.resetRepo.GetPendingByUserID(user.ID)
	if err == nil && existing != nil {
		return errors.New("permintaan reset masih diproses")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("gagal mengecek permintaan reset: %w", err)
	}

	return uc.resetRepo.Create(user.ID)
}
