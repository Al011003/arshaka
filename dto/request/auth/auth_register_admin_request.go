package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterAdminRequest struct {
	Nama        string `json:"nama"`
	NamaLengkap string `json:"nama_lengkap"`
	NRA         string `json:"nra"`
	Role        string `json:"-"`
}

func (r *RegisterAdminRequest) BindAndValidate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        return fmt.Errorf("payload tidak valid")
    }

    // override role
    r.Role = "admin"

    return validation.ValidateStruct(r,
        validation.Field(&r.Nama,
            validation.Required.Error("nama wajib diisi"),
            validation.Length(3, 30).Error("nama harus 3-30 karakter"),
        ),
        validation.Field(&r.NamaLengkap,
            validation.Required.Error("nama lengkap wajib diisi"),
            validation.Length(3, 50).Error("nama lengkap harus 3-50 karakter"),
        ),
        validation.Field(&r.NRA,
            validation.Required.Error("NRA wajib diisi"),
            validation.Length(2, 20).Error("nra harus 2-20 karakter"),
        ),
    )
}
