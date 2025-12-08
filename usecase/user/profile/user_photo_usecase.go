package usecase

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"backend/repo"
)

type UserPhotoUsecase interface {
	UpdatePhoto(userID uint, file *multipart.FileHeader) (string, error)
	DeletePhoto(userID uint) error
}

type userPhotoUsecase struct {
	userRepo repo.UserRepo
}

func NewUserPhotoUsecase(userRepo repo.UserRepo) UserPhotoUsecase {
	return &userPhotoUsecase{userRepo}
}

func (u *userPhotoUsecase) UpdatePhoto(userID uint, file *multipart.FileHeader) (string, error) {

	// Ambil data user
	user, err := u.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return "", errors.New("user tidak ditemukan")
	}

	// Folder upload
	uploadPath := "uploads/profile"
	os.MkdirAll(uploadPath, os.ModePerm)

	// Hapus foto lama kalau ada
	if user.FotoURL != "" {
		oldFile := "." + user.FotoURL
		_ = os.Remove(oldFile) // abaikan error
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%d_%d%s",
		userID,
		time.Now().UnixNano(),
		filepath.Ext(file.Filename),
	)
	fullPath := filepath.Join(uploadPath, filename)

	// Open file (multipart)
	src, err := file.Open()
	if err != nil {
		return "", errors.New("gagal membuka file upload")
	}
	defer src.Close()

	// Buat file tujuan
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", errors.New("gagal menyimpan file")
	}
	defer dst.Close()

	// Copy isi file
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", errors.New("gagal menyalin file")
	}

	// URL yang disimpan ke DB
	url := "/" + fullPath

	// Update DB
	err = u.userRepo.UpdateProfilePhoto(userID, url)
	if err != nil {
		return "", errors.New("gagal menyimpan foto profile")
	}

	return url, nil
}

func (u *userPhotoUsecase) DeletePhoto(userID uint) error {

	user, err := u.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}

	// Hapus file jika ada
	if user.FotoURL != "" {
		_ = os.Remove("." + user.FotoURL)
	}

	// Clear foto di DB
	return u.userRepo.UpdateProfilePhoto(userID, "")
}
