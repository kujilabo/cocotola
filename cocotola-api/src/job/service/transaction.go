//go:generate mockery --output mock --name Transaction
package service

import (
	"context"
)

type Transaction interface {
	Do(ctx context.Context, fn func(rf RepositoryFactory) error) error
}
