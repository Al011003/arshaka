// usecase/barang_usecase.go
package usecase

import (
	requestBarang "backend/dto/request/barang"
	responseBarang "backend/dto/response/barang"
	response "backend/dto/response/common"
	"backend/model"
	"backend/repo"
	"errors"
)

type BarangUseCase interface {
	Create(req requestBarang.CreateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error)
	GetByID(id uint, role string) (interface{}, error) // tetap interface karena beda response by role
	GetAll(filter requestBarang.BarangFilter, role string) ([]responseBarang.BarangListResponse, *response.Pagination, error)
	Update(id uint, req requestBarang.UpdateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error)
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

//
// 🔥 Helper: Convert Model → DTO Admin
//
func (u *barangUseCase) toAdminDetail(m *model.Barang) *responseBarang.BarangAdminDetailResponse {
	return &responseBarang.BarangAdminDetailResponse{
		ID:        m.ID,
		Kode:      m.Kode,
		Nama:      m.Nama,
		Merk:      m.Merk,
		Deskripsi: m.Deskripsi,
		Kategori:  m.Kategori,
		StokTotal: m.StokTotal,
		StokSisa:  m.StokSisa,
		Status:    m.Status, // ← TAMBAH INI
		TahunBeli: m.TahunBeli,
		HargaBeli: m.HargaBeli,
		CoverURL:  m.CoverURL,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

//
// 🔥 Helper: Convert Model → DTO User
//
func (u *barangUseCase) toUserDetail(m *model.Barang) *responseBarang.BarangUserDetailResponse {
	return &responseBarang.BarangUserDetailResponse{
		ID:        m.ID,
		Kode:      m.Kode,
		Nama:      m.Nama,
		Merk:      m.Merk,
		Deskripsi: m.Deskripsi,
		Kategori:  m.Kategori,
		StokTotal: m.StokTotal,
		StokSisa:  m.StokSisa,
		Status:    m.Status, // ← TAMBAH INI
		CoverURL:  m.CoverURL,
	}
}

//
// 🔥 Helper: List view
//
func (u *barangUseCase) toBarangListResponse(m *model.Barang) responseBarang.BarangListResponse {
	return responseBarang.BarangListResponse{
		ID:        m.ID,
		Kode:      m.Kode,
		Nama:      m.Nama,
		Merk:      m.Merk,
		Kategori:  m.Kategori,
		StokTotal: m.StokTotal,
		StokSisa:  m.StokSisa,
		Status:    m.Status, // ← TAMBAH INI
		CoverURL:  m.CoverURL,
	}
}

//
// 🔥 Create Barang
//
func (u *barangUseCase) Create(req requestBarang.CreateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error) {
	// Validasi kode unik
	exists, err := u.barangRepo.IsKodeExists(req.Kode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("kode barang sudah ada")
	}

	// Validasi stok (early validation)
	if req.StokSisa > req.StokTotal {
		return nil, errors.New("stok sisa tidak boleh lebih dari stok total")
	}

	// Map request ke model
	newBarang := &model.Barang{
		Kode:      req.Kode,
		Nama:      req.Nama,
		Merk:      req.Merk,
		Deskripsi: req.Deskripsi,
		Kategori:  req.Kategori,
		StokTotal: req.StokTotal,
		StokSisa:  req.StokSisa,
		TahunBeli: req.TahunBeli,
		HargaBeli: req.HargaBeli,
		CoverURL:  "",
		// Status auto set di BeforeCreate hook (tersedia/habis)
	}

	// Save ke database (BeforeCreate hook akan jalan otomatis)
	if err := u.barangRepo.Create(newBarang); err != nil {
		return nil, err
	}

	// Convert ke response DTO
	return u.toAdminDetail(newBarang), nil
}

//
// 🔥 Get By ID — beda sesuai ROLE
//
func (u *barangUseCase) GetByID(id uint, role string) (interface{}, error) {
	barang, err := u.barangRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// User biasa cuma boleh lihat barang Tersedia/Habis
	if role != "admin" && role != "superadmin" {
		if barang.Status == "nonaktif" {
			return nil, errors.New("barang tidak tersedia")
		}
		return u.toUserDetail(barang), nil
	}

	// Admin / superadmin dapat full detail
	return u.toAdminDetail(barang), nil
}


//
// 🔥 Get All List
//
func (u *barangUseCase) GetAll(filter requestBarang.BarangFilter, role string) ([]responseBarang.BarangListResponse, *response.Pagination, error) {

    if role != "admin" && role != "superadmin" {

        // user biasa --> kalau user memasukkan status manual tetap override
        if filter.Status == "" {
            filter.Status = "tersedia,habis"
        }

    } else {
        // admin / superadmin --> JANGAN hapus status kalau user pakai filter!
        // cuma kalau user tidak set status, biarkan kosong (ambil semua)
        // jadi ga usah apa2
    }

    barangs, pagination, err := u.barangRepo.FindAll(filter)
    if err != nil {
        return nil, nil, err
    }

    var out []responseBarang.BarangListResponse
    for _, b := range barangs {
        out = append(out, u.toBarangListResponse(&b))
    }

    return out, pagination, nil
}



//
// 🔥 Update Barang
//
func (u *barangUseCase) Update(id uint, req requestBarang.UpdateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error) {
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
	if req.Status != nil {
		existingBarang.Status = *req.Status
	}

	// Validasi stok setelah semua field di-update
	if existingBarang.StokSisa > existingBarang.StokTotal {
		return nil, errors.New("stok sisa tidak boleh melebihi stok total")
	}

	// Save ke database (BeforeUpdate hook akan jalan otomatis untuk set status)
	if err := u.barangRepo.Update(existingBarang); err != nil {
		return nil, err
	}

	// Convert ke response DTO
	return u.toAdminDetail(existingBarang), nil
}

//
// 🔥 Delete
//
func (u *barangUseCase) Delete(id uint) error {
	// Cek barang exist
	_, err := u.barangRepo.FindByID(id)
	if err != nil {
		return err
	}

	return u.barangRepo.Delete(id)
}