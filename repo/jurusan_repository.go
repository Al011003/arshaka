package repo

import (
	"backend/model"
	"errors"

	"gorm.io/gorm"
)


type JurusanRepo interface {
Create(j *model.Jurusan) error
GetAll() ([]model.Jurusan, error)
GetByID(id uint) (*model.Jurusan, error)
Update(j *model.Jurusan) error
Delete(id uint) error
GetByFakultasID(fakultasID uint) ([]model.Jurusan, error)
GetByName(name string) (*model.Jurusan, error)
}


type jurusanRepoDB struct { db *gorm.DB }


func NewJurusanRepo(db *gorm.DB) JurusanRepo { return &jurusanRepoDB{db} }


func (r *jurusanRepoDB) Create(j *model.Jurusan) error { return r.db.Create(j).Error }
func (r *jurusanRepoDB) GetAll() ([]model.Jurusan, error) {
var data []model.Jurusan
err := r.db.Find(&data).Error
return data, err
}
func (r *jurusanRepoDB) GetByID(id uint) (*model.Jurusan, error) {
var j model.Jurusan
err := r.db.First(&j, id).Error
if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
return &j, err
}
func (r *jurusanRepoDB) Update(j *model.Jurusan) error { return r.db.Save(j).Error }
func (r *jurusanRepoDB) Delete(id uint) error { return r.db.Delete(&model.Jurusan{}, id).Error }
func (r *jurusanRepoDB) GetByFakultasID(fakultasID uint) ([]model.Jurusan, error) {
	var data []model.Jurusan
	err := r.db.Where("fakultas_id = ?", fakultasID).Find(&data).Error
	return data, err
}
func (r *jurusanRepoDB) GetByName(name string) (*model.Jurusan, error) {
    var j model.Jurusan
    if err := r.db.Where("nama = ?", name).First(&j).Error; err != nil {
        return nil, err
    }
    return &j, nil
}