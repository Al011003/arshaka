package usecase

import (
	request "backend/dto/request/barang"
	response "backend/dto/response/barang"
	"backend/model"
	"backend/repo"
	"errors"
)

type BarangKomponenUseCase interface {
	Add(kodeUnit string, req request.CreateKomponenRequest) error
	Update(id uint, req request.UpdateKomponenRequest) error
	Delete(id uint) error
	GetByUnitKode(kodeUnit string) ([]response.KomponenResponse, error)
}

type barangKomponenUseCase struct {
	komponenRepo repo.BarangKomponenRepository
	unitRepo     repo.BarangUnitRepository
}

func NewBarangKomponenUseCase(
	komponenRepo repo.BarangKomponenRepository,
	unitRepo repo.BarangUnitRepository,
) BarangKomponenUseCase {
	return &barangKomponenUseCase{
		komponenRepo: komponenRepo,
		unitRepo:     unitRepo,
	}
}

func (u *barangKomponenUseCase) Add(kodeUnit string, req request.CreateKomponenRequest) error {
	// Cari unit berdasarkan kode
	unit, err := u.unitRepo.FindByKodeUnit(kodeUnit)
	if err != nil {
		return errors.New("unit tidak ditemukan")
	}

	// Cek apakah unit aktif
	if unit.Status == model.UnitStatusNonAktif {
		return errors.New("tidak bisa menambah komponen: unit tidak aktif")
	}

	// ← TAMBAHAN: Cek duplikasi nama
	exists, err := u.komponenRepo.ExistsByUnitAndNama(unit.ID, req.Nama)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("komponen dengan nama ini sudah ada di unit")
	}

	comp := &model.BarangKomponen{
		BarangUnitID: unit.ID,
		Nama:         req.Nama,
		JumlahWajib:  req.JumlahWajib,
		JumlahAktual: req.JumlahWajib,
		HargaSatuan:  req.HargaSatuan,
	}

	return u.komponenRepo.Create(comp)
}

func (u *barangKomponenUseCase) Update(id uint, req request.UpdateKomponenRequest) error {
	comp, err := u.komponenRepo.FindByID(id)
	if err != nil {
		return errors.New("komponen tidak ditemukan")
	}

	// Update fields
	if req.Nama != nil {
		if *req.Nama != comp.Nama {
			exists, err := u.komponenRepo.ExistsByUnitAndNama(comp.BarangUnitID, *req.Nama)
			if err != nil {
				return err
			}
			if exists {
				return errors.New("komponen dengan nama ini sudah ada di unit")
			}
		}
		comp.Nama = *req.Nama
	}

	if req.JumlahWajib != nil {
		comp.JumlahWajib = *req.JumlahWajib
	}

	if req.JumlahAktual != nil {
		comp.JumlahAktual = *req.JumlahAktual
	}

	if req.HargaSatuan != nil {
		comp.HargaSatuan = *req.HargaSatuan
	}

	// Save komponen
	if err := u.komponenRepo.Update(comp); err != nil {
		return err
	}

	// ← FIX: Trigger update unit dengan ID
	unit, err := u.unitRepo.FindByIDWithComponents(comp.BarangUnitID)
	if err != nil {
		return err
	}

	// Trigger save biar BeforeSave hook jalan & status unit ke-recompute
	return u.unitRepo.Update(unit)
}

func (u *barangKomponenUseCase) Delete(id uint) error {
	comp, err := u.komponenRepo.FindByID(id)
	if err != nil {
		return errors.New("komponen tidak ditemukan")
	}

	unitID := comp.BarangUnitID

	// Delete komponen
	if err := u.komponenRepo.Delete(comp.ID); err != nil {
		return err
	}

	// Trigger update unit biar status ke-recompute
	unit, err := u.unitRepo.FindByIDWithComponents(unitID)
	if err != nil {
		return err
	}

	return u.unitRepo.Update(unit)
}

func (u *barangKomponenUseCase) GetByUnitKode(kodeUnit string) ([]response.KomponenResponse, error) {
	// Cari unit berdasarkan kode
	unit, err := u.unitRepo.FindByKodeUnit(kodeUnit)
	if err != nil {
		return nil, errors.New("unit tidak ditemukan")
	}

	comps, err := u.komponenRepo.FindByUnit(unit.ID)
	if err != nil {
		return nil, err
	}

	var out []response.KomponenResponse
	for _, c := range comps {
		
		out = append(out, response.KomponenResponse{
			ID:                c.ID,
			Nama:              c.Nama,
			JumlahWajib:       c.JumlahWajib,
			JumlahAktual:      c.JumlahAktual,
			Status:            c.Status,
			HargaSatuan:       c.HargaSatuan,
		})
	}

	return out, nil
}