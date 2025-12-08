package usecase

import (
	req "backend/dto/request/auth"
	"backend/repo"
	"backend/utils"
	"context"
	"fmt"
)

type UpdatePasswordUsecase interface {
	UpdatePassword(ctx context.Context, userID uint, req req.UpdatePasswordRequest) error
}

type updatePasswordUsecase struct {
	repo       repo.UserRepo
}

func NewUpdatePasswordUsecase(repo repo.UserRepo) UpdatePasswordUsecase {
    return &updatePasswordUsecase{
        repo: repo,
    }
}


func (u *updatePasswordUsecase) UpdatePassword(ctx context.Context, userID uint, req req.UpdatePasswordRequest) error {

    // 1. Ambil user
    user, err := u.repo.GetByID(userID)
    if err != nil {
        return fmt.Errorf("gagal mengambil user: %w", err)
    }

    // 2. Cek password lama (false = salah)
    if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
        return fmt.Errorf("password lama salah")
    }

    // 3. Hash password baru
    hashedPassword, err := utils.HashPassword(req.NewPassword)
    if err != nil {
        return fmt.Errorf("gagal mengenkripsi password baru: %w", err)
    }

    // 4. Update password baru
    if err := u.repo.UpdatePassword(user.ID, hashedPassword, false); err != nil {
        return fmt.Errorf("gagal update password: %w", err)
    }

    return nil
}
