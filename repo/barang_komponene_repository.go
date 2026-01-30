// repository/barang_komponen_repository.go
package repo

import (
	"backend/model"
	"errors"

	"gorm.io/gorm"
)

type BarangKomponenRepository interface {
	// CREATE
	Create(comp *model.BarangKomponen) error
	BulkCreate(comps []model.BarangKomponen) error

	// READ
	FindByUnit(unitID uint) ([]model.BarangKomponen, error)
	FindByID(id uint) (*model.BarangKomponen, error)

	// UPDATE
	Update(comp *model.BarangKomponen) error
	UpdateStatus(id uint, status string) error
	UpdateJumlahAktual(id uint, jumlahAktual int) error

	// DELETE
	Delete(id uint) error
	DeleteByUnit(unitID uint) error // hapus semua komponen dari unit

	// VALIDATION & CHECKING
	IsUnitComplete(unitID uint) (bool, error)
	HasMissingComponent(unitID uint) (bool, error)
	HasDamagedComponent(unitID uint) (bool, error)
	
	// STATS
	CountByStatus(unitID uint, status string) (int64, error)
	CountIncomplete(unitID uint) (int64, error) // jumlah_aktual < jumlah_wajib
	GetCompletionRate(unitID uint) (float64, error) // persentase kelengkapan
	ExistsByUnitAndNama(unitID uint, nama string) (bool, error)
}

type barangKomponenRepository struct {
	db *gorm.DB
}

func NewBarangKomponenRepository(db *gorm.DB) BarangKomponenRepository {
	return &barangKomponenRepository{db: db}
}

// ================= CREATE =================

func (r *barangKomponenRepository) Create(comp *model.BarangKomponen) error {
	return r.db.Create(comp).Error
}

func (r *barangKomponenRepository) BulkCreate(comps []model.BarangKomponen) error {
	if len(comps) == 0 {
		return nil
	}
	return r.db.Create(&comps).Error
}

// ================= READ =================

func (r *barangKomponenRepository) FindByUnit(unitID uint) ([]model.BarangKomponen, error) {
	var comps []model.BarangKomponen
	err := r.db.Where("barang_unit_id = ?", unitID).
		Order("nama ASC").
		Find(&comps).Error
	return comps, err
}

func (r *barangKomponenRepository) FindByID(id uint) (*model.BarangKomponen, error) {
	var comp model.BarangKomponen
	err := r.db.First(&comp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("komponen tidak ditemukan")
		}
		return nil, err
	}
	return &comp, nil
}

// ================= UPDATE =================

func (r *barangKomponenRepository) Update(comp *model.BarangKomponen) error {
	res := r.db.Save(comp)
	if res.RowsAffected == 0 {
		return errors.New("komponen tidak ditemukan")
	}
	return res.Error
}

func (r *barangKomponenRepository) UpdateStatus(id uint, status string) error {
	res := r.db.Model(&model.BarangKomponen{}).
		Where("id = ?", id).
		Update("status", status)

	if res.RowsAffected == 0 {
		return errors.New("komponen tidak ditemukan")
	}
	return res.Error
}

func (r *barangKomponenRepository) UpdateJumlahAktual(id uint, jumlahAktual int) error {
	res := r.db.Model(&model.BarangKomponen{}).
		Where("id = ?", id).
		Update("jumlah_aktual", jumlahAktual)

	if res.RowsAffected == 0 {
		return errors.New("komponen tidak ditemukan")
	}
	return res.Error
}

// ================= DELETE =================

func (r *barangKomponenRepository) Delete(id uint) error {
	return r.db.Delete(&model.BarangKomponen{}, id).Error
}

func (r *barangKomponenRepository) DeleteByUnit(unitID uint) error {
	return r.db.Where("barang_unit_id = ?", unitID).
		Delete(&model.BarangKomponen{}).Error
}

// ================= VALIDATION & CHECKING =================

// IsUnitComplete - true jika semua komponen statusnya "lengkap"
func (r *barangKomponenRepository) IsUnitComplete(unitID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BarangKomponen{}).
		Where("barang_unit_id = ? AND status != ?", unitID, "lengkap").
		Count(&count).Error
	return count == 0, err
}

// HasMissingComponent - true jika ada komponen kurang/rusak/hilang
func (r *barangKomponenRepository) HasMissingComponent(unitID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BarangKomponen{}).
		Where("barang_unit_id = ? AND status IN ?", unitID, []string{"kurang", "hilang"}).
		Count(&count).Error
	return count > 0, err
}

// HasDamagedComponent - true jika ada komponen rusak
func (r *barangKomponenRepository) HasDamagedComponent(unitID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BarangKomponen{}).
		Where("barang_unit_id = ? AND status = ?", unitID, "rusak").
		Count(&count).Error
	return count > 0, err
}

// ================= STATS =================

func (r *barangKomponenRepository) CountByStatus(unitID uint, status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangKomponen{}).
		Where("barang_unit_id = ? AND status = ?", unitID, status).
		Count(&count).Error
	return count, err
}

// CountIncomplete - hitung komponen yang jumlah_aktual < jumlah_wajib
func (r *barangKomponenRepository) CountIncomplete(unitID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangKomponen{}).
		Where("barang_unit_id = ? AND jumlah_aktual < jumlah_wajib", unitID).
		Count(&count).Error
	return count, err
}

// GetCompletionRate - return persentase kelengkapan komponen (0.0 - 100.0)
func (r *barangKomponenRepository) GetCompletionRate(unitID uint) (float64, error) {
	type Result struct {
		TotalWajib  int64
		TotalAktual int64
	}

	var result Result
	err := r.db.Model(&model.BarangKomponen{}).
		Select("SUM(jumlah_wajib) as total_wajib, SUM(jumlah_aktual) as total_aktual").
		Where("barang_unit_id = ?", unitID).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	if result.TotalWajib == 0 {
		return 100.0, nil // kalau gak ada komponen = dianggap lengkap
	}

	rate := (float64(result.TotalAktual) / float64(result.TotalWajib)) * 100.0
	return rate, nil
}

func (r *barangKomponenRepository) ExistsByUnitAndNama(unitID uint, nama string) (bool, error) {
	var count int64
	err := r.db.Model(&model.BarangKomponen{}).
		Where("barang_unit_id = ? AND nama = ?", unitID, nama).
		Count(&count).Error
	return count > 0, err
}
