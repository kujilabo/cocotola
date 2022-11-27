package student

import (
	"context"
	"errors"
	"io"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/usecase"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type StudentUsecaseProblem interface {
	// problem
	FindProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param service.ProblemSearchCondition) (service.ProblemSearchResult, error)

	FindAllProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (service.ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param service.ProblemIDsCondition) (service.ProblemSearchResult, error)

	FindProblemByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter1) (domain.ProblemModel, error)

	FindProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error)

	AddProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, param service.ProblemAddParameter) ([]domain.ProblemID, error)

	UpdateProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error

	UpdateProblemProperty(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter2) error

	ImportProblems(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, newIterator func(workbookID domain.WorkbookID, problemType string) (service.ProblemAddParameterIterator, error)) error
}

type studentUsecaseProblem struct {
	db             *gorm.DB
	pf             service.ProcessorFactory
	rfFunc         service.RepositoryFactoryFunc
	userRfFunc     userS.RepositoryFactoryFunc
	problemMonitor service.ProblemMonitor
}

func NewStudentUsecaseProblem(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) StudentUsecaseProblem {
	return &studentUsecaseProblem{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *studentUsecaseProblem) FindProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param service.ProblemSearchCondition) (service.ProblemSearchResult, error) {
	var result service.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := workbook.FindProblems(ctx, student, param)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindAllProblemsByWorkbookID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (service.ProblemSearchResult, error) {
	var result service.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindAllProblems(ctx, student)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindProblemsByProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, param service.ProblemIDsCondition) (service.ProblemSearchResult, error) {
	var result service.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblemsByProblemIDs(ctx, student, param)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindProblemByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter1) (domain.ProblemModel, error) {
	var result domain.ProblemModel
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblemByID(ctx, student, id.GetProblemID())
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseProblem) FindProblemIDs(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error) {
	var result []domain.ProblemID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		tmpResult, err := workbook.FindProblemIDs(ctx, student)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseProblem) AddProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, param service.ProblemAddParameter) ([]domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	var result []domain.ProblemID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		studentService, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, param.GetWorkbookID())
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
		return nil, err
	}
	logger.Debug("problem id: %d", result)
	return result, nil
}

func (s *studentUsecaseProblem) UpdateProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		if err := s.updateProblem(ctx, student, workbook, id, param); err != nil {
			return liberrors.Errorf("failed to UpdateProblem. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *studentUsecaseProblem) UpdateProblemProperty(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		if err := s.updateProblemProperty(ctx, student, workbook, id, param); err != nil {
			return liberrors.Errorf("failed to UpdateProblem. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *studentUsecaseProblem) RemoveProblem(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, id service.ProblemSelectParameter2) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.RemoveProblem")

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, id.GetWorkbookID())
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		_, _, removedIDs, err := workbook.RemoveProblem(ctx, student, id)
		if err != nil {
			return liberrors.Errorf("workbook.RemoveProblem. err: %w", err)
		}
		problemType := workbook.GetProblemType()

		{
			event := service.NewProblemEvent(student.GetOrganizationID(), userD.AppUserID(student.GetID()), service.ProblemEventTypeRemove, problemType, removedIDs)
			if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
				return liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
			}
		}
		// if err := student.DecrementQuotaUsage(ctx, problemType, "Size", 1); err != nil {
		// 	return liberrors.Errorf("student.DecrementQuotaUsage. err: %w", err)
		// }

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *studentUsecaseProblem) ImportProblems(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, newIterator func(workbookID domain.WorkbookID, problemType string) (service.ProblemAddParameterIterator, error)) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.ImportProblems")

	var problemType string
	{
		_, workbook, err := s.findStudentAndWorkbook(ctx, s.db, organizationID, operatorID, workbookID)
		if err != nil {
			return liberrors.Errorf("s.findStudentAndWorkbook. err: %w", err)
		}
		problemType = workbook.GetProblemType()
	}
	iterator, err := newIterator(workbookID, problemType)
	if err != nil {
		return err
	}

	for {
		param, err := iterator.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if param == nil {
			continue
		}

		logger.Infof("param.properties: %+v", param.GetProperties())

		if err := s.db.Transaction(func(tx *gorm.DB) error {
			student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
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
			return err
		}
	}
	return nil
}

func (s *studentUsecaseProblem) findStudentAndWorkbook(ctx context.Context, tx *gorm.DB, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (service.Student, service.Workbook, error) {
	repo, err := s.rfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	userRepo, err := s.userRfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	student, err := usecase.FindStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}
	workbook, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, err
	}
	return student, workbook, nil
}

func (s *studentUsecaseProblem) addProblem(ctx context.Context, student service.Student, workbook service.Workbook, param service.ProblemAddParameter) ([]domain.ProblemID, error) {
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

	event := service.NewProblemEvent(student.GetOrganizationID(), userD.AppUserID(student.GetID()), service.ProblemEventTypeAdd, problemType, addedIDs)
	if err := s.problemMonitor.NotifyObservers(ctx, event); err != nil {
		return nil, liberrors.Errorf("student.IncrementQuotaUsage(Size). err: %w", err)
	}

	return addedIDs, nil
}

func (s *studentUsecaseProblem) updateProblem(ctx context.Context, student service.Student, workbook service.Workbook, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error {
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

func (s *studentUsecaseProblem) updateProblemProperty(ctx context.Context, student service.Student, workbook service.Workbook, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error {
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
