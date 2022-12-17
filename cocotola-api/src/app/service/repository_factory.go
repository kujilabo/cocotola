//go:generate mockery --output mock --name RepositoryFactory
//go:generate mockery --output mock --name Transaction
package service

import (
	"context"

	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type RepositoryFactory interface {
	NewWorkbookRepository(ctx context.Context) WorkbookRepository

	NewProblemRepository(ctx context.Context, problemType string) (ProblemRepository, error)

	NewProblemTypeRepository(ctx context.Context) ProblemTypeRepository

	NewStudyTypeRepository(ctx context.Context) StudyTypeRepository

	NewStudyRecordRepository(ctx context.Context) StudyRecordRepository

	NewRecordbookRepository(ctx context.Context) RecordbookRepository

	NewUserQuotaRepository(ctx context.Context) UserQuotaRepository

	NewStatRepository(ctx context.Context) StatRepository

	NewStudyStatRepository(ctx context.Context) StudyStatRepository

	NewUserRepositoryFactory(ctx context.Context) (userS.RepositoryFactory, error)

	NewJobRepositoryFactory(ctx context.Context) (jobS.RepositoryFactory, error)
}

type Transaction interface {
	Do(ctx context.Context, fn func(rf RepositoryFactory) error) error
}
