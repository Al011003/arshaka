package usecase

import (
	resp "backend/dto/response/superadmin"
	"backend/repo"
	"errors"
)

type AdminProfileUsecase interface {
	GetProfile(superAdminID uint) (*resp.SuperAdminProfileResponse, error)
}

type adminProfileUC struct {
	userRepo repo.UserRepo
}


func NewAdminProfileUC(r repo.UserRepo) AdminProfileUsecase {
	return &adminProfileUC{userRepo: r}
}

func (u *adminProfileUC) GetProfile(AdminID uint) (*resp.SuperAdminProfileResponse, error) {
	admin, err := u.userRepo.GetByID(AdminID)
	if err != nil || admin == nil {
		return nil, errors.New("admin tidak ditemukan")
	}

	if admin.Role != "admin" {
		return nil, errors.New("hanya superadmin yang dapat mengakses data ini")
	}

	resp := &resp.SuperAdminProfileResponse{
		NRA:        admin.NRA,
		Nama:       admin.Nama,
		NamaLengkap: admin.NamaLengkap,
		Status:     admin.Status, // bisa juga pakai field khusus status
	}

	return resp, nil
}
