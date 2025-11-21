package usecase

import (
	req "backend/dto/request/auth"
	res "backend/dto/response/auth"
	"backend/model"
	"backend/repo"
	"backend/utils"
	"errors"
)

type RegisterUsecase interface {
    RegisterUser(req req.RegisterRequest, creatorRole string) (*res.RegisterResponse, error)
}

type registerUsecase struct {
    repo repo.AuthRepository
}

func NewRegisterUsecase(r repo.AuthRepository) RegisterUsecase {
    return &registerUsecase{
        repo: r,
    }
}

func (u *registerUsecase) RegisterUser(req req.RegisterRequest, creatorRole string) (*res.RegisterResponse, error) {
// 1. cek hak akses creator
// Admin hanya boleh membuat user
if creatorRole == "admin" {
    if req.Role != "user" {
        return nil, errors.New("admin hanya boleh membuat user")
    }
}

// Only superadmin boleh membuat admin atau superadmin
if creatorRole != "superadmin" && req.Role != "user" {
    return nil, errors.New("hanya superadmin yang boleh membuat admin atau superadmin")
}

// Kalau bukan admin atau superadmin, tolak semua
if creatorRole != "superadmin" && creatorRole != "admin" {
    return nil, errors.New("tidak punya hak buat register user")
}

if creatorRole == "superadmin" && req.Role == "superadmin" {
    return nil, errors.New("superadmin tidak boleh membuat superadmin lain")
}
    // 2. cek NRA unik
    existing, _ := u.repo.FindUserByNra(req.NRA)
    if existing != nil {
        return nil, errors.New("NRA sudah terdaftar")
    }

    // 3. generate password default
    autoPassword := "mapala" + req.NRA
    hashed, err := utils.HashPassword(autoPassword)
    if err != nil {
        return nil, errors.New("gagal generate password")
    }

    // 4. create user model
    user := model.User{
        Nama:             req.Nama,
        NamaLengkap:      req.NamaLengkap,
        NRA:              req.NRA,
        Role:             req.Role,
        Password:         hashed,
        IsDefaultPassword: true,
        AngkatanMapala:   req.AngkatanMapala,
        AngkatanKampus:   req.AngkatanKampus,
        NIM:              req.NIM,
        Jurusan:          req.Jurusan,
        Fakultas:         req.Fakultas,
        NoHP:             req.NoHP,
    }

    // 5. simpan ke DB
    if err := u.repo.CreateUser(&user); err != nil {
        return nil, err
    }

    // 6. mapping ke response DTO
    resp := &res.RegisterResponse{
        NRA:             user.NRA,
        Role:            user.Role,
    }

    return resp, nil
}