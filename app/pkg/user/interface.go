package user

import "context"

type Service interface {
	Auth(ctx context.Context, login, password string) (string, error)
}

type Repository interface {
	GetID(ctx context.Context, login, pass string) (string, error)
}
