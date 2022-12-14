//go:generate mockery --output mock --name StudentUsecaseProblem
package student

import (
	"context"
	"errors"
	"io"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type StudentUsecaseProblem interface {
	// problem
	FindProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param domain.ProblemSearchCondition) (domain.ProblemSearchResult, error)

	FindAllProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (domain.ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param domain.ProblemIDsCondition) (domain.ProblemSearchResult, error)

	FindProblemByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter1) (domain.ProblemModel, error)

	FindProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error)

	AddProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, param domain.ProblemAddParameter) ([]domain.ProblemID, error)

	UpdateProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error

	UpdateProblemProperty(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter2) error

	ImportProblems(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, newIterator func(workbookID domain.WorkbookID, problemType domain.ProblemTypeName) (service.ProblemAddParameterIterator, error)) error
}

type studentUsecaseProblem struct {
	transaction     service.Transaction
	pf              service.ProcessorFactory
	problemMonitor  service.ProblemMonitor
	findStudentFunc FindStudentFunc
}

func NewStudentUsecaseProblem(transaction service.Transaction, pf service.ProcessorFactory, problemMonitor service.ProblemMonitor, findStudentFunc FindStudentFunc) StudentUsecaseProblem {
	return &studentUsecaseProblem{
		transaction:     transaction,
		pf:              pf,
		problemMonitor:  problemMonitor,
		findStudentFunc: findStudentFunc,
	}
}

