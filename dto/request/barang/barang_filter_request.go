package barang

type BarangFilter struct {
	Keyword   string // search di nama, merk, kode
	Kategori  string // filter by kategori
	Tersedia  *bool  // true = stok_sisa > 0, false = stok_sisa = 0, nil = semua
	Page      int    // default: 1
	Limit     int    // default: 10
	SortBy    string // nama, kode, created_at, stok_sisa (default: created_at)
	SortOrder string // asc, desc (default: desc)
}