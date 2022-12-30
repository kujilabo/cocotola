//go:generate mockery --output mock --name StudentUsecaseWorkbook
package student

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/usecase"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

const DefaultPageNo = 1
const DefaultPageSize = 10

type StudentUsecaseWorkbook interface {
	FindWorkbooks(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (domain.WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workBookID domain.WorkbookID) (domain.WorkbookModel, error)

	AddWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, parameter domain.WorkbookAddParameter) (domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int, parameter domain.WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int) error
}

type studentUsecaseWorkbook struct {
	transaction service.Transaction
	pf          service.ProcessorFactory
}

func NewStudentUsecaseWorkbook(transaction service.Transaction, pf service.ProcessorFactory) StudentUsecaseWorkbook {
	return &studentUsecaseWorkbook{
		transaction: transaction,
		pf:          pf,
	}
}

func (s *studentUsecaseWorkbook) FindWorkbooks(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (domain.WorkbookSearchResult, error) {
	var result domain.WorkbookSearchResult
	fn := func(student service.Student) error {
		condition, err := domain.NewWorkbookSearchCondition(DefaultPageNo, DefaultPageSize, []userD.SpaceID{})
		if err != nil {
			return liberrors.Errorf("service.NewWorkbookSearchCondition. err: %w", err)
		}

		tmpResult, err := student.FindWorkbooksFromPersonalSpace(ctx, condition)
		if err != nil {
			return liberrors.Errorf("student.FindWorkbooksFromPersonalSpace. err: %w", err)
		}

		result = tmpResult
		return nil
	}

	if err := s.studentHandle(ctx, organizationID, operatorID, fn); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *studentUsecaseWorkbook) FindWorkbookByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workBookID domain.WorkbookID) (domain.WorkbookModel, error) {
	var result domain.WorkbookModel
	fn := func(student service.Student) error {
		tmpResult, err := student.FindWorkbookByID(ctx, workBookID)
		if err != nil {
			return liberrors.Errorf("student.FindWorkbookByID. err: %w", err)
		}
		result = tmpResult
		return nil
	}

	if err := s.studentHandle(ctx, organizationID, operatorID, fn); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *studentUsecaseWorkbook) AddWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, parameter domain.WorkbookAddParameter) (domain.WorkbookID, error) {
	var addedWorkbookID domain.WorkbookID
	fn := func(student service.Student) error {
		tmpAddedWorkbookID, err := student.AddWorkbookToPersonalSpace(ctx, parameter)
		if err != nil {
			return liberrors.Errorf("student.AddWorkbookToPersonalSpace. err: %w", err)
		}
		addedWorkbookID = tmpAddedWorkbookID
		return nil
	}

	if err := s.studentHandle(ctx, organizationID, operatorID, fn); err != nil {
		return 0, err
	}

	return addedWorkbookID, nil
}

func (s *studentUsecaseWorkbook) UpdateWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int, parameter domain.WorkbookUpdateParameter) error {
	fn := func(student service.Student) error {
		if err := student.UpdateWorkbook(ctx, workbookID, version, parameter); err != nil {
			return liberrors.Errorf("student.UpdateWorkbook. err: %w", err)
		}
		return nil
	}
	return s.studentHandle(ctx, organizationID, operatorID, fn)
}

func (s *studentUsecaseWorkbook) RemoveWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int) error {
	fn := func(student service.Student) error {
		if err := student.RemoveWorkbook(ctx, workbookID, version); err != nil {
			return liberrors.Errorf("student.RemoveWorkbook. err: %w", err)
		}
		return nil
	}
	return s.studentHandle(ctx, organizationID, operatorID, fn)
}

func (s *studentUsecaseWorkbook) studentHandle(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, fn func(service.Student) error) error {
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, err := usecase.FindStudent(ctx, s.pf, rf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("usecase.FindStudent. err: %w", err)
		}
		return fn(student)
	}); err != nil {
		return liberrors.Errorf("studentHandle. err: %w", err)
	}
	return nil
}