func (s *studentUsecaseProblem) FindProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param domain.ProblemSearchCondition) (domain.ProblemSearchResult, error) {
	var result domain.ProblemSearchResult
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblems(ctx, student, param)
		if err != nil {
			return liberrors.Errorf("workbook.FindProblems. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("FindProblemsByWorkbookID. err: %w", err)
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindAllProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (domain.ProblemSearchResult, error) {
	var result domain.ProblemSearchResult
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindAllProblems(ctx, student)
		if err != nil {
			return liberrors.Errorf("workbook.FindAllProblems. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("FindAllProblemsByWorkbookID. err: %w", err)
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindProblemsByProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param domain.ProblemIDsCondition) (domain.ProblemSearchResult, error) {
	var result domain.ProblemSearchResult
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblemsByProblemIDs(ctx, student, param)
		if err != nil {
			return liberrors.Errorf("workbook.FindProblemsByProblemIDs. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("FindProblemsByProblemIDs. err: %w", err)
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindProblemByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter1) (domain.ProblemModel, error) {
	var result domain.ProblemModel
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblemByID(ctx, student, id.GetProblemID())
		if err != nil {
			return liberrors.Errorf("workbook.FindProblemByID. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("FindProblemByID. err: %w", err)
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error) {
	var result []domain.ProblemID
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblemIDs(ctx, student)
		if err != nil {
			return liberrors.Errorf("workbook.FindProblemIDs. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("FindProblemIDs. err: %w", err)
	}
	return result, nil
}

func (s *studentUsecaseProblem) AddProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, param domain.ProblemAddParameter) ([]domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	var result []domain.ProblemID
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		studentService, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, param.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := s.addProblem(ctx, studentService, workbook, param)
		if err != nil {
			return liberrors.Errorf("s.addProblem. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("AddProblem. err: %w", err)
	}
	logger.Debug("problem id: %d", result)
	return result, nil
}

func (s *studentUsecaseProblem) UpdateProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error {
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		if err := s.updateProblem(ctx, student, workbook, id, param); err != nil {
			return liberrors.Errorf("s.updateProblem. err: %w", err)
		}
		return nil
	}); err != nil {
		return liberrors.Errorf("UpdateProblem. err: %w", err)
	}
	return nil
}

func (s *studentUsecaseProblem) UpdateProblemProperty(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error {
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		if err := s.updateProblemProperty(ctx, student, workbook, id, param); err != nil {
			return liberrors.Errorf("failed to UpdateProblem. err: %w", err)
		}
		return nil
	}); err != nil {
		return liberrors.Errorf("UpdateProblemProperty. err: %w", err)
	}
	return nil
}

func (s *studentUsecaseProblem) RemoveProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id domain.ProblemSelectParameter2) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.RemoveProblem")

	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		_, _, removedIDs, err := workbook.RemoveProblem(ctx, student, id)
		if err != nil {
			return liberrors.Errorf("workbook.RemoveProblem. err: %w", err)
		}
		problemType := workbook.GetProblemType()

		{
			event := service.NewProblemEvent(student.GetOrganizationID(), student.GetAppUserID(), service.ProblemEventTypeRemove, problemType, removedIDs)
			if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
				return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
			}
		}
		// if err := student.DecrementQuotaUsage(ctx, problemType, "Size", 1); err != nil {
		// 	return liberrors.Errorf("student.DecrementQuotaUsage. err: %w", err)
		// }

		return nil
	}); err != nil {
		return liberrors.Errorf("RemoveProblem. err: %w", err)
	}
	return nil
}

func (s *studentUsecaseProblem) ImportProblems(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, newIterator func(workbookID domain.WorkbookID, problemType domain.ProblemTypeName) (service.ProblemAddParameterIterator, error)) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.ImportProblems")

	var problemType domain.ProblemTypeName
	{
		if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
			_, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, workbookID)
			if err != nil {
				return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
			}
			problemType = workbook.GetProblemType()
			return nil
		}); err != nil {
			return liberrors.Errorf("GetProblemType. err: %w", err)
		}
	}
	iterator, err := newIterator(workbookID, problemType)
	if err != nil {
		return liberrors.Errorf("newIterator. err: %w", err)
	}

	for {
		param, err := iterator.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return liberrors.Errorf("iterator.Next. err: %w", err)
		}
		if param == nil {
			continue
		}

		// FIXME
		logger.Infof("param.properties: %+v", param.GetProperties())

		if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
			student, workbook, err := s.findStudentAndWorkbook(ctx, rf, organizationID, operatorID, workbookID)
			if err != nil {
				return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
			}

			id, err := s.addProblem(ctx, student, workbook, param)
			if errors.Is(err, service.ErrProblemAlreadyExists) {
				logger.Infof("Problem already exists. param: %+v", param)
				return nil
			}

			if err != nil {
				return liberrors.Errorf("s.addProblem. err: %w", err)
			}
			logger.Infof("%d", id)

			return nil
		}); err != nil {
			return liberrors.Errorf("addProblem. err: %w", err)
		}
	}
	return nil
}

func (s *studentUsecaseProblem) findStudentAndWorkbook(ctx context.Context, rf service.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (service.Student, service.Workbook, error) {
	student, err := s.findStudentFunc(ctx, rf, organizationID, operatorID)
	if err != nil {
		return nil, nil, liberrors.Errorf("FindStudent. err: %w", err)
	}
	workbook, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, liberrors.Errorf("FindWorkbookByID. err: %w", err)
	}
	return student, workbook, nil
}

func (s *studentUsecaseProblem) addProblem(ctx context.Context, student service.Student, workbook service.Workbook, param domain.ProblemAddParameter) ([]domain.ProblemID, error) {
	problemType := workbook.GetProblemType()
	if err := student.CheckQuota(ctx, problemType, "Size"); err != nil {
		return nil, liberrors.Errorf("student.CheckQuota. err: %w", err)
	}
	if err := student.CheckQuota(ctx, problemType, "Update"); err != nil {
		return nil, liberrors.Errorf("student.CheckQuota. err: %w", err)
	}
	addedIDs, _, _, err := workbook.AddProblem(ctx, student, param)
	if err != nil {
		return nil, liberrors.Errorf("workbook.AddProblem. err: %w", err)
	}

	event := service.NewProblemEvent(student.GetOrganizationID(), student.GetAppUserID(), service.ProblemEventTypeAdd, problemType, addedIDs)
	if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
		return nil, liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
	}

	return addedIDs, nil
}

