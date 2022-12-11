//go:generate mockery --output mock --name JobHistoryRepository
package service

import "context"

type JobHistory interface {
}

type JobHistoryRepository interface {
	AddJobHistory(ctx context.Context, jobHistory JobHistory) error
}
