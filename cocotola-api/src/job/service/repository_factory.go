//go:generate mockery --output mock --name RepositoryFactory
package service

import (
	"context"
)

type RepositoryFactory interface {
	NewJobStatusRepository(ctx context.Context) JobStatusRepository
	NewJobHistoryRepository(ctx context.Context) JobHistoryRepository
}
