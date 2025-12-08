package usecase

import (
	req "backend/dto/request/user"
	"backend/repo"
	"context"
	"fmt"
)

type UpdateEmailUsecase interface {
	UpdateEmail(ctx context.Context, userID uint, req req.UpdateEmailRequest) error
}

type updateEmailUsecase struct {
	userRepo repo.UserRepo
}

func NewUpdateEmailUsecase(userRepo repo.UserRepo) UpdateEmailUsecase {
	return &updateEmailUsecase{
		userRepo: userRepo,
	}
}

func (u *updateEmailUsecase) UpdateEmail(ctx context.Context, userID uint, req req.UpdateEmailRequest) error {

	// 1. Cek apakah email sudah digunakan user lain
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)

	// Kalau error bukan record not found => DB error
	if err != nil && err.Error() != "record not found" {
		return fmt.Errorf("gagal mengecek email: %w", err)
	}

	// Kalau ketemu user lain dengan email itu -> tidak boleh
	if err == nil && existingUser.ID != userID {
		return fmt.Errorf("email sudah digunakan")
	}

	// 2. Update email user
	if err := u.userRepo.UpdateEmail(ctx, userID, req.Email); err != nil {
		return fmt.Errorf("gagal mengupdate email: %w", err)
	}

	return nil
}

