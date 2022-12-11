//go:generate mockery --output mock --name RepositoryFactory
//go:generate mockery --output mock --name Transaction
package service

import (
	"context"

	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type RepositoryFactory interface {
	NewWorkbookRepository(ctx context.Context) (WorkbookRepository, error)

	NewProblemRepository(ctx context.Context, problemType string) (ProblemRepository, error)

	NewProblemTypeRepository(ctx context.Context) (ProblemTypeRepository, error)

	NewStudyTypeRepository(ctx context.Context) (StudyTypeRepository, error)

	NewStudyRecordRepository(ctx context.Context) (StudyRecordRepository, error)

	NewRecordbookRepository(ctx context.Context) (RecordbookRepository, error)

	NewUserQuotaRepository(ctx context.Context) (UserQuotaRepository, error)

	NewStatRepository(ctx context.Context) (StatRepository, error)

	NewStudyStatRepository(ctx context.Context) (StudyStatRepository, error)
}

type Transaction interface {
	Do(ctx context.Context, fn func(rf RepositoryFactory, userRf userS.RepositoryFactory) error) error
}
