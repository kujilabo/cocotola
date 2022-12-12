package service

import (
	"context"

	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type Transaction interface {
	Do(ctx context.Context, fn func(rf userS.RepositoryFactory) error) error
}
