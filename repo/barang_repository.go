// repository/barang_repository.go
package repo

import (
	barang "backend/dto/request/barang"
	response "backend/dto/response/common"
	"backend/model"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type BarangRepository interface {
	Create(barang *model.Barang) error
	FindByID(id uint) (*model.Barang, error)
	FindByKode(kode string) (*model.Barang, error)
	FindAll(filter barang.BarangFilter) ([]model.Barang, *response.Pagination, error)
	Update(barang *model.Barang) error
	Delete(id uint) error

	IsKodeExists(kode string) (bool, error)
	IsKodeExistsExcludingID(kode string, excludeID uint) (bool, error) // untuk update
	GetAllKategori() ([]string, error)
	IsActive(barangID uint) (bool, error)
	
	// Stats helper (opsional tapi recommended)
	CountTotalUnits(barangID uint) (int64, error)
	CountUnitsByStatus(barangID uint, status string) (int64, error)

	DeleteByKode(kode string) error
	IsActiveByKode(kode string) (bool, error)

	GetIDByKode(kode string) (uint, error)
	FindSimilarBarang(kategori, merk string, excludeBarangID uint, limit int, nama string) ([]model.Barang, error)
}

type barangRepository struct {
	db *gorm.DB
}

func NewBarangRepository(db *gorm.DB) BarangRepository {
	return &barangRepository{db: db}
}

func (r *barangRepository) Create(barang *model.Barang) error {
	return r.db.Create(barang).Error
}

func (r *barangRepository) FindByID(id uint) (*model.Barang, error) {
	var barang model.Barang
	if err := r.db.First(&barang, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("barang tidak ditemukan")
		}
		return nil, err
	}
	return &barang, nil
}

// FindByIDWithUnits - untuk load barang beserta units & komponennya
func (r *barangRepository) FindByIDWithUnits(id uint) (*model.Barang, error) {
	var barang model.Barang
	if err := r.db.Preload("Units.Komponen").First(&barang, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("barang tidak ditemukan")
		}
		return nil, err
	}
	return &barang, nil
}

func (r *barangRepository) FindByKode(kode string) (*model.Barang, error) {
	var barang model.Barang
	if err := r.db.Where("kode = ?", kode).First(&barang).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("barang tidak ditemukan")
		}
		return nil, err
	}
	return &barang, nil
}

func (r *barangRepository) FindAll(filter barang.BarangFilter) ([]model.Barang, *response.Pagination, error) {
	var barangs []model.Barang
	var total int64

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	offset := (filter.Page - 1) * filter.Limit
	query := r.db.Model(&model.Barang{})

	if filter.Keyword != "" {
		pattern := "%" + filter.Keyword + "%"
		query = query.Where(
			"nama LIKE ? OR merk LIKE ? OR kode LIKE ? OR deskripsi LIKE ?",
			pattern, pattern, pattern, pattern,
		)
	}

	if filter.Kategori != "" {
		query = query.Where("kategori = ?", filter.Kategori)
	}

	if filter.Status != "" {
		statuses := strings.Split(filter.Status, ",")
		query = query.Where("status IN ?", statuses)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	order := filter.SortBy + " " + filter.SortOrder
	if err := query.Offset(offset).Limit(filter.Limit).Order(order).Find(&barangs).Error; err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	pagination := &response.Pagination{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalRows:  int(total),
		TotalPages: totalPages,
	}

	return barangs, pagination, nil
}

func (r *barangRepository) Update(barang *model.Barang) error {
	result := r.db.Save(barang)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("barang tidak ditemukan")
	}
	return nil
}

func (r *barangRepository) Delete(id uint) error {
	// Soft delete dengan GORM
	return r.db.Delete(&model.Barang{}, id).Error
}

func (r *barangRepository) IsKodeExists(kode string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Barang{}).Where("kode = ?", kode).Count(&count).Error
	return count > 0, err
}

// IsKodeExistsExcludingID - untuk validasi saat update (exclude barang yg lg diupdate)
func (r *barangRepository) IsKodeExistsExcludingID(kode string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Barang{}).
		Where("kode = ? AND id != ?", kode, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *barangRepository) GetAllKategori() ([]string, error) {
	var kategori []string
	err := r.db.Model(&model.Barang{}).
		Distinct("kategori").
		Where("kategori != ''").
		Order("kategori ASC").
		Pluck("kategori", &kategori).Error
	return kategori, err
}

func (r *barangRepository) IsActive(barangID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Barang{}).
		Where("id = ? AND status = ?", barangID, "aktif").
		Count(&count).Error
	return count > 0, err
}

// CountTotalUnits - hitung total unit dari barang ini
func (r *barangRepository) CountTotalUnits(barangID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ?", barangID).
		Count(&count).Error
	return count, err
}

// CountUnitsByStatus - hitung unit berdasarkan status (siap, dipinjam, dll)
func (r *barangRepository) CountUnitsByStatus(barangID uint, status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ? AND status = ?", barangID, status).
		Count(&count).Error
	return count, err
}

func (r *barangRepository) DeleteByKode(kode string) error {
	return r.db.
		Where("kode = ?", kode).
		Delete(&model.Barang{}).
		Error
}

func (r *barangRepository) IsActiveByKode(kode string) (bool, error) {
	var count int64
	err := r.db.
		Model(&model.Barang{}).
		Where("kode = ? AND status = ?", kode, "aktif").
		Count(&count).
		Error

	return count > 0, err
}

func (r *barangRepository) GetIDByKode(kode string) (uint, error) {
	var id uint
	err := r.db.
		Model(&model.Barang{}).
		Select("id").
		Where("kode = ?", kode).
		Take(&id).
		Error

	return id, err
}

func (r *barangRepository) FindSimilarBarang(
	kategori string,
	merk string,
	excludeBarangID uint,
	limit int,
	nama string, // 🔥 TAMBAH parameter nama
) ([]model.Barang, error) {
	var barangs []model.Barang

	// 1️⃣ Prioritas 1: Nama mirip (LIKE) + Kategori sama + Merk sama
	if kategori != "" && merk != "" && nama != "" {
		err := r.db.
			Where("status = ?", "aktif").
			Where("id != ?", excludeBarangID).
			Where("nama LIKE ? AND kategori = ? AND merk = ?", "%"+nama+"%", kategori, merk).
			Order("created_at DESC").
			Limit(limit).
			Find(&barangs).Error

		if err != nil {
			return nil, err
		}

		if len(barangs) >= limit {
			return barangs, nil
		}
	}

	// 2️⃣ Prioritas 2: Nama mirip + Kategori sama + Merk beda
	if kategori != "" && nama != "" && len(barangs) < limit {
		var additionalBarangs []model.Barang
		remaining := limit - len(barangs)

		existingIDs := make([]uint, len(barangs))
		for i, b := range barangs {
			existingIDs[i] = b.ID
		}

		query := r.db.
			Where("status = ?", "aktif").
			Where("id != ?", excludeBarangID).
			Where("nama LIKE ? AND kategori = ?", "%"+nama+"%", kategori)

		if merk != "" {
			query = query.Where("merk != ?", merk)
		}

		if len(existingIDs) > 0 {
			query = query.Where("id NOT IN ?", existingIDs)
		}

		err := query.
			Order("created_at DESC").
			Limit(remaining).
			Find(&additionalBarangs).Error

		if err != nil {
			return nil, err
		}

		barangs = append(barangs, additionalBarangs...)
	}

	return barangs, nil
}