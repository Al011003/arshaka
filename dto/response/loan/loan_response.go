package loan



type BarangInLoanResponse struct {
	ID       uint   `json:"id"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	Merk     string `json:"merk"`
	Kategori string `json:"kategori"`
	CoverURL string `json:"cover_url"`
}
