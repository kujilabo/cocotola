//go:generate mockery --output mock --name ProblemQuotaProcessor
package service

import (
	"context"
	"io"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	// pluginCommon "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/domain"
)

// type Added int
// type Updated int

type ProblemAddProcessor interface {
	AddProblem(ctx context.Context, repo RepositoryFactory, operator domain.StudentModel, workbookModel domain.WorkbookModel, param ProblemAddParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)
}

type ProblemUpdateProcessor interface {
	UpdateProblem(ctx context.Context, repo RepositoryFactory, operator domain.StudentModel, workbookModel domain.WorkbookModel, id ProblemSelectParameter2, param ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)

	UpdateProblemProperty(ctx context.Context, repo RepositoryFactory, operator domain.StudentModel, workbookModel domain.WorkbookModel, id ProblemSelectParameter2, param ProblemUpdateParameter) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)
}

type ProblemRemoveProcessor interface {
	RemoveProblem(ctx context.Context, repo RepositoryFactory, operator domain.StudentModel, id ProblemSelectParameter2) ([]domain.ProblemID, []domain.ProblemID, []domain.ProblemID, error)
}

type ProblemImportProcessor interface {
	CreateCSVReader(ctx context.Context, workbookID domain.WorkbookID, reader io.Reader) (ProblemAddParameterIterator, error)
}

type ProblemQuotaProcessor interface {
	// IsExceeded(ctx context.Context, repo RepositoryFactory, operator Student, name string) (bool, error)

	// Increment(ctx context.Context, repo RepositoryFactory, operator Student, name string) (bool, error)

	// Decrement(ctx context.Context, repo RepositoryFactory, operator Student, name string) (bool, error)

	GetUnitForSizeQuota() QuotaUnit

	GetLimitForSizeQuota() int

	GetUnitForUpdateQuota() QuotaUnit

	GetLimitForUpdateQuota() int
}
