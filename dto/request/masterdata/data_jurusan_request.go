package data

type JurusanRequest struct {
    Nama       string `json:"nama" binding:"required"`
    FakultasID uint   `json:"fakultas_id" binding:"required"`
}