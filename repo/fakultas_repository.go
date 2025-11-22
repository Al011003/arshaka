package repo

import (
	"backend/model"
	"errors"

	"gorm.io/gorm"
)


type FakultasRepo interface {
Create(f *model.Fakultas) error
GetAll() ([]model.Fakultas, error)
GetByID(id uint) (*model.Fakultas, error)
Update(f *model.Fakultas) error
Delete(id uint) error
GetByName(name string) (*model.Fakultas, error)
}


type fakultasRepoDB struct { db *gorm.DB }


func NewFakultasRepo(db *gorm.DB) FakultasRepo { return &fakultasRepoDB{db} }


func (r *fakultasRepoDB) Create(f *model.Fakultas) error { return r.db.Create(f).Error }
func (r *fakultasRepoDB) GetAll() ([]model.Fakultas, error) {
var data []model.Fakultas
err := r.db.Preload("Jurusan").Find(&data).Error
return data, err
}
func (r *fakultasRepoDB) GetByID(id uint) (*model.Fakultas, error) {
var f model.Fakultas
err := r.db.Preload("Jurusan").First(&f, id).Error
if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
return &f, err
}
func (r *fakultasRepoDB) Update(f *model.Fakultas) error { return r.db.Save(f).Error }
func (r *fakultasRepoDB) Delete(id uint) error { return r.db.Delete(&model.Fakultas{}, id).Error }
func (r *fakultasRepoDB) GetByName(name string) (*model.Fakultas, error) {
    var f model.Fakultas
    if err := r.db.Where("nama = ?", name).First(&f).Error; err != nil {
        return nil, err
    }
    return &f, nil
}



