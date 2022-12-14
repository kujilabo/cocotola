//go:generate mockery --output mock --name RecordbookRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
)

var ErrStudyResultNotFound = errors.New("StudyResult not found")

type CountAnsweredResult struct {
	WorkbookID    uint
	ProblemTypeID uint
	StudyTypeID   uint
	Answered      int
	Mastered      int
}
type CountAnsweredResults struct {
	Results []CountAnsweredResult
}

type RecordbookRepository interface {
	FindStudyRecords(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType domain.StudyTypeName) (map[domain.ProblemID]domain.StudyRecord, error)

	SetResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType domain.StudyTypeName, problemType domain.ProblemTypeName, problemID domain.ProblemID, studyResult, mastered bool) error

	CountMasteredProblems(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (map[domain.StudyTypeName]int, error)
}
