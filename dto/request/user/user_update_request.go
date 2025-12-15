package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

//
// ========== 1. USER UPDATE DIRI SENDIRI ==========
//
type UserUpdateRequest struct {
	Nama string `json:"nama"`
}

func (r *UserUpdateRequest) BindAndValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("payload tidak valid")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Nama,
			validation.Length(3, 30).Error("nama harus 3-30 karakter"),
		),
	)
}


//
// ========== 2. ADMIN UPDATE DIRI SENDIRI ==========
// admin hanya boleh ubah nama, nama lengkap, dan NRA
//
type AdminSelfUpdateRequest struct {
	Nama        string `json:"nama,omitempty"`
	NamaLengkap string `json:"nama_lengkap,omitempty"`
	NRA         string `json:"nra,omitempty"`
}

func (r *AdminSelfUpdateRequest) BindAndValidate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        return fmt.Errorf("payload tidak valid")
    }

    return validation.ValidateStruct(r,
        validation.Field(&r.Nama,
            validation.Length(3, 30).Error("nama harus 3-30 karakter"),
        ),
        validation.Field(&r.NamaLengkap,
            validation.Length(3, 50).Error("nama lengkap harus 3-50 karakter"),
        ),
        validation.Field(&r.NRA,
            validation.Length(2, 20).Error("nra harus 2-20 karakter"),
        ),
    )
}

//
// ========== 3. ADMIN UPDATE USER LAIN ==========
// admin bisa update user biasa, tapi tidak bisa ubah role
//
type AdminUpdateUserRequest struct {
	Nama           string `json:"nama,omitempty"`
	NamaLengkap    string `json:"nama_lengkap,omitempty"`
	NRA            string `json:"nra,omitempty"`
	AngkatanMapala string `json:"angkatan_mapala,omitempty"`
	AngkatanKampus string `json:"angkatan_kampus,omitempty"`
	NIM            string `json:"nim,omitempty"`
	Jurusan        string `json:"jurusan,omitempty"`
	Fakultas       string `json:"fakultas,omitempty"`
	NoHP           string `json:"no_hp,omitempty"`
	Status string `json:"status,omitempty"`
	FotoURL        string `json:"foto_url,omitempty"`
}

func (r *AdminUpdateUserRequest) BindandValidate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
        return fmt.Errorf("payload tidak valid")
    }

	return validation.ValidateStruct(r,
		validation.Field(&r.Nama,
			validation.Length(3, 30).Error("nama harus 3–30 karakter"),
		),
		validation.Field(&r.NamaLengkap,
			validation.Length(3, 50).Error("nama lengkap harus 3–50 karakter"),
		),
		validation.Field(&r.NRA,
			validation.Length(2, 20).Error("nra harus 2–20 karakter"),
		),
		validation.Field(&r.AngkatanMapala,
			validation.Length(0, 20).Error("angkatan mapala maksimal 20 karakter"),
		),
		validation.Field(&r.AngkatanKampus,
			validation.Min(22).Error("angkatan tertua adalah 22"),
		),
		validation.Field(&r.NIM,
			validation.Length(0, 20).Error("nim maksimal 20 karakter"),
		),
		validation.Field(&r.Jurusan,
			validation.Length(0, 50).Error("jurusan maksimal 50 karakter"),
		),
		validation.Field(&r.Fakultas,
			validation.Length(0, 50).Error("fakultas maksimal 50 karakter"),
		),
		validation.Field(&r.NoHP,
			validation.Length(10, 15).Error("nomor HP harus 10–15 digit"),
		),
		validation.Field(&r.FotoURL,
			validation.Length(0, 255).Error("foto_url maksimal 255 karakter"),
		),
		validation.Field(&r.Status,
			validation.In("aktif", "lulus", "keluar").Error("status harus salah satu dari 'aktif', 'lulus', atau 'keluar'"),
		),
	)
}

//
// ========== 4. SUPER ADMIN UPDATE DIRI SENDIRI ==========
// sama seperti admin: hanya nama, nama lengkap, NRA
//
type SuperAdminSelfUpdateRequest struct {
	Nama        string `json:"nama,omitempty"`
	NamaLengkap string `json:"nama_lengkap,omitempty"`
	NRA         string `json:"nra,omitempty"`
}

func (r SuperAdminSelfUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Nama,
			validation.Length(3, 30).Error("nama harus 3–30 karakter"),
		),
		validation.Field(&r.NamaLengkap,
			validation.Length(3, 50).Error("nama lengkap harus 3–50 karakter"),
		),
		validation.Field(&r.NRA,
			validation.Length(2, 20).Error("nra harus 2–20 karakter"),
		),
	)
}

//
// ========== 5. SUPER ADMIN UPDATE USER ATAU ADMIN ==========
// super admin boleh update password dan role
//
type SuperAdminUpdateUserRequest struct {
	Nama           string `json:"nama,omitempty"`
	NamaLengkap    string `json:"nama_lengkap,omitempty"`
	NRA            string `json:"nra,omitempty"`
	AngkatanMapala string `json:"angkatan_mapala,omitempty"`
	AngkatanKampus string `json:"angkatan_kampus,omitempty"`
	NIM            string `json:"nim,omitempty"`
	Jurusan        string `json:"jurusan,omitempty"`
	Fakultas       string `json:"fakultas,omitempty"`
	NoHP           string `json:"no_hp,omitempty"`
	FotoURL        string `json:"foto_url,omitempty"`
	Password       string `json:"password,omitempty"`
	Role           string `json:"role,omitempty"`
}

func (r SuperAdminUpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Nama,
			validation.Length(3, 30).Error("nama harus 3–30 karakter"),
		),
		validation.Field(&r.NamaLengkap,
			validation.Length(3, 50).Error("nama lengkap harus 3–50 karakter"),
		),
		validation.Field(&r.NRA,
			validation.Length(2, 20).Error("nra harus 2–20 karakter"),
		),
		validation.Field(&r.AngkatanMapala,
			validation.Length(0, 20).Error("angkatan mapala maksimal 20 karakter"),
		),
		validation.Field(&r.AngkatanKampus,
			validation.Min(22).Error("angkatan tertua adalah 22"),
		),
		validation.Field(&r.NIM,
			validation.Length(0, 20).Error("nim maksimal 20 karakter"),
		),
		validation.Field(&r.Jurusan,
			validation.Length(0, 50).Error("jurusan maksimal 50 karakter"),
		),
		validation.Field(&r.Fakultas,
			validation.Length(0, 50).Error("fakultas maksimal 50 karakter"),
		),
		validation.Field(&r.NoHP,
			validation.Length(10, 15).Error("nomor HP harus 10–15 digit"),
		),
		validation.Field(&r.FotoURL,
			validation.Length(0, 255).Error("foto_url maksimal 255 karakter"),
		),
		validation.Field(&r.Role,
			validation.Length(3, 20).Error("role harus 3–20 karakter"),
		),
	)
}
