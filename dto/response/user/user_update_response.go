package user

import "time"

type UpdateUserResponse struct {
	ID     uint   `json:"id"`
	NRA    string `json:"nra"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type UserDetailResponse struct {
	ID             uint      `json:"id"`
	FotoURL 		string `json:"foto_url"`
	Nama           string    `json:"nama"`
	NamaLengkap    string    `json:"nama_lengkap"`
	NRA            string    `json:"nra"`
	Jurusan        string    `json:"jurusan"`
	Fakultas       string    `json:"fakultas"`
	AngkatanMapala string    `json:"angkatan_mapala"`
	AngkatanKampus string    `json:"angkatan_kampus"`
	NIM            string    `json:"nim"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}