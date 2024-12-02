package offsetpagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HGV/x/pointerx"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		p        Paginator[any]
		expected Paginator[any]
	}{
		{
			p: New[any](0, 0),
			expected: Paginator[any]{
				page:     1,
				pageSize: defaultPageSize,
			},
		},
		{
			p: New[any](-1, -100),
			expected: Paginator[any]{
				page:     1,
				pageSize: defaultPageSize,
			},
		},
		{
			p: New[any](3, 250),
			expected: Paginator[any]{
				page:     3,
				pageSize: 250,
			},
		},
		{
			p: New[any](3, 250, WithMaxPageSize[any](100)),
			expected: Paginator[any]{
				page:        3,
				pageSize:    100,
				maxPageSize: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.p)
		})
	}
}

func TestOffsetAndLimit(t *testing.T) {
	tests := []struct {
		p              Paginator[any]
		expectedOffset int
		expectedLimit  int
	}{
		{
			p:              New[any](1, 250),
			expectedOffset: 0,
			expectedLimit:  251,
		},
		{
			p:              New[any](3, 250),
			expectedOffset: 500,
			expectedLimit:  251,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expectedOffset, tt.p.Offset())
			assert.Equal(t, tt.expectedLimit, tt.p.Limit())
		})
	}
}

func TestPaginate(t *testing.T) {
	p := New[int](1, 250)

	items := make([]int, 1000)
	for i := 0; i < len(items); i++ {
		items[i] = i
	}

	t.Run("should have a next page", func(t *testing.T) {
		result := p.Paginate(items)
		assert.Equal(t, len(result.Items), 250)
		assert.Equal(t, result.NextPage, pointerx.Ptr(2))
	})

	t.Run("should not have a next page", func(t *testing.T) {
		result := p.Paginate(items[790:])
		assert.Equal(t, len(result.Items), 210)
		assert.Nil(t, result.NextPage)
	})
}

func TestParse(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/items?page=3&page_size=250", nil)
	p := Parse[any](r)
	assert.Equal(t, p.page, 3)
	assert.Equal(t, p.pageSize, 250)
}
