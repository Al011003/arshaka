package usecase

import (
	request "backend/dto/request/barang"
	responseBarang "backend/dto/response/barang"
	response "backend/dto/response/common"
	"backend/model"
	"backend/repo"
	"errors"
	"fmt"
)

type BarangUseCase interface {
	Create(req request.CreateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error)
	GetByKode(kode string, role string) (interface{}, error)
	GetAll(filter request.BarangFilter, role string) ([]responseBarang.BarangListResponse, *response.Pagination, error)
	UpdateByKode(kode string, req request.UpdateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error)
	DeleteByKode(kode string) error

	GetBarangStatsByKode(kode string) (*responseBarang.BarangStatsResponse, error)
	SetNonAktif(kode string) error // ← NEW
	SetAktif(kode string) error    // ← NEW (optional, untuk reactivate)
	CheckAndSuggestBarangStatus(kode string) (*responseBarang.BarangStatusSuggestion, error)
}

type barangUseCase struct {
	barangRepo repo.BarangRepository
	unitRepo   repo.BarangUnitRepository
}

func NewBarangUseCase(
	barangRepo repo.BarangRepository,
	unitRepo repo.BarangUnitRepository,
) BarangUseCase {
	return &barangUseCase{barangRepo, unitRepo}
}

func (u *barangUseCase) GetBarangStatsByKode(kode string) (*responseBarang.BarangStatsResponse, error) {
	barangID, err := u.barangRepo.GetIDByKode(kode)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	total, err := u.unitRepo.CountTotal(barangID)
	if err != nil {
		return nil, err
	}

	available, err := u.unitRepo.CountAvailable(barangID)
	if err != nil {
		return nil, err
	}

	dipinjam, err := u.unitRepo.CountByStatus(barangID, model.UnitStatusDipinjam)
	if err != nil {
		return nil, err
	}

	maintenance, err := u.unitRepo.CountNeedsMaintenance(barangID)
	if err != nil {
		return nil, err
	}

	return &responseBarang.BarangStatsResponse{
		TotalUnit:   total,
		Siap:        available,
		Dipinjam:    dipinjam,
		Maintenance: maintenance,
	}, nil
}

//
// ================= CRUD =================
//

func (u *barangUseCase) Create(req request.CreateBarangRequest) (*responseBarang.BarangAdminDetailResponse, error) {
	exists, _ := u.barangRepo.IsKodeExists(req.Kode)
	if exists {
		return nil, errors.New("kode barang sudah ada")
	}

	barang := &model.Barang{
		Kode:      req.Kode,
		Nama:      req.Nama,
		Merk:      req.Merk,
		Deskripsi: req.Deskripsi,
		Kategori:  req.Kategori,
		Harga:     req.Harga,
		Status:    model.BarangStatusAktif,
	}

	if err := u.barangRepo.Create(barang); err != nil {
		return nil, err
	}

	return &responseBarang.BarangAdminDetailResponse{
		ID:       barang.ID,
		Kode:     barang.Kode,
		Nama:     barang.Nama,
		Kategori: barang.Kategori,
		Status:   barang.Status,
		Harga:    barang.Harga,
		Stats:    &responseBarang.BarangStatsResponse{},
	}, nil
}

func (u *barangUseCase) GetByKode(kode string, role string) (interface{}, error) {
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil {
		return nil, err
	}

	total, _ := u.unitRepo.CountTotal(barang.ID)
	available, _ := u.unitRepo.CountAvailable(barang.ID)
	
	// ← Cek apakah semua unit nonaktif
	allUnitsInactive := false
	if total > 0 && available == 0 {
		nonaktif, _ := u.unitRepo.CountByStatus(barang.ID, model.UnitStatusNonAktif)
		if nonaktif == total {
			allUnitsInactive = true
		}
	}

	if role != "admin" {
		return &responseBarang.BarangUserDetailResponse{
			ID:            barang.ID,
			Kode:          barang.Kode,          // ← TAMBAH
			Nama:          barang.Nama,
			Merk:          barang.Merk,          // ← TAMBAH
			Deskripsi:     barang.Deskripsi,     // ← TAMBAH
			Kategori:      barang.Kategori,      // ← TAMBAH
			CoverURL:      barang.CoverURL,      // ← TAMBAH
			Harga:         barang.Harga,         // ← TAMBAH (optional)
			TotalUnit:     total,
			AvailableUnit: available,
		}, nil
	}

	stats, err := u.GetBarangStatsByKode(kode)
	if err != nil {
		return nil, err
	}

	resp := MapBarangToDetailResponse(barang, stats)
	resp.AllUnitsInactive = allUnitsInactive // ← tambah field ini
	
	return &resp, nil
}
func (u *barangUseCase) GetAll(
	filter request.BarangFilter,
	role string,
) ([]responseBarang.BarangListResponse, *response.Pagination, error) {

	if role != "admin" {
		filter.Status = model.BarangStatusAktif
	}

	list, pag, err := u.barangRepo.FindAll(filter)
	if err != nil {
		return nil, nil, err
	}

	var out []responseBarang.BarangListResponse
	for _, b := range list {
		total, _ := u.unitRepo.CountTotal(b.ID)
		available, _ := u.unitRepo.CountAvailable(b.ID)

		out = append(out, MapBarangToListResponse(b, total, available))
	}

	return out, pag, nil
}

