package repo

import (
	"backend/model"
	"errors"

	"gorm.io/gorm"
)

type AngkatanMapalaRepo interface {
    Create(f *model.AngkatanMapala) error
    GetAll() ([]model.AngkatanMapala, error)
    GetByID(id uint) (*model.AngkatanMapala, error)
    Update(f *model.AngkatanMapala) error
    Delete(id uint) error
    GetByName(name string) (*model.AngkatanMapala, error)
}

type AngkatanMapalaRepoDB struct {
    db *gorm.DB
}

func NewAngkatanMapalaRepo(db *gorm.DB) AngkatanMapalaRepo {
    return &AngkatanMapalaRepoDB{db: db}
}

func (r *AngkatanMapalaRepoDB) Create(f *model.AngkatanMapala) error {
    return r.db.Create(f).Error
}

func (r *AngkatanMapalaRepoDB) GetAll() ([]model.AngkatanMapala, error) {
    var data []model.AngkatanMapala
    err := r.db.Find(&data).Error
    return data, err
}

func (r *AngkatanMapalaRepoDB) GetByID(id uint) (*model.AngkatanMapala, error) {
    var data model.AngkatanMapala
    err := r.db.First(&data, id).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New("angkatan mapala tidak ditemukan")
    }

    return &data, err
}

func (r *AngkatanMapalaRepoDB) Update(f *model.AngkatanMapala) error {
    return r.db.Save(f).Error
}

func (r *AngkatanMapalaRepoDB) Delete(id uint) error {
    return r.db.Delete(&model.AngkatanMapala{}, id).Error
}
func (r *AngkatanMapalaRepoDB) GetByName(name string) (*model.AngkatanMapala, error) {
    var a model.AngkatanMapala
    if err := r.db.Where("nama = ?", name).First(&a).Error; err != nil {
        return nil, err
    }
    return &a, nil
}