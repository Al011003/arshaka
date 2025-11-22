package usecase

import (
	req "backend/dto/request/masterdata"
	res "backend/dto/response/masterdata"
	"backend/model"
	"backend/repo"
	"errors"
)

type FakultasUsecase interface {
	Create(req req.FakultasRequest) (*res.FakultasResponse, error)
	GetAll() ([]res.FakultasResponse, error)
	Update(id uint, req req.FakultasRequest) (*res.FakultasResponse, error)
	Delete(id uint) error
}

type fakultasUsecase struct {
	repo repo.FakultasRepo
}

func NewFakultasUsecase(r repo.FakultasRepo) FakultasUsecase {
	return &fakultasUsecase{
		repo: r,
	}
}

func (u *fakultasUsecase) Create(req req.FakultasRequest) (*res.FakultasResponse, error) {
	data := model.Fakultas{
		Nama: req.Nama,
	}

	if err := u.repo.Create(&data); err != nil {
		return nil, err
	}

	return &res.FakultasResponse{
		ID:   data.ID,
		Nama: data.Nama,
	}, nil
}

func (u *fakultasUsecase) GetAll() ([]res.FakultasResponse, error) {
	list, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]res.FakultasResponse, len(list))
	for i, f := range list {
		result[i] = res.FakultasResponse{
			ID:   f.ID,
			Nama: f.Nama,
		}
	}

	return result, nil
}

func (u *fakultasUsecase) Update(id uint, req req.FakultasRequest) (*res.FakultasResponse, error) {
	data, err := u.repo.GetByID(id)
	if err != nil || data == nil {
		return nil, errors.New("fakultas tidak ditemukan")
	}

	data.Nama = req.Nama
	if err := u.repo.Update(data); err != nil {
		return nil, err
	}

	return &res.FakultasResponse{
		ID:   data.ID,
		Nama: data.Nama,
	}, nil
}

func (u *fakultasUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
