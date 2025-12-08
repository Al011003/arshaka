package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type LoginRequest struct {
	NRA      string `json:"nra"`
	Password string `json:"password"`
}

func (r *LoginRequest) BindAndValidate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        return fmt.Errorf("payload tidak valid")
    }

    return validation.ValidateStruct(r,
        validation.Field(
            &r.NRA,
            validation.Required.Error("NRA wajib diisi"),
            validation.Length(2, 20).Error("NRA harus 2–20 karakter"),
        ),
        validation.Field(
            &r.Password,
            validation.Required.Error("password wajib diisi"),
        ),
    )
}
