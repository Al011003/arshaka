package usecase

import (
	req "backend/dto/request/user"
	res "backend/dto/response/user"
	"backend/repo"
	"errors"
)

type AdminUpdateUserUsecase interface {
	AdminUpdateUser(userID uint, req req.AdminUpdateUserRequest, adminRole string) (*res.UpdateUserResponse, error)
}

type adminUpdateUserUsecase struct {
	authRepo     repo.AuthRepository
	fakultasRepo repo.FakultasRepo
	jurusanRepo  repo.JurusanRepo
	angkatanRepo repo.AngkatanMapalaRepo
}

func NewAdminUpdateUserUsecase(
	a repo.AuthRepository,
	f repo.FakultasRepo,
	j repo.JurusanRepo,
	am repo.AngkatanMapalaRepo,
) AdminUpdateUserUsecase {
	return &adminUpdateUserUsecase{
		authRepo:     a,
		fakultasRepo: f,
		jurusanRepo:  j,
		angkatanRepo: am,
	}
}

func (u *adminUpdateUserUsecase) AdminUpdateUser(
	userID uint,
	req req.AdminUpdateUserRequest,
	adminRole string,
) (*res.UpdateUserResponse, error) {

	targetUser, err := u.authRepo.FindUserByID(userID)
    if err != nil || targetUser == nil {
        return nil, errors.New("user tidak ditemukan")
    }


	// ---------------------------------------------------------
	// 1. Validasi Role Admin (sama persis gaya register)
	// ---------------------------------------------------------
if adminRole == "admin" && targetUser.Role != "user" {
    return nil, errors.New("admin hanya boleh mengedit user saja")
}

// Superadmin tidak boleh update superadmin lain
if adminRole == "superadmin" && targetUser.Role == "superadmin" {
    return nil, errors.New("superadmin tidak boleh mengedit superadmin lain")
}

	// ---------------------------------------------------------
	// 2. Ambil user
	// ---------------------------------------------------------

	user, err := u.authRepo.FindUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// ---------------------------------------------------------
	// 3. Validasi NRA (kalau diubah)
	// ---------------------------------------------------------

	if req.NRA != "" && req.NRA != user.NRA {
		exist, _ := u.authRepo.FindUserByNra(req.NRA)
		if exist != nil {
			return nil, errors.New("NRA sudah terdaftar")
		}
		user.NRA = req.NRA
	}

	// ---------------------------------------------------------
	// 4. Validasi Fakultas dan Jurusan (sama dengan register)
	// ---------------------------------------------------------

	if req.Fakultas != "" {

		fakultas, err := u.fakultasRepo.GetByName(req.Fakultas)
		if err != nil || fakultas == nil {
			return nil, errors.New("fakultas tidak ditemukan")
		}

		if req.Jurusan != "" {
			jurusan, err := u.jurusanRepo.GetByName(req.Jurusan)
			if err != nil || jurusan == nil {
				return nil, errors.New("jurusan tidak ditemukan")
			}

			if jurusan.FakultasID != fakultas.ID {
				return nil, errors.New("jurusan tidak sesuai dengan fakultas")
			}

			user.Jurusan = req.Jurusan
		}

		user.Fakultas = req.Fakultas
	}

	// ---------------------------------------------------------
	// 5. Validasi Angkatan Mapala (sama dengan register)
	// ---------------------------------------------------------

	if req.AngkatanMapala != "" {
		angkatan, err := u.angkatanRepo.GetByName(req.AngkatanMapala)
		if err != nil || angkatan == nil {
			return nil, errors.New("angkatan mapala tidak ditemukan")
		}
		user.AngkatanMapala = req.AngkatanMapala
	}

	// ---------------------------------------------------------
	// 6. Update field lain (safe update)
	// ---------------------------------------------------------

	if req.Nama != "" {
		user.Nama = req.Nama
	}

	if req.NamaLengkap != "" {
		user.NamaLengkap = req.NamaLengkap
	}

	if req.NoHP != "" {
		user.NoHP = req.NoHP
	}


	if req.AngkatanKampus != "" {
		user.AngkatanKampus = req.AngkatanKampus
	}

	if req.NIM != "" {
		user.NIM = req.NIM
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	// ---------------------------------------------------------
	// 7. Simpan perubahan
	// ---------------------------------------------------------

	if err := u.authRepo.UpdateUser(user); err != nil {
		return nil, errors.New("gagal menyimpan perubahan user")
	}

	// ---------------------------------------------------------
	// 8. Response (style register)
	// ---------------------------------------------------------

	resp := &res.UpdateUserResponse{
		ID:     user.ID,
		NRA:    user.NRA,
		Role:   user.Role,
		Status: "updated",
	}

	return resp, nil
}
