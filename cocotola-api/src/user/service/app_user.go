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
	spaceRepo SpaceRepository
}

func NewAppUser(ctx context.Context, rf RepositoryFactory, appUserModel domain.AppUserModel) (AppUser, error) {
	if rf == nil {
		return nil, liberrors.Errorf("rf is nil. err: %w", libD.ErrInvalidArgument)
	}
	if appUserModel == nil {
		return nil, liberrors.Errorf("appUserModel is nil. err: %w", libD.ErrInvalidArgument)
	}
	spaceRepo := rf.NewSpaceRepository(ctx)

	m := &appUser{
		AppUserModel: appUserModel,
		spaceRepo:    spaceRepo,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *appUser) GetDefaultSpace(ctx context.Context) (Space, error) {
	space, err := m.spaceRepo.FindDefaultSpace(ctx, m)
	if err != nil {
		return nil, liberrors.Errorf("m.spaceRepo.FindDefaultSpace. err: %w", err)
	}

	return space, nil
}

func (m *appUser) GetPersonalSpace(ctx context.Context) (Space, error) {
	space, err := m.spaceRepo.FindPersonalSpace(ctx, m)
	if err != nil {
		return nil, liberrors.Errorf("m.spaceRepo.FindPersonalSpace. err: %w", err)
	}

	return space, nil
}
