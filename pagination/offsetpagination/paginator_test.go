package offsetpagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
			p: New[any](3, 100),
			expected: Paginator[any]{
				page:     3,
				pageSize: 100,
			},
		},
		{
			p: New[any](3, 300),
			expected: Paginator[any]{
				page:     3,
				pageSize: 100,
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
			p:              New[any](1, 100),
			expectedOffset: 0,
			expectedLimit:  101,
		},
		{
			p:              New[any](3, 100),
			expectedOffset: 200,
			expectedLimit:  101,
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
	p := New[int](1, 75)

	items := make([]int, 1000)
	for i := 0; i < len(items); i++ {
		items[i] = i
	}

	t.Run("should have a next page", func(t *testing.T) {
		result := p.Paginate(items)
		assert.Equal(t, len(result.Items), 75)
		assert.Equal(t, result.NextPage, 2)
		assert.True(t, result.HasNextPage())
	})

	t.Run("should not have a next page", func(t *testing.T) {
		result := p.Paginate(items[950:])
		assert.Equal(t, len(result.Items), 50)
		assert.Zero(t, result.NextPage)
		assert.False(t, result.HasNextPage())
	})
}

func TestParse(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/items?page=3&page_size=100", nil)
	p, _ := Parse[any](r.URL.Query())
	assert.Equal(t, p.page, 3)
	assert.Equal(t, p.pageSize, 100)
}
