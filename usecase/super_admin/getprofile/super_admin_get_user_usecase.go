package usecase

import (
	request "backend/dto/request/user"
	resp "backend/dto/response/admin"
	"backend/repo"
	"errors"
)

type SuperAdminGetUserUsecase interface {
  GetUsers(superAdminID uint, filter request.UserFilter) ([]resp.UserListResponse, int, error)
}

type superAdminGetUserUC struct {
    userRepo repo.UserRepo
}

func NewSuperAdminGetUserUsecase(r repo.UserRepo) SuperAdminGetUserUsecase {
    return &superAdminGetUserUC{userRepo: r}
}

func (u *superAdminGetUserUC) GetUsers(superAdminID uint, filter request.UserFilter) ([]resp.UserListResponse, int, error) {

	// Validasi admin
	currentUser, err := u.userRepo.GetByID(superAdminID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	if currentUser.Role != "superadmin" {
		return nil, 0, errors.New("forbidden: only superadmin can access this")
	}

	// Exclude admin & superadmin
	filter.ExcludeRole = "superadmin"

	// Ambil data user dari repo
	users, total, err := u.userRepo.FindUsers(filter)
	if err != nil {
		return nil, 0, err
	}

	// Mapping ke response minimalis
	result := make([]resp.UserListResponse, 0, len(users))
	for _, u := range users {
		result = append(result, resp.UserListResponse{
			ID:             u.ID,
			Nama:           u.Nama,
			AngkatanMapala: u.AngkatanMapala,
			Status:         u.Status,
			Role:           u.Role,
		})
	}

	return result, total, nil
}