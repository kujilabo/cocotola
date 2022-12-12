//go:generate mockery --output mock --name StudentUsecaseStudy
package student

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/usecase"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type StudentUsecaseStudy interface {

	// study
	FindResults(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.StudyRecordWithProblemID, error)

	GetCompletionRate(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (map[string]int, error)

	// FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error)
	SetResult(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result, mastered bool) error
}

type studentUsecaseStudy struct {
	transaction service.Transaction
	pf          service.ProcessorFactory
}

func NewStudentUsecaseStudy(transaction service.Transaction, pf service.ProcessorFactory) StudentUsecaseStudy {
	return &studentUsecaseStudy{
		transaction: transaction,
		pf:          pf,
	}
}

func (s *studentUsecaseStudy) FindResults(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.StudyRecordWithProblemID, error) {
	var results []domain.StudyRecordWithProblemID
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		recordbook, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		tmpResults, err := recordbook.GetResultsSortedLevel(ctx)
		if err != nil {
			return liberrors.Errorf("failed to GetResultsSortedLevel. err: %w", err)
		}
		results = tmpResults
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *studentUsecaseStudy) GetCompletionRate(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (map[string]int, error) {
	var results map[string]int
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		recordbookSummary, err := student.FindRecordbookSummary(ctx, workbookID)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		tmpResults, err := recordbookSummary.GetCompletionRate(ctx)
		if err != nil {
			return liberrors.Errorf("failed to GetResultsSortedLevel. err: %w", err)
		}
		results = tmpResults
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *studentUsecaseStudy) SetResult(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result, mastered bool) error {
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		recordbook, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		if err := recordbook.SetResult(ctx, workbook.GetProblemType(), problemID, result, mastered); err != nil {
			return liberrors.Errorf("failed to SetResult. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *studentUsecaseStudy) findStudent(ctx context.Context, rf service.RepositoryFactory, userRf userS.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}

	return student, nil
}
