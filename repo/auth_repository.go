package repo

import (
	"backend/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *model.User) error
	FindUserByNra(nra string) (*model.User, error)
	FindUserByID(id uint) (*model.User, error)
	UpdateUser(user *model.User) error
	
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindUserByNra(nra string) (*model.User, error) {
	var user model.User
	err := r.db.Where("nra = ?", nra).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}
