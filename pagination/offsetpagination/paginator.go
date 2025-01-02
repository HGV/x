package offsetpagination

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	defaultPageSize int = 50
	maxPageSize     int = 100
)

type (
	Paged interface {
		Offset() int
		Limit() int
	}
	Paginator[T any] struct {
		page     int
		pageSize int
	}
	Result[T any] struct {
		Items    []T `json:"items"`
		NextPage int `json:"next_page,omitempty"`
	}
)

var _ Paged = new(Paginator[any])

func New[T any](page, pageSize int) Paginator[T] {
	p := Paginator[T]{
		page:     page,
		pageSize: pageSize,
	}

	if p.page <= 0 {
		p.page = 1
	}
	if p.pageSize <= 0 {
		p.pageSize = defaultPageSize
	}
	if p.pageSize > maxPageSize {
		p.pageSize = maxPageSize
	}

	return p
}

func (p *Paginator[T]) Offset() int {
	return (p.page - 1) * p.pageSize
}

func (p *Paginator[T]) Limit() int {
	return p.pageSize + 1
}

func (p *Paginator[T]) Paginate(items []T) Result[T] {
	if len(items) > p.pageSize {
		return Result[T]{
			Items:    items[:p.pageSize],
			NextPage: p.page + 1,
		}
	}
	return Result[T]{
		Items: items,
	}
}

func (r Result[T]) HasNextPage() bool {
	return r.NextPage > 0
}

var (
	ErrInvalidPage      = errors.New("query parameter `page` must be an integer, got string")
	ErrPageTooSmall     = errors.New("query parameter `page` must be non-negative")
	ErrInvalidPageSize  = errors.New("query parameter `page_size` must be an integer, got string")
	ErrPageSizeTooSmall = errors.New("query parameter `page_size` must be non-negative")
)

func Parse[T any](q url.Values) (*Paginator[T], error) {
	var page, pageSize int
	var err error

	if pageParam := q.Get("page"); pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			return nil, ErrInvalidPage
		}
		if page <= 0 {
			return nil, ErrPageTooSmall
		}
	}

	if pageSizeParam := q.Get("page_size"); pageSizeParam != "" {
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil {
			return nil, ErrInvalidPageSize
		}
		if pageSize <= 0 {
			return nil, ErrPageSizeTooSmall
		}
	}

	p := New[T](page, pageSize)
	return &p, nil
}
