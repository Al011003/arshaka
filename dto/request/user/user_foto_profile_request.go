package user

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserPhotoRequest struct {
	File *multipart.FileHeader `json:"-" form:"file"`
}

func (r *UserPhotoRequest) BindAndValidate(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errors.New("file wajib diupload")
	}

	r.File = file

	// Validasi ukuran maksimal 2MB
	if file.Size > 2*1024*1024 {
		return errors.New("ukuran foto maksimal 2MB")
	}

	// Validasi extensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return errors.New("format foto harus jpg, jpeg, atau png")
	}

	return nil
}
