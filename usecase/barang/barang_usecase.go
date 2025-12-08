// usecase/barang_usecase.go
package usecase

import (
	barang "backend/dto/request/barang"
	response "backend/dto/response/common"
	"backend/model"
	"backend/repo"
	"errors"
)

type BarangUseCase interface {
    Create(req barang.CreateBarangRequest) (*model.Barang, error)
    GetByID(id uint) (*model.Barang, error)
    GetAll(filter barang.BarangFilter) ([]model.Barang, *response.Pagination, error)
    Update(id uint, req barang.UpdateBarangRequest) (*model.Barang, error)
    Delete(id uint) error
}

type barangUseCase struct {
    barangRepo repo.BarangRepository
}

func NewBarangUseCase(barangRepo repo.BarangRepository) BarangUseCase {
    return &barangUseCase{
        barangRepo: barangRepo,
    }
}

func (u *barangUseCase) Create(req barang.CreateBarangRequest) (*model.Barang, error) {
    // Validasi kode unik
    exists, err := u.barangRepo.IsKodeExists(req.Kode)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("kode barang sudah ada")
    }
    
    // Validasi stok sisa tidak boleh lebih dari stok total
    if req.StokSisa > req.StokTotal {
        return nil, errors.New("stok sisa tidak boleh lebih dari stok total")
    }
    
    // Map request ke model
    newBarang := &model.Barang{
        Kode:        req.Kode,
        Nama:        req.Nama,
        Merk:        req.Merk,
        Deskripsi:   req.Deskripsi,
        Kategori:    req.Kategori,
        StokTotal:   req.StokTotal,
        StokSisa:    req.StokSisa,
        TahunBeli:   req.TahunBeli,
        HargaBeli:   req.HargaBeli,
        CoverURL:    "", // Kosong dulu, diisi kalau upload foto
    }
    
    // Save ke database
    if err := u.barangRepo.Create(newBarang); err != nil {
        return nil, err
    }
    
    return newBarang, nil
}

func (u *barangUseCase) GetByID(id uint) (*model.Barang, error) {
    barang, err := u.barangRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    return barang, nil
}

func (u *barangUseCase) GetAll(filter barang.BarangFilter) ([]model.Barang, *response.Pagination, error) {
    return u.barangRepo.FindAll(filter)
}

func (u *barangUseCase) Update(id uint, req barang.UpdateBarangRequest) (*model.Barang, error) {
    // Cek barang exist
    existingBarang, err := u.barangRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Update HANYA field yang dikirim (tidak nil)
    if req.Kode != nil {
        // Validasi kode unik (kecuali kode yang sama)
        if *req.Kode != existingBarang.Kode {
            exists, err := u.barangRepo.IsKodeExists(*req.Kode)
            if err != nil {
                return nil, err
            }
            if exists {
                return nil, errors.New("kode barang sudah ada")
            }
        }
        existingBarang.Kode = *req.Kode
    }
    
    if req.Nama != nil {
        existingBarang.Nama = *req.Nama
    }
    
    if req.Merk != nil {
        existingBarang.Merk = *req.Merk
    }
    
    if req.Deskripsi != nil {
        existingBarang.Deskripsi = *req.Deskripsi
    }
    
    if req.Kategori != nil {
        existingBarang.Kategori = *req.Kategori
    }
    
    if req.StokTotal != nil {
        existingBarang.StokTotal = *req.StokTotal
    }
    
    if req.StokSisa != nil {
        existingBarang.StokSisa = *req.StokSisa
    }
    
    if req.TahunBeli != nil {
        existingBarang.TahunBeli = *req.TahunBeli
    }
    
    if req.HargaBeli != nil {
        existingBarang.HargaBeli = *req.HargaBeli
    }
    
    // Validasi stok setelah semua field di-update
    if existingBarang.StokSisa > existingBarang.StokTotal {
        return nil, errors.New("stok sisa tidak boleh lebih dari stok total")
    }
    
    // Save ke database
    if err := u.barangRepo.Update(existingBarang); err != nil {
        return nil, err
    }
    
    return existingBarang, nil
}


func (u *barangUseCase) Delete(id uint) error {
    // Cek barang exist
    _, err := u.barangRepo.FindByID(id)
    if err != nil {
        return err
    }
    
    return u.barangRepo.Delete(id)
}