package usecase

import (
	service "backend/email"
	"backend/model"
	"backend/repo"
	"backend/utils"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type PasswordResetUsecase interface {
	SendOTP(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, otp, newPassword string) error
}

type passwordResetUsecase struct {
	userRepo     repo.UserRepo
	prRepo       repo.PasswordResetRepository
	emailService service.EmailService
}

func NewPasswordResetUsecase(
	userRepo repo.UserRepo,
	prRepo repo.PasswordResetRepository,
	emailService service.EmailService,
) PasswordResetUsecase {
	return &passwordResetUsecase{
		userRepo:     userRepo,
		prRepo:       prRepo,
		emailService: emailService,
	}
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (u *passwordResetUsecase) SendOTP(ctx context.Context, email string) error {

	// 1. Cek apakah user dengan email tsb ada
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	// 2. CEK COOLDOWN - Cek apakah ada OTP yang baru dibuat (< 1 menit yang lalu)
	lastOTP, err := u.prRepo.GetLatestOTPByUserID(ctx, user.ID)
	if err == nil && lastOTP != nil {
		// Hitung berapa lama sejak OTP terakhir dibuat
		timeSinceLastOTP := time.Since(lastOTP.CreatedAt)
		
		// Jika belum lewat 1 menit, tolak request
		if timeSinceLastOTP < 1*time.Minute {
			remainingSeconds := int(60 - timeSinceLastOTP.Seconds())
			return fmt.Errorf("mohon tunggu %d detik lagi sebelum request OTP baru", remainingSeconds)
		}
	}

	// 3. EXPIRE OTP LAMA - Tandai semua OTP lama sebagai expired/used
	if err := u.prRepo.ExpireAllOTPByUserID(ctx, user.ID); err != nil {
		return fmt.Errorf("gagal mengexpire OTP lama: %w", err)
	}

	// 4. Generate OTP baru
	otp := generateOTP()

	// 5. Simpan OTP baru
	pr := &model.PasswordReset{
		UserID:    user.ID,
		OTP:       otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Used:      false,
		CreatedAt: time.Now(), // PENTING: pastikan field ini ada
	}

	if err := u.prRepo.Create(ctx, pr); err != nil {
		return fmt.Errorf("gagal menyimpan OTP: %w", err)
	}

	// 6. Kirim email
	if err := u.emailService.SendOTP(email, user.Nama, otp); err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}

func (u *passwordResetUsecase) ResetPassword(ctx context.Context, email, otp, newPassword string) error {

	// 1. cek user berdasarkan email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	// 2. ambil OTP valid
	pr, err := u.prRepo.GetValidOTP(ctx, user.ID, otp)
	if err != nil {
		return errors.New("otp salah atau sudah kadaluarsa")
	}

	// 3. hash password baru
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("gagal menghash password: %w", err)
	}

	// 4. update password user
	if err := u.userRepo.UpdatePassword(user.ID, hashedPassword, false); err != nil {
		return fmt.Errorf("gagal mengupdate password: %w", err)
	}

	// 5. tandai OTP sudah dipakai
	if err := u.prRepo.MarkUsed(ctx, pr.ID); err != nil {
		return fmt.Errorf("gagal menandai OTP sebagai used: %w", err)
	}

	return nil
}