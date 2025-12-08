package usecase

import (
	req "backend/dto/request/auth"
	res "backend/dto/response/auth"
	"backend/repo"
	"backend/token"
	"backend/utils"
	"errors"
	"fmt"
	"os"
)

type LoginUsecase interface {
	Login(req req.LoginRequest) (*res.LoginResponse, error)
}

type loginUsecase struct {
	repo repo.AuthRepository
}

func NewLoginUsecase(r repo.AuthRepository) LoginUsecase {
	return &loginUsecase{
		repo: r,
	}
}

func (u *loginUsecase) Login(req req.LoginRequest) (*res.LoginResponse, error) {
	// 1. cek user ada atau tidak
	user, err := u.repo.FindUserByNra(req.NRA)
	if err != nil || user == nil {
		return nil, errors.New("NRA atau password salah")
	}

	// 2. cek password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("NRA atau password salah")
	}

	// 3. generate access token
	fmt.Println("ROLE DARI DATABASE:", user.Role)
	fmt.Println("SECRET_KEY VALIDATE =", os.Getenv("SECRET_KEY"))
	accessToken, err := token.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal generate access token")
	}

	// 4. generate refresh token
	refreshToken, err := token.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal generate refresh token")
	}

	// 5. simpan refresh token di DB
	user.RefreshToken = refreshToken
	if err := u.repo.UpdateUser(user); err != nil {
		return nil, errors.New("gagal update refresh token")
	}

	// 6. mapping ke response DTO
	resp := &res.LoginResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		MustChangePassword: user.IsDefaultPassword,
		MustFillEmail: user.IsEmailNull,
		Role: user.Role,
	}

	return resp, nil
}
