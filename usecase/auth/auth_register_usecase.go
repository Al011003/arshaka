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
	authRepo    repo.AuthRepository
	fakultasRepo repo.FakultasRepo
	jurusanRepo  repo.JurusanRepo
	angkatanRepo repo.AngkatanMapalaRepo
}

func NewRegisterUsecase(
	a repo.AuthRepository,
	f repo.FakultasRepo,
	j repo.JurusanRepo,
	am repo.AngkatanMapalaRepo,
) RegisterUsecase {
	return &registerUsecase{
		authRepo:     a,
		fakultasRepo: f,
		jurusanRepo:  j,
		angkatanRepo: am,
	}
}

func (u *registerUsecase) RegisterUser(req req.RegisterRequest, creatorRole string) (*res.RegisterResponse, error) {

	// ---------------------------------------------------------
	// 1. Validasi Role Creator
	// ---------------------------------------------------------

	// Admin hanya boleh membuat user
	if creatorRole == "admin" && req.Role != "user" {
		return nil, errors.New("admin hanya boleh membuat user")
	}

	// Selain superadmin tidak bisa buat admin/superadmin
	if creatorRole != "superadmin" && req.Role != "user" {
		return nil, errors.New("hanya superadmin yang boleh membuat admin atau superadmin")
	}

	// Jika bukan admin / superadmin
	if creatorRole != "admin" && creatorRole != "superadmin" {
		return nil, errors.New("tidak punya hak membuat user")
	}

	// Superadmin tidak boleh buat superadmin lain
	if creatorRole == "superadmin" && req.Role == "superadmin" {
		return nil, errors.New("superadmin tidak boleh membuat superadmin lain")
	}

	// ---------------------------------------------------------
	// 2. Validasi NRA unik
	// ---------------------------------------------------------

	existing, _ := u.authRepo.FindUserByNra(req.NRA)
	if existing != nil {
		return nil, errors.New("NRA sudah terdaftar")
	}

	// ---------------------------------------------------------
	// 3. Validasi Fakultas
	// ---------------------------------------------------------

	fakultas, err := u.fakultasRepo.GetByName(req.Fakultas)
    if err != nil || fakultas == nil {
        return nil, errors.New("fakultas tidak ditemukan")
    }

    // Cek jurusan
    jurusan, err := u.jurusanRepo.GetByName(req.Jurusan)
    if err != nil || jurusan == nil {
        return nil, errors.New("jurusan tidak ditemukan")
    }

    if jurusan.FakultasID != fakultas.ID {
        return nil, errors.New("jurusan tidak sesuai dengan fakultas")
    }

    // Cek angkatan mapala
    angkatan, err := u.angkatanRepo.GetByName(req.AngkatanMapala)
    if err != nil || angkatan == nil {
        return nil, errors.New("angkatan mapala tidak ditemukan")
    }

	// ---------------------------------------------------------
	// 6. Generate Password Default
	// ---------------------------------------------------------

	autoPassword := "mapala" + req.NRA
	hashed, err := utils.HashPassword(autoPassword)
	if err != nil {
		return nil, errors.New("gagal generate password")
	}

	// ---------------------------------------------------------
	// 7. Buat User Model
	// ---------------------------------------------------------

	user := model.User{
		Nama:              req.Nama,
		NamaLengkap:       req.NamaLengkap,
		NRA:               req.NRA,
		Role:              req.Role,
		Password:          hashed,
		IsDefaultPassword: true,
		AngkatanMapala:  req.AngkatanMapala,
		AngkatanKampus:    req.AngkatanKampus,
		NIM:               req.NIM,
		Jurusan:         req.Jurusan,
		Fakultas:        req.Fakultas,
		NoHP:              req.NoHP,
	}

	// ---------------------------------------------------------
	// 8. Simpan User
	// ---------------------------------------------------------

	if err := u.authRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	// ---------------------------------------------------------
	// 9. Response
	// ---------------------------------------------------------

	resp := &res.RegisterResponse{
		NRA:  user.NRA,
		Role: user.Role,
	}

	return resp, nil
}
