package pagination

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaginationValues(t *testing.T) {
	cases := []struct {
		values  *PaginationValues
		err     bool
		page    int
		count   int
		perPage int
	}{
		{&PaginationValues{Limit: 2, Offset: 0}, false, 1, 5, 2},
		{&PaginationValues{Limit: 2, Offset: 2}, false, 2, 5, 2},
		{&PaginationValues{Limit: 2, Offset: 4}, false, 3, 5, 2},
		{nil, true, 4, 5, 2},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			values, err := GetPaginationValues(tt.page, tt.count, tt.perPage)

			assert.Equal(t, tt.values, values)
			assert.Equal(t, tt.err, err != nil)
		})
	}
}
func TestGetPages(t *testing.T) {
	cases := []struct {
		expectedPrev int
		expectedNext int
		page         int
		totalPages   int
	}{
		{-1, 2, 1, 3},
		{1, 3, 2, 3},
		{2, -1, 3, 3},
		{-1, -1, 0, 3},
		{-1, -1, 4, 3},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			prev, next := GetPages(tt.page, tt.totalPages)
			assert.Equal(t, tt.expectedPrev, prev)
			assert.Equal(t, tt.expectedNext, next)
		})
	}
}

func TestGetLastPage(t *testing.T) {
	cases := []struct {
		expected int
		count    int
		perPage  int
	}{
		{3, 5, 2},
		{0, 0, 2},
		{5, 10, 2},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert.Equal(t, tt.expected, GetLastPage(tt.count, tt.perPage))
		})
	}
}
