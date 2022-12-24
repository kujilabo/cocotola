//go:generate mockery --output mock --name GuestStudent
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type GuestStudent interface {
	domain.StudentModel

	GetDefaultSpace(ctx context.Context) (userS.Space, error)

	FindWorkbooksFromPublicSpace(ctx context.Context, condition WorkbookSearchCondition) (WorkbookSearchResult, error)
}

type guestStudent struct {
	domain.StudentModel
	rf        RepositoryFactory
	pf        ProcessorFactory
	spaceRepo userS.SpaceRepository
}

func NewGuestStudent(ctx context.Context, pf ProcessorFactory, rf RepositoryFactory, studentModel domain.StudentModel) (GuestStudent, error) {
	userRf, err := rf.NewUserRepositoryFactory(ctx)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	spaceRepo := userRf.NewSpaceRepository(ctx)

	m := &guestStudent{
		StudentModel: studentModel,
		pf:           pf,
		rf:           rf,
		spaceRepo:    spaceRepo,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (s *guestStudent) GetDefaultSpace(ctx context.Context) (userS.Space, error) {
	space, err := s.spaceRepo.FindDefaultSpace(ctx, s)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return space, nil
}

func (s *guestStudent) FindWorkbooksFromPublicSpace(ctx context.Context, condition WorkbookSearchCondition) (WorkbookSearchResult, error) {
	return nil, errors.New("aaa")
	// space, err := s.GetPersonalSpace(ctx)
	// if err != nil {
	// 	return nil, liberrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	// }

	// // specify space
	// newCondition, err := NewWorkbookSearchCondition(condition.GetPageNo(), condition.GetPageSize(), []userD.SpaceID{userD.SpaceID(space.GetID())})
	// if err != nil {
	// 	return nil, liberrors.Errorf("failed to NewWorkbookSearchCondition. err: %w", err)
	// }

	// workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	// if err != nil {
	// 	return nil, liberrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	// }

	// return workbookRepo.FindPersonalWorkbooks(ctx, s, newCondition)
}
