//go:generate mockery --output mock --name GuestStudent
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
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
		return nil, err
	}

	spaceRepo, err := userRf.NewSpaceRepository(ctx)
	if err != nil {
		return nil, err
	}

	m := &guestStudent{
		StudentModel: studentModel,
		pf:           pf,
		rf:           rf,
		spaceRepo:    spaceRepo,
	}

	return m, libD.Validator.Struct(m)
}

func (s *guestStudent) GetDefaultSpace(ctx context.Context) (userS.Space, error) {
	return s.spaceRepo.FindDefaultSpace(ctx, s)
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
