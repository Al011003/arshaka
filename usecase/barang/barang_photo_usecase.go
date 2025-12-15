// usecase/barang_photo_usecase.go
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

type BarangPhotoUsecase interface {
	UpdatePhoto(barangID uint, file *multipart.FileHeader) (string, error)
	DeletePhoto(barangID uint) error
}

type barangPhotoUsecase struct {
	barangRepo repo.BarangRepository
}

func NewBarangPhotoUsecase(barangRepo repo.BarangRepository) BarangPhotoUsecase {
	return &barangPhotoUsecase{
		barangRepo: barangRepo,
	}
}

func (u *barangPhotoUsecase) UpdatePhoto(barangID uint, file *multipart.FileHeader) (string, error) {
	// Ambil data barang
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil || barang == nil {
		return "", errors.New("barang tidak ditemukan")
	}

	// Folder upload
	uploadPath := "uploads/barang"
	os.MkdirAll(uploadPath, os.ModePerm)

	// Hapus foto lama kalau ada
	if barang.CoverURL != "" {
		oldFile := "." + barang.CoverURL
		_ = os.Remove(oldFile) // abaikan error
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%d_%d%s",
		barangID,
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
	barang.CoverURL = url
	err = u.barangRepo.Update(barang)
	if err != nil {
		// Kalau gagal update DB, hapus file yang baru diupload
		_ = os.Remove(fullPath)
		return "", errors.New("gagal menyimpan foto barang")
	}

	return url, nil
}

func (u *barangPhotoUsecase) DeletePhoto(barangID uint) error {
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil || barang == nil {
		return errors.New("barang tidak ditemukan")
	}

	// Cek apakah ada foto
	if barang.CoverURL == "" {
		return errors.New("barang tidak memiliki foto")
	}

	// Hapus file jika ada
	if barang.CoverURL != "" {
		oldFile := "." + barang.CoverURL
		_ = os.Remove(oldFile) // abaikan error jika file tidak ada
	}

	// Clear foto di DB
	barang.CoverURL = ""
	return u.barangRepo.Update(barang)
}