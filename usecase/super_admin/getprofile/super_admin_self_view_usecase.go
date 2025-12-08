package usecase

import (
	superadminresp "backend/dto/response/superadmin"
	"backend/repo"
	"errors"
)

type SuperAdminProfileUsecase interface {
	GetProfile(superAdminID uint) (*superadminresp.SuperAdminProfileResponse, error)
}

type superAdminProfileUC struct {
	userRepo repo.UserRepo
}


func NewSuperAdminProfileUC(r repo.UserRepo) SuperAdminProfileUsecase {
	return &superAdminProfileUC{userRepo: r}
}

func (u *superAdminProfileUC) GetProfile(superAdminID uint) (*superadminresp.SuperAdminProfileResponse, error) {
	superadmin, err := u.userRepo.GetByID(superAdminID)
	if err != nil || superadmin == nil {
		return nil, errors.New("superadmin tidak ditemukan")
	}

	if superadmin.Role != "superadmin" {
		return nil, errors.New("hanya superadmin yang dapat mengakses data ini")
	}

	resp := &superadminresp.SuperAdminProfileResponse{
		NRA:        superadmin.NRA,
		Nama:       superadmin.Nama,
		NamaLengkap: superadmin.NamaLengkap,
		Status:     superadmin.Status, // bisa juga pakai field khusus status
	}

	return resp, nil
}