func (u *barangUseCase) UpdateByKode(
	kode string,
	req request.UpdateBarangRequest,
) (*responseBarang.BarangAdminDetailResponse, error) {

	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil {
		return nil, err
	}

	if req.Nama != nil {
		barang.Nama = *req.Nama
	}
	if req.Merk != nil {
		barang.Merk = *req.Merk
	}
	if req.Deskripsi != nil {
		barang.Deskripsi = *req.Deskripsi
	}
	if req.Kategori != nil {
		barang.Kategori = *req.Kategori
	}
	if req.Status != nil {
		barang.Status = *req.Status
	}
	if req.Harga != nil {
		barang.Harga = *req.Harga
	}

	if err := u.barangRepo.Update(barang); err != nil {
		return nil, err
	}

	stats, _ := u.GetBarangStatsByKode(kode)
	resp := MapBarangToDetailResponse(barang, stats)
	return &resp, nil
}

func (u *barangUseCase) DeleteByKode(kode string) error {
	barangID, err := u.barangRepo.GetIDByKode(kode)
	if err != nil {
		return errors.New("barang tidak ditemukan")
	}

	dipinjam, _ := u.unitRepo.CountByStatus(barangID, model.UnitStatusDipinjam)
	if dipinjam > 0 {
		return errors.New("barang masih dipinjam")
	}

	return u.barangRepo.DeleteByKode(kode)
}

func (u *barangUseCase) SetNonAktif(kode string) error {
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil {
		return errors.New("barang tidak ditemukan")
	}

	// Cek apakah ada unit yang sedang dipinjam
	dipinjam, _ := u.unitRepo.CountByStatus(barang.ID, model.UnitStatusDipinjam)
	if dipinjam > 0 {
		return errors.New("tidak bisa menonaktifkan barang: masih ada unit yang sedang dipinjam")
	}

	// Set barang jadi nonaktif
	barang.Status = model.BarangStatusNonAktif
	if err := u.barangRepo.Update(barang); err != nil {
		return err
	}

	// Cascade: set semua unit jadi nonaktif
	if err := u.unitRepo.SetAllStatusByBarangID(barang.ID, model.UnitStatusNonAktif); err != nil {
		return err
	}

	return nil
}

func (u *barangUseCase) SetAktif(kode string) error {
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil {
		return errors.New("barang tidak ditemukan")
	}

	// Set barang jadi aktif
	barang.Status = model.BarangStatusAktif
	if err := u.barangRepo.Update(barang); err != nil {
		return err
	}

	// Cascade: set semua unit kembali ke status normal
	// Units akan auto-compute statusnya berdasarkan kondisi (via BeforeSave hook)
	// Kita perlu trigger update untuk setiap unit
	if err := u.unitRepo.ReactivateAllByBarangID(barang.ID); err != nil {
		return err
	}

	return nil
}

// usecase/barang/barang.go

// CheckAndSuggestBarangStatus - cek apakah barang perlu dinonaktifkan
func (u *barangUseCase) CheckAndSuggestBarangStatus(kode string) (*responseBarang.BarangStatusSuggestion, error) {
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	total, _ := u.unitRepo.CountTotal(barang.ID)
	
	// Case 1: Tidak ada unit sama sekali
	if total == 0 {
		shouldChange := barang.Status != model.BarangStatusNonAktif // ← FIX
		return &responseBarang.BarangStatusSuggestion{
			CurrentStatus:   barang.Status,
			SuggestedStatus: model.BarangStatusNonAktif,
			Reason:          "Tidak ada unit sama sekali",
			ShouldChange:    shouldChange,
		}, nil
	}

	nonaktif, _ := u.unitRepo.CountByStatus(barang.ID, model.UnitStatusNonAktif)
	
	// Case 2: Semua unit nonaktif
	if nonaktif == total {
		shouldChange := barang.Status != model.BarangStatusNonAktif // ← FIX
		return &responseBarang.BarangStatusSuggestion{
			CurrentStatus:   barang.Status,
			SuggestedStatus: model.BarangStatusNonAktif,
			Reason:          fmt.Sprintf("Semua unit (%d) sudah nonaktif", total),
			ShouldChange:    shouldChange,
		}, nil
	}

	// Case 3: Ada unit yang masih aktif
	shouldChange := barang.Status != model.BarangStatusAktif // ← FIX
	return &responseBarang.BarangStatusSuggestion{
		CurrentStatus:   barang.Status,
		SuggestedStatus: model.BarangStatusAktif,
		Reason:          fmt.Sprintf("%d dari %d unit masih aktif", total-nonaktif, total),
		ShouldChange:    shouldChange,
	}, nil
}
//
// ================= MAPPER =================
//

func MapBarangToListResponse(
	b model.Barang,
	total, available int64,
) responseBarang.BarangListResponse {
	return responseBarang.BarangListResponse{
		ID:            b.ID,
		Kode:          b.Kode,
		Nama:          b.Nama,
		Merk:          b.Merk,
		Kategori:      b.Kategori,
		Status:        b.Status,
		CoverURL:      b.CoverURL,
		Harga:         b.Harga,
		TotalUnit:     total,
		AvailableUnit: available,
	}
}

func MapBarangToDetailResponse(
	b *model.Barang,
	stats *responseBarang.BarangStatsResponse,
) responseBarang.BarangAdminDetailResponse {
	return responseBarang.BarangAdminDetailResponse{
		ID:        b.ID,
		Kode:      b.Kode,
		Nama:      b.Nama,
		Merk:      b.Merk,
		Deskripsi: b.Deskripsi,
		Kategori:  b.Kategori,
		Status:    b.Status,
		CoverURL:  b.CoverURL,
		Harga:     b.Harga,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Stats:     stats,
	}
}