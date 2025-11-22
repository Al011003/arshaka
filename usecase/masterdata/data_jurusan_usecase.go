package usecase

import (
	req "backend/dto/request/masterdata"
	res "backend/dto/response/masterdata"
	"backend/model"
	"backend/repo"
	"errors"
)

type JurusanUsecase interface {
	Create(req req.JurusanRequest) (*res.JurusanResponse, error)
	GetAll() ([]res.JurusanResponse, error)
	Update(id uint, req req.JurusanRequest) (*res.JurusanResponse, error)
	Delete(id uint) error
	GetByFakultas(fakultasID uint) ([]model.Jurusan, error)

}

type jurusanUsecase struct {
	jurusanRepo  repo.JurusanRepo
	fakultasRepo repo.FakultasRepo
}

func NewJurusanUsecase(j repo.JurusanRepo, f repo.FakultasRepo) JurusanUsecase {
	return &jurusanUsecase{
		jurusanRepo:  j,
		fakultasRepo: f,
	}
}

// Create
func (u *jurusanUsecase) Create(req req.JurusanRequest) (*res.JurusanResponse, error) {
	// cek fakultas
	f, err := u.fakultasRepo.GetByID(req.FakultasID)
	if err != nil || f == nil {
		return nil, errors.New("fakultas tidak ditemukan")
	}

	data := model.Jurusan{
		Nama:       req.Nama,
		FakultasID: req.FakultasID,
	}

	if err := u.jurusanRepo.Create(&data); err != nil {
		return nil, err
	}

	return &res.JurusanResponse{
		ID:         data.ID,
		Nama:       data.Nama,
		FakultasID: data.FakultasID,
	}, nil
}

// GetAll
func (u *jurusanUsecase) GetAll() ([]res.JurusanResponse, error) {
	list, err := u.jurusanRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]res.JurusanResponse, len(list))
	for i, j := range list {
		result[i] = res.JurusanResponse{
			ID:         j.ID,
			Nama:       j.Nama,
			FakultasID: j.FakultasID,
		}
	}

	return result, nil
}

// Update
func (u *jurusanUsecase) Update(id uint, req req.JurusanRequest) (*res.JurusanResponse, error) {
	data, err := u.jurusanRepo.GetByID(id)
	if err != nil || data == nil {
		return nil, errors.New("jurusan tidak ditemukan")
	}

	// cek fakultas
	f, err := u.fakultasRepo.GetByID(req.FakultasID)
	if err != nil || f == nil {
		return nil, errors.New("fakultas tidak ditemukan")
	}

	data.Nama = req.Nama
	data.FakultasID = req.FakultasID

	if err := u.jurusanRepo.Update(data); err != nil {
		return nil, err
	}

	return &res.JurusanResponse{
		ID:         data.ID,
		Nama:       data.Nama,
		FakultasID: data.FakultasID,
	}, nil
}

// Delete
func (u *jurusanUsecase) Delete(id uint) error {
	return u.jurusanRepo.Delete(id)
}

func (u *jurusanUsecase) GetByFakultas(fakultasID uint) ([]model.Jurusan, error) {
	return u.jurusanRepo.GetByFakultasID(fakultasID)
}
