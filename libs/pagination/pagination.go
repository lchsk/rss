package pagination

import (
	"errors"
	"math"
)

type PaginationValues struct {
	Limit  int
	Offset int
}

// page - from 1...last page
// count - total number of items
// perPage - number of items per page
func GetPaginationValues(page int, count int, perPage int) (*PaginationValues, error) {
	if page < 1 {
		return nil, errors.New("page must be > 0")
	}

	if lastPage := GetLastPage(count, perPage); page > lastPage {
		return &PaginationValues{Limit: 0, Offset: 0}, nil
	}

	offset := (page - 1) * perPage

	values := &PaginationValues{Limit: perPage, Offset: offset}

	return values, nil
}

func GetPages(page int, totalPages int) (int, int) {
	var prev int = -1
	var next int = -1

	if page <= 0 || page > totalPages {
		return prev, next
	}

	if page < totalPages {
		next = page + 1
	}

	if page > 1 {
		prev = page - 1
	}

	return prev, next
}

func GetLastPage(count int, perPage int) int {
	return int(math.Ceil(float64(count) / float64(perPage)))
}
