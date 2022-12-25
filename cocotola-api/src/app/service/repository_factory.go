//go:generate mockery --output mock --name RepositoryFactory
package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type RepositoryFactory interface {
	NewWorkbookRepository(ctx context.Context) WorkbookRepository

	NewProblemRepository(ctx context.Context, problemType domain.ProblemTypeName) (ProblemRepository, error)

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
