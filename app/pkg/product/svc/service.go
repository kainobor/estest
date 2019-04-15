package svc

import (
	"context"

	"github.com/kainobor/estest/app/pkg/product"
)

type Service struct {
	repo product.Repository
}

func New(repo product.Repository) product.Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, query string, filter product.Filter, sort *product.Sort, page *product.Page) ([]product.Product, error) {
	return s.repo.List(ctx, query, filter, sort, page)
}
