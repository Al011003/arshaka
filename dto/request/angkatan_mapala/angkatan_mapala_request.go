package angkatan

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type AngkatanMapalaRequest struct {
	Nama string `json:"nama" example:"Mapala 2023"`
}

func (r *AngkatanMapalaRequest)  BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(
			&r.Nama,
			validation.Required.Error("nama wajib diisi"),
			validation.Length(2, 100).Error("nama harus 2-100 karakter"),
		),
	)
}
