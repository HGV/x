package offsetpagination

import (
	"net/http"
	"strconv"
)

const (
	defaultPageSize int = 100
)

type (
	Paginator[T any] struct {
		page                  int
		pageSize, maxPageSize int
	}
	Result[T any] struct {
		Items    []T  `json:"items"`
		NextPage *int `json:"next_page,omitempty"`
	}
	Option[T any] func(*Paginator[T])
)

func New[T any](page, pageSize int, opts ...Option[T]) Paginator[T] {
	p := Paginator[T]{
		page:     page,
		pageSize: pageSize,
	}

	for _, opt := range opts {
		opt(&p)
	}

	if p.page <= 0 {
		p.page = 1
	}
	if p.pageSize <= 0 {
		p.pageSize = defaultPageSize
	}
	if p.maxPageSize > 0 && p.pageSize > p.maxPageSize {
		p.pageSize = p.maxPageSize
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
		nextPage := p.page + 1
		return Result[T]{
			Items:    items[:p.pageSize],
			NextPage: &nextPage,
		}
	}
	return Result[T]{
		Items: items,
	}
}

func Parse[T any](r *http.Request, opts ...Option[T]) Paginator[T] {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	return New[T](page, pageSize, opts...)
}

func WithMaxPageSize[T any](maxPageSize int) Option[T] {
	return func(opt *Paginator[T]) {
		opt.maxPageSize = maxPageSize
	}
}
