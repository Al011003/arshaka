package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdatePasswordRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewConfirmPassword string `json:"new_confirm_password"`
}

func (r *UpdatePasswordRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	// Validasi ozzo
	return validation.ValidateStruct(r,
		validation.Field(&r.OldPassword,
			validation.Required.Error("password lama wajib diisi"),
		),
		validation.Field(&r.NewPassword,
			validation.Required.Error("password baru wajib diisi"),
			validation.Length(6, 50).Error("password baru minimal 6 karakter"),
		),
		validation.Field(&r.NewConfirmPassword,
			validation.Required.Error("konfirmasi password wajib diisi"),
			validation.By(func(value interface{}) error {
				if value.(string) != r.NewPassword {
					return fmt.Errorf("konfirmasi password tidak cocok")
				}
				return nil
			}),
		),
	)
}
