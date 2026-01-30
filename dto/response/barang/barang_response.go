// dto/response/barang/barang_response.go
package barang

import "time"

// dto/response/barang/barang.go

type BarangStatsResponse struct {
    TotalUnit    int64 `json:"total_unit"`
    Siap         int64 `json:"siap"`
    Dipinjam     int64 `json:"dipinjam"`
    SiapTerbatas int64 `json:"siap_terbatas"`
    Maintenance  int64 `json:"maintenance"`
}

type BarangAdminDetailResponse struct {
    ID        uint      `json:"id"`
    Kode      string    `json:"kode"`
    Nama      string    `json:"nama"`
    Merk      string    `json:"merk"`
    Deskripsi string    `json:"deskripsi"`
    Kategori  string    `json:"kategori"`
    Status    string    `json:"status"`
    CoverURL  string    `json:"cover_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Harga float64 `json:"haraga"`
    
    Stats *BarangStatsResponse `json:"stats"`

    AllUnitsInactive bool   `json:"all_units_inactive"`
	StatusWarning    string `json:"status_warning,omitempty"`
}

type BarangUserDetailResponse struct {
	ID            uint    `json:"id"`
	Kode          string  `json:"kode"`           // ← TAMBAH
	Nama          string  `json:"nama"`
	Merk          string  `json:"merk"`           // ← TAMBAH
	Deskripsi     string  `json:"deskripsi"`      // ← TAMBAH
	Kategori      string  `json:"kategori"`       // ← TAMBAH
	CoverURL      string  `json:"cover_url"`      // ← TAMBAH
	Harga         float64 `json:"harga"`          // ← TAMBAH (optional, tergantung mau ditampilin ke user atau engga)
	TotalUnit     int64   `json:"total_unit"`
	AvailableUnit int64   `json:"available_unit"`
}

type BarangListResponse struct {
    ID            uint   `json:"id"`
    Kode          string `json:"kode"`
    Nama          string `json:"nama"`
    Merk          string `json:"merk"`
    Kategori      string `json:"kategori"`
    Status        string `json:"status"`
    CoverURL      string `json:"cover_url"`
    Harga float64 `json:"harga"`
    TotalUnit     int64  `json:"total_unit"`
    AvailableUnit int64  `json:"available_unit"`
}


type UnitListResponse struct {
	ID        uint   `json:"id"`
	KodeUnit string `json:"kode_unit"`
	Status   string `json:"status"`
	Kondisi  string `json:"kondisi"`
}

type KomponenResponse struct {
	ID           uint   `json:"id"`
	Nama         string `json:"nama"`
	JumlahWajib  int    `json:"jumlah_wajib"`
	JumlahAktual int    `json:"jumlah_aktual"`
	Status       string `json:"status"`
    HargaSatuan      int64  `json:"harga_satuan"`
}

type UnitCreateResponse struct {
	Message  string   `json:"message"`
	KodeUnit []string `json:"kode_unit"`
}

// dto/response/barang/unit.go

type UnitDetailResponse struct {
	ID              uint   `json:"id"`
	BarangKode      string `json:"barang_kode"`
	BarangNama      string `json:"barang_nama"`
	KodeUnit        string `json:"kode_unit"`
	Status          string `json:"status"`
	Kondisi         string `json:"kondisi"`
	FixPriority     int    `json:"fix_priority"`
	TahunPerolehan  int    `json:"tahun_perolehan"`
	LokasiPenyimpan string `json:"lokasi_penyimpan"`
	Catatan         string `json:"catatan"`
	// Bisa tambahin Komponen []KomponenResponse nanti
}

type BarangStatusSuggestion struct {
	CurrentStatus   string `json:"current_status"`
	SuggestedStatus string `json:"suggested_status"`
	Reason          string `json:"reason"`
	ShouldChange    bool   `json:"should_change"`
}