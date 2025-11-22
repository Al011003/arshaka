package auth

import validation "github.com/go-ozzo/ozzo-validation"

type RegisterRequest struct {
	Nama           string `json:"nama"`
	NamaLengkap    string `json:"nama_lengkap"`
	NRA            string `json:"nra"`
	AngkatanMapala string `json:"angkatan_mapala"`
	AngkatanKampus string `json:"angkatan_kampus"`
	NIM            string `json:"nim"`
	Jurusan        string `json:"jurusan"`
	Fakultas       string `json:"fakultas"`
	NoHP           string `json:"no_hp"`
	Role string `json:"role"`
}

func (r RegisterRequest) Validate() error {

	// VALIDASI YANG WAJIB UNTUK SEMUA ROLE
	baseRules := []*validation.FieldRules{
		validation.Field(
			&r.Nama,
			validation.Required.Error("nama wajib diisi"),
			validation.Length(3, 30).Error("nama harus 3-30 karakter"),
		),
		validation.Field(
			&r.NamaLengkap,
			validation.Required.Error("nama lengkap wajib diisi"),
			validation.Length(3, 50).Error("nama lengkap harus 3-50 karakter"),
		),
		validation.Field(
			&r.NRA,
			validation.Required.Error("NRA wajib diisi"),
			validation.Length(2, 20).Error("nra harus 2-20 karakter"),
		),
		validation.Field(
			&r.Role,
			validation.Required.Error("role wajib diisi"),
		),
	}

	// TAMBAH VALIDASI KHUSUS UNTUK ROLE ANGGOTA
	if r.Role == "anggota" {
		baseRules = append(baseRules,
			validation.Field(
				&r.AngkatanMapala,
				validation.Required.Error("angkatan mapala wajib diisi untuk anggota"),
			),
			validation.Field(
				&r.AngkatanKampus,
				validation.Required.Error("angkatan kampus wajib diisi untuk anggota"),validation.Min(22).Error("angkatan tertua adalah 22"),
			),
			validation.Field(
				&r.NIM,
				validation.Required.Error("NIM wajib diisi untuk anggota"),
			),
			validation.Field(
				&r.Jurusan,
				validation.Required.Error("jurusan wajib diisi untuk anggota"),
			),
			validation.Field(
				&r.Fakultas,
				validation.Required.Error("fakultas wajib diisi untuk anggota"),
			),
			validation.Field(
				&r.NoHP,
				validation.Required.Error("nomor HP wajib diisi untuk anggota"),
				validation.Length(10, 15).Error("nomor HP harus 10–15 digit"),
			),
		)
	}

	// ROLE ADMIN atau SUPERADMIN tidak butuh field anggota
	// jadi cukup baseRules saja

	return validation.ValidateStruct(&r, baseRules...)
}
