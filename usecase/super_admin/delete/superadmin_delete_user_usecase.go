package usecase

import (
	"backend/repo"
	"errors"
)

type SuperAdminDeleteUserUsecase interface {
	DeleteUser(superAdminID, targetID uint) error
}

type superAdminDeleteUserUC struct {
	userRepo repo.UserRepo
}

// Constructor
func NewSuperAdminDeleteUserUsecase(r repo.UserRepo) SuperAdminDeleteUserUsecase {
	return &superAdminDeleteUserUC{userRepo: r}
}

func (u *superAdminDeleteUserUC) DeleteUser(superAdminID, targetID uint) error {

	// 1. Ambil data superadmin
	admin, err := u.userRepo.GetByID(superAdminID)
	if err != nil || admin == nil {
		return errors.New("super admin tidak ditemukan")
	}

	// 2. Validasi role
	if admin.Role != "superadmin" {
		return errors.New("forbidden: hanya super admin yang bisa menghapus user")
	}

	// 3. Super admin tidak boleh hapus dirinya sendiri
	if superAdminID == targetID {
		return errors.New("super admin tidak boleh menghapus dirinya sendiri")
	}

	// 4. Cek user target
	targetUser, err := u.userRepo.GetByID(targetID)
	if err != nil || targetUser == nil {
		return errors.New("user tidak ditemukan")
	}

	// 5. Eksekusi delete
	if err := u.userRepo.Delete(targetID); err != nil {
		return errors.New("gagal menghapus user")
	}

	return nil
}