func (s *studentUsecaseProblem) updateProblem(ctx context.Context, student service.Student, workbook service.Workbook, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error {
	problemType := workbook.GetProblemType()
	if err := student.CheckQuota(ctx, problemType, "Size"); err != nil {
		return liberrors.Errorf("student.CheckQuota(size). err: %w", err)
	}
	if err := student.CheckQuota(ctx, problemType, "Update"); err != nil {
		return liberrors.Errorf("student.CheckQuota(udpate). err: %w", err)
	}
	addedIDs, updatedIDs, removedIDs, err := workbook.UpdateProblem(ctx, student, id, param)
	if err != nil {
		return liberrors.Errorf("failed to UpdateProblem. err: %w", err)
	}

	{
		event := service.NewProblemEvent(student.GetOrganizationID(), student.GetAppUserID(), service.ProblemEventTypeAdd, problemType, addedIDs)
		if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
			return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
		}
	}

	{
		event := service.NewProblemEvent(student.GetOrganizationID(), student.GetAppUserID(), service.ProblemEventTypeUpdate, problemType, updatedIDs)
		if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
			return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
		}
	}

	{
		event := service.NewProblemEvent(student.GetOrganizationID(), student.GetAppUserID(), service.ProblemEventTypeRemove, problemType, removedIDs)
		if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
			return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
		}
	}

	// 	if err := student.IncrementQuotaUsage(ctx, problemType, "Size", int(added)); err != nil {
	// 		return err
	// 	}
	// } else if added < 0 {
	// 	if err := student.DecrementQuotaUsage(ctx, problemType, "Size", -int(added)); err != nil {
	// 		return err
	// 	}
	// }
	// if updated > 0 {
	// 	if err := student.IncrementQuotaUsage(ctx, problemType, "Update", int(updated)); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (s *studentUsecaseProblem) updateProblemProperty(ctx context.Context, student service.Student, workbook service.Workbook, id domain.ProblemSelectParameter2, param domain.ProblemUpdateParameter) error {
	problemType := workbook.GetProblemType()
	if err := student.CheckQuota(ctx, problemType, "Size"); err != nil {
		return liberrors.Errorf("student.CheckQuota(size). err: %w", err)
	}
	if err := student.CheckQuota(ctx, problemType, "Update"); err != nil {
		return liberrors.Errorf("student.CheckQuota(udpate). err: %w", err)
	}
	addedIDs, updatedIDs, removedIDs, err := workbook.UpdateProblemProperty(ctx, student, id, param)
	if err != nil {
		return liberrors.Errorf("workbook.UpdateProblemProperty. err: %w", err)
	}

	{
		event := service.NewProblemEvent(student.GetOrganizationID(), userD.AppUserID(student.GetID()), service.ProblemEventTypeAdd, problemType, addedIDs)
		if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
			return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
		}
	}

	{
		event := service.NewProblemEvent(student.GetOrganizationID(), userD.AppUserID(student.GetID()), service.ProblemEventTypeUpdate, problemType, updatedIDs)
		if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
			return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
		}
	}

	{
		event := service.NewProblemEvent(student.GetOrganizationID(), userD.AppUserID(student.GetID()), service.ProblemEventTypeRemove, problemType, removedIDs)
		if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
			return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
		}
	}

	// if added > 0 {
	// 	if err := student.IncrementQuotaUsage(ctx, problemType, "Size", int(added)); err != nil {
	// 		return err
	// 	}
	// } else if added < 0 {
	// 	if err := student.DecrementQuotaUsage(ctx, problemType, "Size", -int(added)); err != nil {
	// 		return err
	// 	}
	// }
	// if updated > 0 {
	// 	if err := student.IncrementQuotaUsage(ctx, problemType, "Update", int(updated)); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
