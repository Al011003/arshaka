// errors/errors.go
package errors

import "errors"

var (
    ErrNotFound      = errors.New("data tidak ditemukan")
    ErrUnauthorized  = errors.New("tidak memiliki akses")
    ErrBadRequest    = errors.New("request tidak valid")
    ErrInternalServer = errors.New("terjadi kesalahan server")
)

// IsNotFound checks if error is not found
func IsNotFound(err error) bool {
    return errors.Is(err, ErrNotFound)
}

func IsUnauthorized(err error) bool {
    return errors.Is(err, ErrUnauthorized)
}