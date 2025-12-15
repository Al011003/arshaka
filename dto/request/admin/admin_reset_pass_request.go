package request

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type ForgotPasswordRequest struct {
	NRA  string `json:"nra"`
	Name string `json:"name"`
}

func (r *ForgotPasswordRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return errors.New("payload tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.NRA,
			validation.Required.Error("NRA wajib diisi"),
			validation.Length(3, 50).Error("NRA harus 3–50 karakter"),
		),
		validation.Field(&r.Name,
			validation.Required.Error("nama wajib diisi"),
			validation.Length(3, 100).Error("nama harus 3–100 karakter"),
		),
	)
}
