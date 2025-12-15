// repository/barang_repository.go
package repo

import (
	barang "backend/dto/request/barang"
	response "backend/dto/response/common"
	"backend/model"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BarangRepository interface {
    // CRUD Operations
    Create(barang *model.Barang) error
    FindByID(id uint) (*model.Barang, error)
    FindByKode(kode string) (*model.Barang, error)
    FindAll(filter barang.BarangFilter) ([]model.Barang, *response.Pagination, error)
    Update(barang *model.Barang) error
    Delete(id uint) error
    
    // Stok Management
    KurangiStok(barangID uint, jumlah int) error
    TambahStok(barangID uint, jumlah int) error
    
    // Utility
    IsKodeExists(kode string) (bool, error)
    GetAllKategori() ([]string, error)
}

type barangRepository struct {
    db *gorm.DB
}

func NewBarangRepository(db *gorm.DB) BarangRepository {
    return &barangRepository{db: db}
}

// Create - membuat barang baru
func (r *barangRepository) Create(barang *model.Barang) error {
    return r.db.Create(barang).Error
}

// FindByID - mencari barang berdasarkan ID
func (r *barangRepository) FindByID(id uint) (*model.Barang, error) {
    var barang model.Barang
    err := r.db.First(&barang, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("barang tidak ditemukan")
        }
        return nil, err
    }
    return &barang, nil
}

// FindByKode - mencari barang berdasarkan kode
func (r *barangRepository) FindByKode(kode string) (*model.Barang, error) {
    var barang model.Barang
    err := r.db.Where("kode = ?", kode).First(&barang).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("barang tidak ditemukan")
        }
        return nil, err
    }
    return &barang, nil
}

// FindAll - mengambil semua barang dengan filter dan pagination
func (r *barangRepository) FindAll(filter barang.BarangFilter) ([]model.Barang, *response.Pagination, error) {
	var barangs []model.Barang
	var total int64

	// Default values
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

	// Build base query
	query := r.db.Model(&model.Barang{})

	// Apply keyword search
	if filter.Keyword != "" {
		searchPattern := "%" + filter.Keyword + "%"
		query = query.Where(
			"nama LIKE ? OR merk LIKE ? OR kode LIKE ? OR deskripsi LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Filter kategori
	if filter.Kategori != "" {
		query = query.Where("kategori = ?", filter.Kategori)
	}

	// Filter by multiple status
	if filter.Status != "" {
		statuses := strings.Split(filter.Status, ",")
		query = query.Where("status IN ?", statuses)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Fetch paginated data with sorting
	orderClause := filter.SortBy + " " + filter.SortOrder
	if err := query.Offset(offset).Limit(filter.Limit).Order(orderClause).Find(&barangs).Error; err != nil {
		return nil, nil, err
	}

	// Calculate total pages
	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	// Build pagination response
	pagination := &response.Pagination{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalRows:  int(total),
		TotalPages: totalPages,
	}

	return barangs, pagination, nil
}

// Update - update data barang
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

// Delete - soft delete barang (jika pakai gorm.DeletedAt) atau hard delete
func (r *barangRepository) Delete(id uint) error {
    result := r.db.Delete(&model.Barang{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("barang tidak ditemukan")
    }
    return nil
}

// KurangiStok - mengurangi stok sisa (untuk peminjaman)
func (r *barangRepository) KurangiStok(barangID uint, jumlah int) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        var barang model.Barang
        
        // Lock row untuk prevent race condition
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&barang, barangID).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return errors.New("barang tidak ditemukan")
            }
            return err
        }
        
        // Validasi stok
        if barang.StokSisa < jumlah {
            return errors.New("stok tidak mencukupi")
        }
        
        // Update stok
        newStok := barang.StokSisa - jumlah
        return tx.Model(&barang).Update("stok_sisa", newStok).Error
    })
}

// TambahStok - menambah stok sisa (untuk pengembalian atau penambahan barang)
func (r *barangRepository) TambahStok(barangID uint, jumlah int) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        var barang model.Barang
        
        // Lock row
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&barang, barangID).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return errors.New("barang tidak ditemukan")
            }
            return err
        }
        
        // Validasi stok tidak melebihi total
        newStok := barang.StokSisa + jumlah
        if newStok > barang.StokTotal {
            return errors.New("stok sisa tidak boleh melebihi stok total")
        }
        
        // Update stok
        return tx.Model(&barang).Update("stok_sisa", newStok).Error
    })
}

// IsKodeExists - cek apakah kode barang sudah ada
func (r *barangRepository) IsKodeExists(kode string) (bool, error) {
    var count int64
    err := r.db.Model(&model.Barang{}).Where("kode = ?", kode).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// GetAllKategori - mengambil semua kategori yang ada (unique)
func (r *barangRepository) GetAllKategori() ([]string, error) {
    var kategoris []string
    err := r.db.Model(&model.Barang{}).
        Distinct("kategori").
        Where("kategori != ''").
        Order("kategori ASC").
        Pluck("kategori", &kategoris).Error
    
    if err != nil {
        return nil, err
    }
    return kategoris, nil
}