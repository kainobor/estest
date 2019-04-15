package product

import "context"

type Service interface {
	List(ctx context.Context, query string, filters Filter, sort *Sort, page *Page) ([]Product, error)
}

type Repository interface {
	List(ctx context.Context, query string, filters Filter, sort *Sort, page *Page) ([]Product, error)
}
