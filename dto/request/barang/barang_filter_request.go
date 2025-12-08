package barang

type BarangFilter struct {
	Keyword   string `form:"keyword"`  // ?keyword=anjay
	Kategori  string `form:"kategori"` // ?kategori=Elektronik
	Status    string `form:"status"`
	Page      int    `form:"page"`       // ?page=2
	Limit     int    `form:"limit"`      // ?limit=20
	SortBy    string `form:"sort_by"`    // ?sort_by=nama
	SortOrder string `form:"sort_order"` // ?sort_order=asc
}
