//go:generate mockery --output mock --name RepositoryFactory
//go:generate mockery --output mock --name Transaction
package service

import (
	"context"
)

type RepositoryFactory interface {
	NewJobStatusRepository(ctx context.Context) (JobStatusRepository, error)
	NewJobHistoryRepository(ctx context.Context) (JobHistoryRepository, error)
}

type Transaction interface {
	Do(ctx context.Context, fn func(rf RepositoryFactory) error) error
}
