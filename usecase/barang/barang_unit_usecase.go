package usecase

import (
	request "backend/dto/request/barang"
	response "backend/dto/response/barang"
	"backend/model"
	"backend/repo"
	"errors"
	"fmt"
	"time"
)

type BarangUnitUseCase interface {
	Create(req request.CreateUnitRequest) (*response.UnitCreateResponse, error)
	GetByKode(kode string) (*response.UnitDetailResponse, error)
	GetByBarangKode(kode string) ([]response.UnitListResponse, error)
	
	// Update unit info (metadata)
	UpdateByKode(kode string, req request.UpdateUnitRequest) (*response.UnitDetailResponse, error)
	
	// Explicit status actions
	SetMaintenance(kode string, priority int, alasan string) error
	SetNonAktif(kode string, alasan string) error
	SetAktifKembali(kode string) error
	
	DeleteByKode(kode string) error
}

type barangUnitUseCase struct {
	unitRepo     repo.BarangUnitRepository
	barangRepo   repo.BarangRepository
	komponenRepo repo.BarangKomponenRepository
}

func NewBarangUnitUseCase(
	unitRepo repo.BarangUnitRepository,
	barangRepo repo.BarangRepository,
	komponenRepo repo.BarangKomponenRepository,
) BarangUnitUseCase {
	return &barangUnitUseCase{
		unitRepo:     unitRepo,
		barangRepo:   barangRepo,
		komponenRepo: komponenRepo,
	}
}

func (u *barangUnitUseCase) Create(req request.CreateUnitRequest) (*response.UnitCreateResponse, error) {
	barang, err := u.barangRepo.FindByKode(req.KodeBarang)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	// Cek apakah barang aktif
	if barang.Status != model.BarangStatusAktif {
		return nil, errors.New("tidak bisa menambah unit: barang tidak aktif")
	}

	jumlah := req.Jumlah
	if jumlah <= 0 {
		jumlah = 1
	}

	lastNumber, err := u.unitRepo.GetLastUnitNumber(barang.ID)
	if err != nil {
		return nil, err
	}

	var units []model.BarangUnit
	var kodeUnits []string

	for i := 1; i <= jumlah; i++ {
		no := lastNumber + i
		kodeUnit := fmt.Sprintf("%s-%02d", barang.Kode, no)

		units = append(units, model.BarangUnit{
			BarangID:        barang.ID,
			KodeUnit:        kodeUnit,
			Status:          model.UnitStatusSiap,
			Kondisi:         model.UnitKondisiBaik,
			TahunPerolehan:  req.TahunPerolehan,
			LokasiPenyimpan: req.LokasiPenyimpan,
			Catatan:         req.Catatan,
		})

		kodeUnits = append(kodeUnits, kodeUnit)
	}

	if err := u.unitRepo.BulkCreate(units); err != nil {
		return nil, err
	}

	return &response.UnitCreateResponse{
		Message:  fmt.Sprintf("%d unit %s berhasil ditambahkan", jumlah, barang.Kode),
		KodeUnit: kodeUnits,
	}, nil
}

func (u *barangUnitUseCase) GetByKode(kode string) (*response.UnitDetailResponse, error) {
	unit, err := u.unitRepo.FindByKodeWithComponents(kode)
	if err != nil {
		return nil, err
	}

	barang, err := u.barangRepo.FindByID(unit.BarangID)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	return &response.UnitDetailResponse{
		ID:              unit.ID,
		BarangKode:      barang.Kode,
		BarangNama:      barang.Nama,
		KodeUnit:        unit.KodeUnit,
		Status:          unit.Status,
		Kondisi:         unit.Kondisi,
		FixPriority:     unit.FixPriority,
		TahunPerolehan:  unit.TahunPerolehan,
		LokasiPenyimpan: unit.LokasiPenyimpan,
		Catatan:         unit.Catatan,
		// Komponen bisa ditambahin di response nanti
	}, nil
}

func (u *barangUnitUseCase) GetByBarangKode(kode string) ([]response.UnitListResponse, error) {
	barang, err := u.barangRepo.FindByKode(kode)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	units, err := u.unitRepo.FindAllByBarangID(barang.ID)
	if err != nil {
		return nil, err
	}

	var out []response.UnitListResponse
	for _, unit := range units {
		out = append(out, response.UnitListResponse{
			ID:       unit.ID,
			KodeUnit: unit.KodeUnit,
			Status:   unit.Status,
			Kondisi:  unit.Kondisi,
		})
	}

	return out, nil
}

