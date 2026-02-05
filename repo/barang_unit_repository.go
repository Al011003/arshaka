// repository/barang_unit_repository.go
package repo

import (
	"backend/model"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type BarangUnitRepository interface {
	// === CREATE ===
	Create(unit *model.BarangUnit) error
	BulkCreate(units []model.BarangUnit) error

	// === READ ===
	FindByID(id uint) (*model.BarangUnit, error)
	FindByKodeUnit(kode string) (*model.BarangUnit, error)
	FindAllByBarangID(barangID uint) ([]model.BarangUnit, error)

	// Ambil unit yang bisa dipinjam
	FindAvailableUnits(barangID uint, limit int) ([]model.BarangUnit, error)

	// Ambil unit beserta komponennya (buat detail pinjaman)
	FindWithComponents(unitID uint) (*model.BarangUnit, error)

	// === UPDATE ===
	Update(unit *model.BarangUnit) error
	UpdateStatus(unitID uint, status string) error
	UpdateKondisi(unitID uint, kondisi string) error

	// === DELETE ===
	Delete(id uint) error

	// === VALIDATION ===
	IsKodeUnitExists(kodeUnit string) (bool, error)
	IsKodeUnitExistsExcludingID(kodeUnit string, excludeID uint) (bool, error)

	// === COUNT (pengganti stok lama) ===
	CountTotal(barangID uint) (int64, error)
	CountAvailable(barangID uint) (int64, error)
	CountByStatus(barangID uint, status string) (int64, error)
	
	// === STATS ===
	CountByKondisi(barangID uint, kondisi string) (int64, error)
	CountNeedsMaintenance(barangID uint) (int64, error) // status = maintenance_required

	FindByKodeWithComponents(kode string) (*model.BarangUnit, error)

	UpdateStatusByKode(kode string, status string) error
	UpdateKondisiByKode(kode string, kondisi string) error

	DeleteByKode(kode string) error


	// 🔥 TAMBAHAN BARU
	GetLastUnitNumber(barangID uint) (int, error)
		SetAllStatusByBarangID(barangID uint, status string) error
	ReactivateAllByBarangID(barangID uint) error

	FindByIDWithComponents(id uint) (*model.BarangUnit, error)
	FindByBarangID(barangID uint) ([]model.BarangUnit, error)
}

type barangUnitRepository struct {
	db *gorm.DB
}

func NewBarangUnitRepository(db *gorm.DB) BarangUnitRepository {
	return &barangUnitRepository{db: db}
}

// ================= CREATE =================

func (r *barangUnitRepository) Create(unit *model.BarangUnit) error {
	return r.db.Create(unit).Error
}

func (r *barangUnitRepository) BulkCreate(units []model.BarangUnit) error {
	if len(units) == 0 {
		return nil
	}
	return r.db.Create(&units).Error
}

// ================= READ =================

func (r *barangUnitRepository) FindByID(id uint) (*model.BarangUnit, error) {
	var unit model.BarangUnit
	err := r.db.First(&unit, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unit tidak ditemukan")
		}
		return nil, err
	}
	return &unit, nil
}

func (r *barangUnitRepository) FindByKodeUnit(kode string) (*model.BarangUnit, error) {
	var unit model.BarangUnit
	err := r.db.Where("kode_unit = ?", kode).First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unit tidak ditemukan")
		}
		return nil, err
	}
	return &unit, nil
}

func (r *barangUnitRepository) FindAllByBarangID(barangID uint) ([]model.BarangUnit, error) {
	var units []model.BarangUnit
	err := r.db.Where("barang_id = ?", barangID).
		Order("kode_unit ASC").
		Find(&units).Error
	return units, err
}

// Ambil unit yang bisa dipinjam
func (r *barangUnitRepository) FindAvailableUnits(barangID uint, limit int) ([]model.BarangUnit, error) {
	var units []model.BarangUnit

	err := r.db.
		Where("barang_id = ? AND status IN ?", barangID, []string{"siap", "siap_terbatas"}).
		Order("status ASC, created_at ASC"). // "siap" diprioritaskan daripada "siap_terbatas"
		Limit(limit).
		Find(&units).Error

	return units, err
}

// Unit + komponennya (buat detail pinjaman)
func (r *barangUnitRepository) FindWithComponents(unitID uint) (*model.BarangUnit, error) {
	var unit model.BarangUnit
	err := r.db.
		Preload("Komponen").
		First(&unit, unitID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unit tidak ditemukan")
		}
		return nil, err
	}
	return &unit, nil
}

// ================= UPDATE =================

func (r *barangUnitRepository) Update(unit *model.BarangUnit) error {
	res := r.db.Save(unit)
	if res.RowsAffected == 0 {
		return errors.New("unit tidak ditemukan")
	}
	return res.Error
}

func (r *barangUnitRepository) UpdateStatus(unitID uint, status string) error {
	res := r.db.Model(&model.BarangUnit{}).
		Where("id = ?", unitID).
		Update("status", status)

	if res.RowsAffected == 0 {
		return errors.New("unit tidak ditemukan")
	}
	return res.Error
}

