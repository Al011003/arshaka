package usecase

import (
	service "backend/email"
	"backend/model"
	"backend/repo"
	"backend/utils"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"time"
)

type PasswordResetUsecase interface {
	SendOTP(ctx context.Context, email string) error
	VerifyOTPAndGenerateToken(ctx context.Context, email, otp string) (resetToken string, expiresIn int, err error)
	ResetPasswordWithToken(ctx context.Context, resetToken, newPassword string) error
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
	mathrand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", mathrand.Intn(1000000))
}

// Generate secure random token
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "reset_" + hex.EncodeToString(bytes), nil
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
		timeSinceLastOTP := time.Since(lastOTP.CreatedAt)
		
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
		CreatedAt: time.Now(),
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

// Verify OTP dan generate reset token
func (u *passwordResetUsecase) VerifyOTPAndGenerateToken(ctx context.Context, email, otp string) (string, int, error) {
	// 1. Cek user berdasarkan email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, errors.New("email tidak ditemukan")
	}

	// 2. Validasi OTP
	pr, err := u.prRepo.GetValidOTP(ctx, user.ID, otp)
	if err != nil {
		return "", 0, errors.New("otp salah atau sudah kadaluarsa")
	}

	// 3. Generate reset token
	resetToken, err := generateResetToken()
	if err != nil {
		return "", 0, fmt.Errorf("gagal generate reset token: %w", err)
	}

	// 4. Token expires dalam 10 menit
	tokenExpiresAt := time.Now().Add(10 * time.Minute)
	expiresIn := 600 // 10 menit dalam detik

	// 5. Update record dengan reset token
	if err := u.prRepo.UpdateResetToken(ctx, pr.ID, resetToken, tokenExpiresAt); err != nil {
		return "", 0, fmt.Errorf("gagal menyimpan reset token: %w", err)
	}

	// 6. Mark OTP as used (opsional, bisa juga biarkan tetap valid sampai token dipakai)
	if err := u.prRepo.MarkUsed(ctx, pr.ID); err != nil {
		return "", 0, fmt.Errorf("gagal menandai OTP sebagai used: %w", err)
	}

	return resetToken, expiresIn, nil
}

// Reset password menggunakan token
func (u *passwordResetUsecase) ResetPasswordWithToken(ctx context.Context, resetToken, newPassword string) error {
	// 1. Validasi token
	pr, err := u.prRepo.GetByResetToken(ctx, resetToken)
	if err != nil {
		return errors.New("token tidak valid atau sudah kadaluarsa")
	}

	// 2. Hash password baru
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("gagal menghash password: %w", err)
	}

	// 3. Update password user
	if err := u.userRepo.UpdatePassword(pr.UserID, hashedPassword, false); err != nil {
		return fmt.Errorf("gagal mengupdate password: %w", err)
	}

	// 4. Invalidate token
	if err := u.prRepo.MarkTokenUsed(ctx, pr.ID); err != nil {
		return fmt.Errorf("gagal menginvalidasi token: %w", err)
	}

	return nil
}