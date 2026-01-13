package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type RequestOTPRequest struct {
	Email string `json:"email" example:"user@mail.com"`
}

func (r *RequestOTPRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required.Error("email wajib diisi"),
			is.Email.Error("format email tidak valid"),
		),
	)
}

// ==============================

type VerifyResetOTPRequest struct {
	Email string `json:"email" example:"user@mail.com"`
	OTP   string `json:"otp" example:"123456"`
}

func (r *VerifyResetOTPRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required.Error("email wajib diisi"),
			is.Email.Error("format email tidak valid"),
		),
		validation.Field(&r.OTP,
			validation.Required.Error("otp wajib diisi"),
			validation.Length(6, 6).Error("otp harus 6 digit"),
		),
	)
}

// ==============================

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token" example:"abc123"`
	NewPassword string `json:"new_password" example:"password123"`
}

func (r *ResetPasswordRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("format request tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.ResetToken,
			validation.Required.Error("reset token wajib diisi"),
		),
		validation.Field(&r.NewPassword,
			validation.Required.Error("password baru wajib diisi"),
			validation.Length(6, 0).Error("password minimal 6 karakter"),
		),
	)
}