func (r *barangUnitRepository) UpdateKondisi(unitID uint, kondisi string) error {
	res := r.db.Model(&model.BarangUnit{}).
		Where("id = ?", unitID).
		Update("kondisi", kondisi)

	if res.RowsAffected == 0 {
		return errors.New("unit tidak ditemukan")
	}
	return res.Error
}

// ================= DELETE =================

func (r *barangUnitRepository) Delete(id uint) error {
	// Soft delete via GORM
	return r.db.Delete(&model.BarangUnit{}, id).Error
}

// ================= VALIDATION =================

func (r *barangUnitRepository) IsKodeUnitExists(kodeUnit string) (bool, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("kode_unit = ?", kodeUnit).
		Count(&count).Error
	return count > 0, err
}

func (r *barangUnitRepository) IsKodeUnitExistsExcludingID(kodeUnit string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("kode_unit = ? AND id != ?", kodeUnit, excludeID).
		Count(&count).Error
	return count > 0, err
}

// ================= COUNT =================

func (r *barangUnitRepository) CountTotal(barangID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ?", barangID).
		Count(&count).Error
	return count, err
}

func (r *barangUnitRepository) CountAvailable(barangID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ? AND status = ?", barangID, "siap").
		Count(&count).Error
	return count, err
}

func (r *barangUnitRepository) CountByStatus(barangID uint, status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ? AND status = ?", barangID, status).
		Count(&count).Error
	return count, err
}

// ================= STATS =================

func (r *barangUnitRepository) CountByKondisi(barangID uint, kondisi string) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ? AND kondisi = ?", barangID, kondisi).
		Count(&count).Error
	return count, err
}

func (r *barangUnitRepository) CountNeedsMaintenance(barangID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ? AND status = ?", barangID, "maintenance_required").
		Count(&count).Error
	return count, err
}

func (r *barangUnitRepository) FindByKodeWithComponents(kode string) (*model.BarangUnit, error) {
	var unit model.BarangUnit
	err := r.db.
		Preload("Komponen").
		Where("kode_unit = ?", kode).
		First(&unit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unit tidak ditemukan")
		}
		return nil, err
	}
	return &unit, nil
}

func (r *barangUnitRepository) UpdateStatusByKode(kode string, status string) error {
	res := r.db.Model(&model.BarangUnit{}).
		Where("kode_unit = ?", kode).
		Update("status", status)

	if res.RowsAffected == 0 {
		return errors.New("unit tidak ditemukan")
	}
	return res.Error
}

func (r *barangUnitRepository) UpdateKondisiByKode(kode string, kondisi string) error {
	res := r.db.Model(&model.BarangUnit{}).
		Where("kode_unit = ?", kode).
		Update("kondisi", kondisi)

	if res.RowsAffected == 0 {
		return errors.New("unit tidak ditemukan")
	}
	return res.Error
}

func (r *barangUnitRepository) DeleteByKode(kode string) error {
	res := r.db.
		Where("kode_unit = ?", kode).
		Delete(&model.BarangUnit{})

	if res.RowsAffected == 0 {
		return errors.New("unit tidak ditemukan")
	}
	return res.Error
}

func (r *barangUnitRepository) GetLastUnitNumber(barangID uint) (int, error) {
	var lastKode string

	// Unscoped → include soft-deleted
	err := r.db.Unscoped().
		Model(&model.BarangUnit{}).
		Select("kode_unit").
		Where("barang_id = ?", barangID).
		Order("kode_unit DESC").
		Limit(1).
		Scan(&lastKode).Error

	if err != nil {
		return 0, err
	}

	if lastKode == "" {
		return 0, nil // belum ada unit
	}

	// contoh: TND-07 → ambil 07
	parts := strings.Split(lastKode, "-")
	if len(parts) < 2 {
		return 0, nil
	}

	num, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, nil
	}

	return num, nil
}

func (r *barangUnitRepository) SetAllStatusByBarangID(barangID uint, status string) error {
	return r.db.Model(&model.BarangUnit{}).
		Where("barang_id = ?", barangID).
		Update("status", status).Error
}

func (r *barangUnitRepository) ReactivateAllByBarangID(barangID uint) error {
	// Ambil semua units
	var units []model.BarangUnit
	if err := r.db.Preload("Komponen").Where("barang_id = ?", barangID).Find(&units).Error; err != nil {
		return err
	}

	// Update satu-satu biar BeforeSave hook jalan (auto-compute status)
	for _, unit := range units {
		unit.Status = model.UnitStatusSiap // temporary, akan di-recompute
		if err := r.db.Save(&unit).Error; err != nil {
			return err
		}
	}

	return nil
}


func (r *barangUnitRepository) FindByIDWithComponents(id uint) (*model.BarangUnit, error) {
	var unit model.BarangUnit
	err := r.db.Preload("Komponen").First(&unit, id).Error
	return &unit, err
}

func (r *barangUnitRepository) FindByBarangID(barangID uint) ([]model.BarangUnit, error) {
	var units []model.BarangUnit
	err := r.db.Where("barang_id = ?", barangID).Find(&units).Error
	return units, err
}