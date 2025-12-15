package repo

import (
	user "backend/dto/request/user"
	"backend/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetPaginated(page, limit int) ([]model.User, int, error)
	GetByID(id uint) (*model.User, error)
	SearchPaginated(keyword string, page, limit int) ([]model.User, int, error)
	Update(u *model.User) error
	Delete(id uint) error
	FindUsers(filter user.UserFilter) ([]model.User, int, error)
    UpdatePassword(userID uint, newPassword string, mustChange bool) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
    UpdateEmail(ctx context.Context, userID uint, email string) error
    ResetPassword(userID uint, newPassword string, mustChange bool) error
    UpdateProfilePhoto(userID uint, photoURL string) error
    FindByNRAAndName(nra string, nama string) (*model.User, error)
}

type UserRepoDB struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UserRepoDB{db: db}
}

func (r *UserRepoDB) GetPaginated(page, limit int) ([]model.User, int, error) {
    var users []model.User
    offset := (page - 1) * limit

    // total data
    var total int64
    r.db.Model(&model.User{}).Count(&total)

    // ambil data berdasar page
    err := r.db.
        Order("id DESC").
        Limit(limit).
        Offset(offset).
        Find(&users).Error

    return users, int(total), err
}

func (r *UserRepoDB) GetByID(id uint) (*model.User, error) {
	var user model.User

	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}


func (r *UserRepoDB) SearchPaginated(keyword string, page, limit int) ([]model.User, int, error) {
	var users []model.User

	offset := (page - 1) * limit

	query := r.db.Model(&model.User{})

	if keyword != "" {
		query = query.Where(
			"nama LIKE ? OR nama_lengkap LIKE ? OR nra LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	// total rows
	var total int64
	query.Count(&total)

	// data pagination
	err := query.Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, int(total), err
}



func (r *UserRepoDB) Update(u *model.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepoDB) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepoDB) FindUsers(filter user.UserFilter) ([]model.User, int, error) {
    var users []model.User
    db := r.db.Model(&model.User{})
	roles := strings.Split(filter.ExcludeRole, ",")


    // Exclude role - buat permission (dari usecase)
       if len(roles) == 1 {
        fmt.Println("Using != query")
        db = db.Where("role != ?", roles[0])
    } else {
        fmt.Println("Using NOT IN query")
        db = db.Where("role NOT IN ?", roles)
    }

    // Universal search - cari di SEMUA field termasuk role
   if filter.Search != "" {
    keyword := "%" + filter.Search + "%"
    db = db.Where(
        "nama LIKE ? OR nama_lengkap LIKE ? OR nra LIKE ? OR fakultas LIKE ? OR jurusan LIKE ? OR angkatan_mapala LIKE ? OR role LIKE ?",
        keyword, keyword, keyword, keyword, keyword, keyword, keyword,
        // ^ 7 placeholder, 7 values ✅
    )
}

    var total int64
    db.Count(&total)

   // Default pagination values
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    // Apply pagination
    offset := (filter.Page - 1) * filter.Limit
    err := db.Order("id DESC").
        Limit(filter.Limit).
        Offset(offset).
        Find(&users).Error
    
    return users, int(total), err
}

func (r *UserRepoDB) UpdatePassword(userID uint, newPassword string, mustChange bool) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":            newPassword,
			"is_default_password": false,
		}).Error
}

func (r *UserRepoDB) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepoDB) UpdateEmail(ctx context.Context, userID uint, email string) error {
    return r.db.WithContext(ctx).
        Model(&model.User{}).
        Where("id = ?", userID).
        Updates(map[string]interface{}{
            "email":         email,
            "is_email_null": false,
        }).Error
}

func (r *UserRepoDB) ResetPassword(userID uint, newPassword string, mustChange bool) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":            newPassword,
			"is_default_password": true,
		}).Error
}

func (r *UserRepoDB) UpdateProfilePhoto(userID uint, photoURL string) error {
    return r.db.Model(&model.User{}).
        Where("id = ?", userID).
        Update("foto_url", photoURL).
        Error
}

func (r *UserRepoDB) FindByNRAAndName(nra string, nama string) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("nra = ? AND nama LIKE ?", nra, "%"+nama+"%").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
