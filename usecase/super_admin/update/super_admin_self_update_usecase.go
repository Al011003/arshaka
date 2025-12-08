package usecase

import (
	req "backend/dto/request/user"
	"backend/repo"
	"errors"
)

// ================== Interface ==================
type SuperAdminSelfUpdateUsecase interface {
	SuperAdminUpdateSelf(superAdminID uint, req req.SuperAdminSelfUpdateRequest) error
}

type superAdminSelfUpdateUsecase struct {
	userRepo repo.UserRepo
}

func NewSuperAdminSelfUpdateUsecase(r repo.UserRepo) SuperAdminSelfUpdateUsecase {
	return &superAdminSelfUpdateUsecase{userRepo: r}
}

// ================== Implementasi ==================
func (u *superAdminSelfUpdateUsecase) SuperAdminUpdateSelf(superAdminID uint, req req.SuperAdminSelfUpdateRequest) error {

	superadmin, err := u.userRepo.GetByID(superAdminID)
	if err != nil || superadmin == nil {
		return errors.New("superadmin tidak ditemukan")
	}

	if superadmin.Role != "superadmin" {
		return errors.New("hanya superadmin yang dapat mengupdate data ini")
	}

	// Hanya update field yang divalidasi (Nama, NamaLengkap, NRA)
	if req.Nama != "" {
		superadmin.Nama = req.Nama
	}
	if req.NamaLengkap != "" {
		superadmin.NamaLengkap = req.NamaLengkap
	}
	if req.NRA != "" {
		superadmin.NRA = req.NRA
	}


	return u.userRepo.Update(superadmin)
}
