package usecase

import (
	req "backend/dto/request/masterdata"
	res "backend/dto/response/masterdata"
	"backend/model"
	"backend/repo"
	"errors"
)

type AngkatanMapalaUsecase interface {
    Create(req req.AngkatanMapalaRequest) (*res.AngkatanMapalaResponse, error)
    GetAll() ([]res.AngkatanMapalaResponse, error)
    Update(id uint, req req.AngkatanMapalaRequest) (*res.AngkatanMapalaResponse, error)
    Delete(id uint) error
}

type angkatanMapalaUsecase struct {
    repo repo.AngkatanMapalaRepo
}

func NewAngkatanMapalaUsecase(r repo.AngkatanMapalaRepo) AngkatanMapalaUsecase {
    return &angkatanMapalaUsecase{
        repo: r,
    }
}

func (u *angkatanMapalaUsecase) Create(req req.AngkatanMapalaRequest) (*res.AngkatanMapalaResponse, error) {
    data := model.AngkatanMapala{
        Nama: req.Nama,
    }

    if err := u.repo.Create(&data); err != nil {
        return nil, err
    }

    return &res.AngkatanMapalaResponse{
        ID:   data.ID,
        Nama: data.Nama,
    }, nil
}

func (u *angkatanMapalaUsecase) GetAll() ([]res.AngkatanMapalaResponse, error) {
    list, err := u.repo.GetAll()
    if err != nil {
        return nil, err
    }

    result := make([]res.AngkatanMapalaResponse, len(list))
    for i, f := range list {
        result[i] = res.AngkatanMapalaResponse{
            ID:   f.ID,
            Nama: f.Nama,
        }
    }

    return result, nil
}

func (u *angkatanMapalaUsecase) Update(id uint, req req.AngkatanMapalaRequest) (*res.AngkatanMapalaResponse, error) {
    data, err := u.repo.GetByID(id)
    if err != nil || data == nil {
        return nil, errors.New("angkatan mapala tidak ditemukan")
    }

    data.Nama = req.Nama

    if err := u.repo.Update(data); err != nil {
        return nil, err
    }

    return &res.AngkatanMapalaResponse{
        ID:   data.ID,
        Nama: data.Nama,
    }, nil
}

func (u *angkatanMapalaUsecase) Delete(id uint) error {
    return u.repo.Delete(id)
}
