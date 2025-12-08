package usecase

import (
	request "backend/dto/request/user"
	resp "backend/dto/response/admin"
	"backend/repo"
	"errors"
)

type AdminGetUserUsecase interface {
	GetUsers(adminID uint, filter request.UserFilter) ([]resp.UserListResponse, int, error)
}

type adminGetUserUC struct {
	userRepo repo.UserRepo
}

func NewAdminGetUserUsecase(r repo.UserRepo) AdminGetUserUsecase {
	return &adminGetUserUC{userRepo: r}
}

func (u *adminGetUserUC) GetUsers(adminID uint, filter request.UserFilter) ([]resp.UserListResponse, int, error) {

	// Validasi admin
	currentUser, err := u.userRepo.GetByID(adminID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	if currentUser.Role != "admin" {
		return nil, 0, errors.New("forbidden: only admin can access this")
	}

	// Exclude admin & superadmin
	filter.ExcludeRole = "admin,superadmin"

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
