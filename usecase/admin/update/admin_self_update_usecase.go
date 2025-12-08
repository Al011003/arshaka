package usecase

import (
	req "backend/dto/request/user"
	apperrors "backend/errors"
	"backend/repo"
)

// ================== Interface ==================
type AdminSelfUpdateUsecase interface {
	AdminUpdateSelf(adminID uint, req req.AdminSelfUpdateRequest) error
}

type adminSelfUpdateUsecase struct {
	userRepo repo.UserRepo
}

func NewAdminSelfUpdateUsecase(r repo.UserRepo) AdminSelfUpdateUsecase {
	return &adminSelfUpdateUsecase{userRepo: r}
}

// ================== Implementasi ==================
func (u *adminSelfUpdateUsecase) AdminUpdateSelf(adminID uint, req req.AdminSelfUpdateRequest) error {
    admin, err := u.userRepo.GetByID(adminID)
    if err != nil {
        return apperrors.ErrInternalServer
    }
    if admin == nil {
        return apperrors.ErrNotFound
    }

    if req.Nama != "" {
        admin.Nama = req.Nama
    }
    if req.NamaLengkap != "" {
        admin.NamaLengkap = req.NamaLengkap
    }
    if req.NRA != "" {
        admin.NRA = req.NRA
    }

    if err := u.userRepo.Update(admin); err != nil {
        return apperrors.ErrInternalServer
    }
    
    return nil
}