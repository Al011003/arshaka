package utils

import (
	"errors"
	"time"
)

const DateLayout = "2006-01-02"

func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("tanggal tidak boleh kosong")
	}

	date, err := time.Parse(DateLayout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("format tanggal harus YYYY-MM-DD")
	}

	return date.Truncate(24 * time.Hour), nil
}
