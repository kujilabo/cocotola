//go:generate mockery --output mock --name SystemStudent
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type SystemStudent interface {
	domain.SystemStudentModel

	FindWorkbookFromSystemSpace(ctx context.Context, name string) (Workbook, error)

	AddWorkbookToSystemSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error)
}

type systemStudent struct {
	domain.SystemStudentModel
	rf RepositoryFactory
}

func NewSystemStudent(rf RepositoryFactory, systemStudentModel domain.SystemStudentModel) (SystemStudent, error) {
	m := &systemStudent{
		SystemStudentModel: systemStudentModel,
		rf:                 rf,
	}

	return m, libD.Validator.Struct(m)
}

func (s *systemStudent) FindWorkbookFromSystemSpace(ctx context.Context, name string) (Workbook, error) {
	systemSpaceID := GetSystemSpaceID()

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	workbook, err := workbookRepo.FindWorkbookByName(ctx, s, systemSpaceID, name)
	if err != nil {
		return nil, liberrors.Errorf("failed to FindWorkbookByName. err: %w", err)
	}

	return workbook, nil
}

func (s *systemStudent) AddWorkbookToSystemSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error) {
	systemSpaceID := GetSystemSpaceID()
	if uint(systemSpaceID) == 0 {
		return 0, errors.New("invalid system space ID")
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	workbookID, err := workbookRepo.AddWorkbook(ctx, s, systemSpaceID, parameter)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddWorkbook. err: %w", err)
	}

	return workbookID, nil
}
