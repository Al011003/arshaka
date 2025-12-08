package usecase

import (
	resA "backend/dto/response/admin"
	"backend/repo"
	"errors"
)
type UserDetailUsecase interface {
    GetDetail(userID uint) (interface{}, error)
}


type userDetailUC struct {
    userRepo repo.UserRepo
}

func NewUserDetailUsecase(r repo.UserRepo) UserDetailUsecase {
    return &userDetailUC{userRepo: r}
}

func (uc *userDetailUC) GetDetail(userID uint) (interface{}, error) {
    // Ambil data user
    user, err := uc.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }

    if user == nil {
        return nil, errors.New("user not found")
    }

    // Cek role target user
    switch user.Role {
    case "user":
        // Return detail lengkap
        res := &resA.UserDetailResponse{
            ID:             user.ID,
            Nama:           user.Nama,
            NamaLengkap:    user.NamaLengkap,
            NRA:            user.NRA,
            Jurusan:        user.Jurusan,
            Fakultas:       user.Fakultas,
            AngkatanMapala: user.AngkatanMapala,
            AngkatanKampus: user.AngkatanKampus,
            NIM:            user.NIM,
            Status:         user.Status,
            CreatedAt:      user.CreatedAt,
        }
        return res, nil

    case "admin":
        // Return detail terbatas untuk admin
        res := &resA.AdminDetailResponse{
            ID:          user.ID,
            Nama:        user.Nama,
            NamaLengkap: user.NamaLengkap,
            NRA:         user.NRA,
            Role:        user.Role,
        }
        return res, nil
    }

    return nil, errors.New("tidak bisa melihat super admin")
}
