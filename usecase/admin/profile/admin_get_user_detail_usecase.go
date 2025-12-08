package usecase

import (
	res "backend/dto/response/admin"
	"backend/repo"
	"errors"
)
type AdminUserDetailUsecase interface {
    GetDetail(userID uint) (interface{}, error)
}


type adminUserDetailUC struct {
    userRepo repo.UserRepo
}

func NewAdminUserDetailUsecase(r repo.UserRepo) AdminUserDetailUsecase {
    return &adminUserDetailUC{userRepo: r}
}

func (uc *adminUserDetailUC) GetDetail(userID uint) (interface{}, error) {
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
        res := &res.UserDetailResponse{
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
        return nil, errors.New("tidak bisa melihat profile sesama admin")
    }

    return nil, errors.New("tidak bisa melihat super admin")
}
