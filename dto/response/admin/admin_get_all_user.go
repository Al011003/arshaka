package admin

type UserListResponse struct {
	ID             uint   `json:"id"`
	Nama           string `json:"nama"`
	AngkatanMapala string `json:"angkatan_mapala"`
	Status         string `json:"status"`
	Role           string `json:"role"`
}
