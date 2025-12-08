package superadmin

type SuperAdminProfileResponse struct {
	NRA         string `json:"nra"`
	Nama        string `json:"nama"`
	NamaLengkap string `json:"nama_lengkap"`
	Status      string `json:"status"`
}