// UpdateByKode - update metadata unit (kondisi, tahun, lokasi, catatan)
func (u *barangUnitUseCase) UpdateByKode(kode string, req request.UpdateUnitRequest) (*response.UnitDetailResponse, error) {
	unit, err := u.unitRepo.FindByKodeWithComponents(kode)
	if err != nil {
		return nil, err
	}

	// Update fields yang boleh diubah
	if req.Kondisi != nil {
		// Validate kondisi
		if *req.Kondisi != model.UnitKondisiBaik && 
		   *req.Kondisi != model.UnitKondisiLecet && 
		   *req.Kondisi != model.UnitKondisiRusak {
			return nil, errors.New("kondisi tidak valid")
		}
		unit.Kondisi = *req.Kondisi
	}

	if req.TahunPerolehan != nil {
		unit.TahunPerolehan = *req.TahunPerolehan
	}

	if req.LokasiPenyimpan != nil {
		unit.LokasiPenyimpan = *req.LokasiPenyimpan
	}

	if req.Catatan != nil {
		unit.Catatan = *req.Catatan
	}

	// Status akan auto-compute di BeforeSave hook
	if err := u.unitRepo.Update(unit); err != nil {
		return nil, err
	}

	return u.GetByKode(kode)
}

// SetMaintenance - set unit ke status maintenance dengan priority
// usecase/barang/barang_unit.go

// SetMaintenance - set unit ke status maintenance dengan priority
func (u *barangUnitUseCase) SetMaintenance(kode string, priority int, alasan string) error {
	unit, err := u.unitRepo.FindByKodeUnit(kode)
	if err != nil {
		return err
	}

	if priority < model.FixPriorityLow || priority > model.FixPriorityHigh {
		return errors.New("fix priority harus antara 1 (low), 2 (medium), atau 3 (high)")
	}

	if unit.Status == model.UnitStatusDipinjam {
		return errors.New("unit sedang dipinjam, tidak bisa di-maintenance")
	}

	now := time.Now()
	priorityText := getPriorityText(priority)
	
	unit.FixPriority = priority
	unit.Status = model.UnitStatusMaintenanceRequired
	
	// OVERWRITE catatan dengan status terbaru
	unit.Catatan = fmt.Sprintf("[%s] Maintenance (%s): %s", 
		now.Format("2006-01-02 15:04"), 
		priorityText,
		alasan,
	)

	return u.unitRepo.Update(unit)
}

func (u *barangUnitUseCase) SetNonAktif(kode string, alasan string) error {
	unit, err := u.unitRepo.FindByKodeUnit(kode)
	if err != nil {
		return err
	}

	if unit.Status == model.UnitStatusDipinjam {
		return errors.New("unit sedang dipinjam, tidak bisa dinonaktifkan")
	}

	now := time.Now()
	unit.Status = model.UnitStatusNonAktif
	unit.FixPriority = model.FixPriorityNormal
	
	// OVERWRITE catatan
	unit.Catatan = fmt.Sprintf("[%s] Dinonaktifkan: %s", 
		now.Format("2006-01-02 15:04"),
		alasan,
	)

	return u.unitRepo.Update(unit)
}

func (u *barangUnitUseCase) SetAktifKembali(kode string) error {
	unit, err := u.unitRepo.FindByKodeWithComponents(kode)
	if err != nil {
		return err
	}

	if unit.Status != model.UnitStatusNonAktif && unit.Status != model.UnitStatusMaintenanceRequired {
		return errors.New("unit sudah aktif")
	}

	now := time.Now()
	unit.FixPriority = model.FixPriorityNormal
	unit.Status = unit.HitungStatus()
	
	// OVERWRITE catatan
	unit.Catatan = fmt.Sprintf("[%s] Barang aktif kembali", 
		now.Format("2006-01-02 15:04"),
	)

	return u.unitRepo.Update(unit)
}

func getPriorityText(priority int) string {
	switch priority {
	case model.FixPriorityLow:
		return "Low"
	case model.FixPriorityMedium:
		return "Medium"
	case model.FixPriorityHigh:
		return "High"
	default:
		return "Normal"
	}
}

func (u *barangUnitUseCase) DeleteByKode(kode string) error {
	unit, err := u.unitRepo.FindByKodeUnit(kode)
	if err != nil {
		return err
	}

	// Cek apakah sedang dipinjam
	if unit.Status == model.UnitStatusDipinjam {
		return errors.New("unit sedang dipinjam, tidak bisa dihapus")
	}

	return u.unitRepo.DeleteByKode(kode)
}