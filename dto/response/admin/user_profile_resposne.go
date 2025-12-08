package admin

type AdminDetailResponse struct {
	ID          uint   `json:"id"`
	Nama        string `json:"nama"`
	NamaLengkap string `json:"nama_lengkap"`
	NRA         string `json:"nra"`
	Role        string `json:"role"`
}
