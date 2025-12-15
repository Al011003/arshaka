package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func (r *UpdateEmailRequest) BindAndValidate(c *gin.Context) error {
	// Bind JSON dulu
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	// Validation Ozzo
	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required.Error("email wajib diisi"),
			is.Email.Error("format email tidak valid"),
		),
	)
}


