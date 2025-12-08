package utils

import "math"

func ValidatePagination(page int, limit int) (int, int) {
    if page < 1 {
        page = 1
    }

    if limit < 10 { 
        limit = 10
    }

    // limit kelipatan 5
    if limit%5 != 0 {
        limit = 10 // fallback aman
    }

    return page, limit
}

func CalcTotalPages(totalRows int, limit int) int {
    return int(math.Ceil(float64(totalRows) / float64(limit)))
}
