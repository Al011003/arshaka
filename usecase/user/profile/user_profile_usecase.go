package usecase

import (
	res "backend/dto/response/user"
	"backend/repo"
	"errors"
)

type UserProfileUsecase interface {
	GetProfile(userID uint) (*res.UserDetailResponse, error)
}

type userProfileUC struct {
	userRepo repo.UserRepo
}


func NewUserProfileUsecase(r repo.UserRepo) UserProfileUsecase {
	return &userProfileUC{userRepo: r}
}

func (u *userProfileUC) GetProfile(userID uint) (*res.UserDetailResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user tidak ditemukan")
	}


	resp := &res.UserDetailResponse{
		FotoURL:             user.FotoURL,
        Nama:           user.Nama,
        NamaLengkap:    user.NamaLengkap,
        NRA:            user.NRA,
        Jurusan:        user.Jurusan,
        Fakultas:       user.Fakultas,
        AngkatanMapala: user.AngkatanMapala,
        AngkatanKampus: user.AngkatanKampus,
        NIM:            user.NIM,
        Status:         user.Status,
	}

	return resp, nil
}
