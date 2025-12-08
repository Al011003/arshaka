package usecase

import (
	req "backend/dto/request/auth"
	res "backend/dto/response/auth"
	"backend/model"
	"backend/repo"
	"backend/utils"
	"errors"
)

// ================= Interface =================
type RegisterUsecase interface {
	RegisterUser(req req.RegisterUserRequest) (*res.RegisterResponse, error)
	RegisterAdmin(req req.RegisterAdminRequest) (*res.RegisterResponse, error)
}

type registerUsecase struct {
	authRepo     repo.AuthRepository
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

// ================= Register User (anggota) =================
func (u *registerUsecase) RegisterUser(req req.RegisterUserRequest) (*res.RegisterResponse, error) {

	// Validasi NRA unik
	existing, _ := u.authRepo.FindUserByNra(req.NRA)
	if existing != nil {
		return nil, errors.New("NRA sudah terdaftar")
	}

	// Validasi fakultas
	fakultas, err := u.fakultasRepo.GetByName(req.Fakultas)
	if err != nil || fakultas == nil {
		return nil, errors.New("fakultas tidak ditemukan")
	}

	jurusan, err := u.jurusanRepo.GetByName(req.Jurusan)
	if err != nil || jurusan == nil {
		return nil, errors.New("jurusan tidak ditemukan")
	}
	if jurusan.FakultasID != fakultas.ID {
		return nil, errors.New("jurusan tidak sesuai dengan fakultas")
	}

	// Validasi angkatan
	angkatan, err := u.angkatanRepo.GetByName(req.AngkatanMapala)
	if err != nil || angkatan == nil {
		return nil, errors.New("angkatan mapala tidak ditemukan")
	}

	// Generate password default
	hashed, err := utils.HashPassword("mapala" + req.NRA)
	if err != nil {
		return nil, errors.New("gagal generate password")
	}

	// Buat user model
	user := model.User{
		Nama:              req.Nama,
		NamaLengkap:       req.NamaLengkap,
		NRA:               req.NRA,
		Role:              "user",
		Password:          hashed,
		IsDefaultPassword: true,
		Email:             nil,        // set email ke nil
		IsEmailNull:       true,       // set flag email null
		AngkatanMapala:    req.AngkatanMapala,
		AngkatanKampus:    req.AngkatanKampus,
		NIM:               req.NIM,
		Jurusan:           req.Jurusan,
		Fakultas:          req.Fakultas,
		NoHP:              req.NoHP,
	}

	if err := u.authRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &res.RegisterResponse{
		NRA:  user.NRA,
		Role: user.Role,
	}, nil
}

// ================= Register Admin / Superadmin =================
func (u *registerUsecase) RegisterAdmin(req req.RegisterAdminRequest) (*res.RegisterResponse, error) {
	
	existing, _ := u.authRepo.FindUserByNra(req.NRA)
	if existing != nil {
		return nil, errors.New("NRA sudah terdaftar")
	}
	
	// generate password
	hashed, err := utils.HashPassword("mapala" + req.NRA)
	if err != nil {
		return nil, errors.New("gagal generate password")
	}

	user := model.User{
		Nama:              req.Nama,
		NamaLengkap:       req.NamaLengkap,
		NRA:               req.NRA,
		Role:              "admin", // default role
		Password:          hashed,
		IsDefaultPassword: true,
		Email:             nil,     // set email ke nil
		IsEmailNull:       true,    // set flag email null
	}

	if err := u.authRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &res.RegisterResponse{
		NRA:  user.NRA,
		Role: user.Role,
	}, nil
}