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
	rf RepositoryFactory
	domain.AppUserModel
}

func NewAppUser(rf RepositoryFactory, appUserModel domain.AppUserModel) (AppUser, error) {
	if rf == nil {
		return nil, liberrors.Errorf("rf is nil. err: %w", libD.ErrInvalidArgument)
	}
	if appUserModel == nil {
		return nil, liberrors.Errorf("appUserModel is nil. err: %w", libD.ErrInvalidArgument)
	}

	m := &appUser{
		rf:           rf,
		AppUserModel: appUserModel,
	}

	return m, libD.Validator.Struct(m)
}

func (a *appUser) GetDefaultSpace(ctx context.Context) (Space, error) {
	spaceRepo, err := a.rf.NewSpaceRepository()
	if err != nil {
		return nil, err
	}
	return spaceRepo.FindDefaultSpace(ctx, a)
}

func (a *appUser) GetPersonalSpace(ctx context.Context) (Space, error) {
	spaceRepo, err := a.rf.NewSpaceRepository()
	if err != nil {
		return nil, err
	}
	return spaceRepo.FindPersonalSpace(ctx, a)
}
