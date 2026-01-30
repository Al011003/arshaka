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
	UpdatePhotoByKode(kode string, file *multipart.FileHeader) (string, error)
	DeletePhotoByKode(kode string) error
}

type barangPhotoUsecase struct {
	barangRepo repo.BarangRepository
}

func NewBarangPhotoUsecase(barangRepo repo.BarangRepository) BarangPhotoUsecase {
	return &barangPhotoUsecase{
		barangRepo: barangRepo,
	}
}

// =====================
// UPDATE PHOTO (BY KODE)
// =====================
func (u *barangPhotoUsecase) UpdatePhotoByKode(
	kode string,
	file *multipart.FileHeader,
) (string, error) {

	// Ambil barang via KODE
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil || barang == nil {
		return "", errors.New("barang tidak ditemukan")
	}

	// Folder upload
	uploadPath := "uploads/barang"
	_ = os.MkdirAll(uploadPath, os.ModePerm)

	// Hapus foto lama jika ada
	if barang.CoverURL != "" {
		oldFile := "." + barang.CoverURL
		_ = os.Remove(oldFile)
	}

	// Generate filename unik
	filename := fmt.Sprintf(
		"%s_%d%s",
		kode,
		time.Now().UnixNano(),
		filepath.Ext(file.Filename),
	)

	fullPath := filepath.Join(uploadPath, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", errors.New("gagal membuka file upload")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", errors.New("gagal menyimpan file")
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New("gagal menyalin file")
	}

	// URL untuk DB
	url := "/uploads/barang/" + filename

	// Update DB
	barang.CoverURL = url
	if err := u.barangRepo.Update(barang); err != nil {
		_ = os.Remove(fullPath)
		return "", errors.New("gagal menyimpan foto barang")
	}

	return url, nil
}

// =====================
// DELETE PHOTO (BY KODE)
// =====================
func (u *barangPhotoUsecase) DeletePhotoByKode(kode string) error {
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil || barang == nil {
		return errors.New("barang tidak ditemukan")
	}

	if barang.CoverURL == "" {
		return errors.New("barang tidak memiliki foto")
	}

	// Hapus file
	oldFile := "." + barang.CoverURL
	_ = os.Remove(oldFile)

	// Clear DB
	barang.CoverURL = ""
	return u.barangRepo.Update(barang)
}
