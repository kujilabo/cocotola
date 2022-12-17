package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

// type AppUserID uint

type AppUser interface {
	domain.AppUserModel

	GetDefaultSpace(ctx context.Context) (Space, error)
	GetPersonalSpace(ctx context.Context) (Space, error)
}

type appUser struct {
	domain.AppUserModel
	rf        RepositoryFactory
	spaceRepo SpaceRepository
}

func NewAppUser(ctx context.Context, rf RepositoryFactory, appUserModel domain.AppUserModel) (AppUser, error) {
	if rf == nil {
		return nil, liberrors.Errorf("rf is nil. err: %w", libD.ErrInvalidArgument)
	}
	if appUserModel == nil {
		return nil, liberrors.Errorf("appUserModel is nil. err: %w", libD.ErrInvalidArgument)
	}
	spaceRepo, err := rf.NewSpaceRepository(ctx)
	if err != nil {
		return nil, err
	}

	m := &appUser{
		AppUserModel: appUserModel,
		rf:           rf,
		spaceRepo:    spaceRepo,
	}

	return m, libD.Validator.Struct(m)
}

func (m *appUser) GetDefaultSpace(ctx context.Context) (Space, error) {
	return m.spaceRepo.FindDefaultSpace(ctx, m)
}

func (m *appUser) GetPersonalSpace(ctx context.Context) (Space, error) {
	return m.spaceRepo.FindPersonalSpace(ctx, m)
}
