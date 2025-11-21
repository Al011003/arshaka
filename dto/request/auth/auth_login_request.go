package auth

import validation "github.com/go-ozzo/ozzo-validation"

type LoginRequest struct {
	NRA      string `json:"nra"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {

	return validation.ValidateStruct(&r, 
		validation.Field(
			&r.NRA,
			validation.Required.Error("NRA wajib diisi"),
			validation.Length(2, 20).Error("nra harus 2-20 karakter"),
		),
		validation.Field(
			&r.Password,
			validation.Required.Error("password wajib diisi"),
		),
	)

}