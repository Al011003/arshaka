// dto/response/barang/barang_response.go
package barang

import "time"

type BarangAdminDetailResponse struct {
	ID        uint      `json:"id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	Merk      string    `json:"merk"`
	Deskripsi string    `json:"deskripsi"`
	Kategori  string    `json:"kategori"`
	StokTotal int       `json:"stok_total"`
	StokSisa  int       `json:"stok_sisa"`
	TahunBeli int       `json:"tahun_beli"`
	HargaBeli float64       `json:"harga_beli"`
	CoverURL  string    `json:"cover_url"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BarangListResponse - untuk list (bisa lebih ringkas kalau mau)
type BarangListResponse struct {
	ID        uint    `json:"id"`
	Kode      string  `json:"kode"`
	Nama      string  `json:"nama"`
	Merk      string  `json:"merk"`
	Kategori  string  `json:"kategori"`
	StokTotal int     `json:"stok_total"`
	StokSisa  int     `json:"stok_sisa"`
	CoverURL  string  `json:"cover_url"`
	Status string `json:"status"`
}

type BarangUserDetailResponse struct {
	ID        uint   `json:"id"`
	Kode      string `json:"kode"`
	Nama      string `json:"nama"`
	Merk      string `json:"merk"`
	Deskripsi string `json:"deskripsi"`
	Kategori  string `json:"kategori"`
	StokTotal int    `json:"stok_total"`
	StokSisa  int    `json:"stok_sisa"`
	CoverURL  string `json:"cover_url"`
	Status string `json:"status"`
}