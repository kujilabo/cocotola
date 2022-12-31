//go:generate mockery --output mock --name ProblemRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
)

var ErrProblemAlreadyExists = errors.New("problem already exists")
var ErrProblemNotFound = errors.New("problem not found")
var ErrProblemOtherError = errors.New("problem other error")

type ProblemRepository interface {
	// FindProblems searches for problems based on search condition
	FindProblems(ctx context.Context, operator domain.StudentModel, param domain.ProblemSearchCondition) (domain.ProblemSearchResult, error)

	FindAllProblems(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (domain.ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, operator domain.StudentModel, param domain.ProblemIDsCondition) (domain.ProblemSearchResult, error)

	FindProblemsByCustomCondition(ctx context.Context, operator domain.StudentModel, condition interface{}) ([]domain.ProblemModel, error)

	FindProblemByID(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter1) (Problem, error)

	FindProblemIDs(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) ([]domain.ProblemID, error)

	// AddProblem register a new problem
	AddProblem(ctx context.Context, operator domain.StudentModel, param domain.ProblemAddParameter) (domain.ProblemID, error)

	UpdateProblem(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error

	UpdateProblemProperty(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, operator domain.StudentModel, id domain.ProblemSelectParameter2) error

	CountProblems(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (int, error)
}
