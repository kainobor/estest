package svc

import (
	"context"

	"github.com/kainobor/estest/app/pkg/user"
)

type Service struct {
	repo user.Repository
}

func New(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Auth(ctx context.Context, login, pass string) (string, error) {
	return s.repo.GetID(ctx, login, pass)
}
