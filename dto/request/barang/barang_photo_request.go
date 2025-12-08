// dto/request/barang/photo.go
package barang

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type BarangPhotoRequest struct {
	File *multipart.FileHeader `json:"-" form:"file"`
}

func (r *BarangPhotoRequest) BindAndValidate(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errors.New("file wajib diupload")
	}

	r.File = file

	// Validasi ukuran maksimal 5MB
	if file.Size > 5*1024*1024 {
		return errors.New("ukuran foto maksimal 5MB")
	}

	// Validasi extensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return errors.New("format foto harus jpg, jpeg, png, atau webp")
	}

	return nil
}