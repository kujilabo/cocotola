//go:generate mockery --output mock --name StudentUsecaseStudy
package student

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type StudentUsecaseStudy interface {

	// study
	FindResults(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType domain.StudyTypeName) ([]domain.StudyRecordWithProblemID, error)

	GetCompletionRate(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (map[domain.StudyTypeName]int, error)

	// FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error)
	SetResult(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType domain.StudyTypeName, problemID domain.ProblemID, result, mastered bool) error
}

type studentUsecaseStudy struct {
	transaction     service.Transaction
	pf              service.ProcessorFactory
	studyMonitor    service.StudyMonitor
	findStudentFunc FindStudentFunc
}

func NewStudentUsecaseStudy(transaction service.Transaction, pf service.ProcessorFactory, studyMonitor service.StudyMonitor, findStudentFunc FindStudentFunc) StudentUsecaseStudy {
	return &studentUsecaseStudy{
		transaction:     transaction,
		pf:              pf,
		studyMonitor:    studyMonitor,
		findStudentFunc: findStudentFunc,
	}
}

func (s *studentUsecaseStudy) FindResults(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType domain.StudyTypeName) ([]domain.StudyRecordWithProblemID, error) {
	var results []domain.StudyRecordWithProblemID
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, organizationID, operatorID)
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
		return nil, liberrors.Errorf("FindResults. err: %w", err)
	}

	return results, nil
}

func (s *studentUsecaseStudy) GetCompletionRate(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (map[domain.StudyTypeName]int, error) {
	var results map[domain.StudyTypeName]int
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, organizationID, operatorID)
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
		return nil, liberrors.Errorf("GetCompletionRate. err: %w", err)
	}

	return results, nil
}

func (s *studentUsecaseStudy) SetResult(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType domain.StudyTypeName, problemID domain.ProblemID, result, mastered bool) error {
	logger := log.FromContext(ctx)
	var problemType domain.ProblemTypeName
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return liberrors.Errorf("student.FindWorkbookByID. err: %w", err)
		}

		problemType = workbook.GetProblemType()
		recordbook, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		if err := recordbook.SetResult(ctx, workbook.GetProblemType(), problemID, result, mastered); err != nil {
			return liberrors.Errorf("failed to SetResult. err: %w", err)
		}
		return nil
	}); err != nil {
		return liberrors.Errorf("SetResult. err: %w", err)
	}

	studyEvent := service.NewStudyEvent(organizationID, operatorID, service.StudyEventTypeAnswer, problemType, studyType, problemID)
	if err := s.studyMonitor.NotifyObservers(ctx, studyEvent); err != nil {
		logger.Errorf("NotifyObservers. err: %v", err)
	}

	return nil
}

func (s *studentUsecaseStudy) findStudent(ctx context.Context, rf service.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	student, err := s.findStudentFunc(ctx, rf, organizationID, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}

	return student, nil
}
