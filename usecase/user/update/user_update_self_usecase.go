package usecase

import (
	req "backend/dto/request/user"
	"backend/repo"
	"errors"
)

type UserSelfUsecase interface {
	UpdateSelf(userID uint, r req.UserUpdateRequest) error
}

type userSelfUsecase struct {
	userRepo repo.UserRepo
}

func NewUserSelfUsecase(userRepo repo.UserRepo) UserSelfUsecase {
	return &userSelfUsecase{
		userRepo: userRepo,
	}
}

func (u *userSelfUsecase) UpdateSelf(userID uint, r req.UserUpdateRequest) error {

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("gagal mengambil data user")
	}
	if user == nil {
		return errors.New("user tidak ditemukan")
	}

	user.Nama = r.Nama

	if err := u.userRepo.Update(user); err != nil {
		return errors.New("gagal menyimpan perubahan user")
	}

	return nil
}
