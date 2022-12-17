package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type repositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
	return &repositoryFactory{
		db: db,
	}, nil
}

func (f *repositoryFactory) NewOrganizationRepository(ctx context.Context) service.OrganizationRepository {
	return NewOrganizationRepository(ctx, f.db)
}

func (f *repositoryFactory) NewSpaceRepository(ctx context.Context) service.SpaceRepository {
	return NewSpaceRepository(ctx, f.db)
}

func (f *repositoryFactory) NewAppUserRepository(ctx context.Context) service.AppUserRepository {
	return NewAppUserRepository(ctx, f, f.db)
}

func (f *repositoryFactory) NewAppUserGroupRepository(ctx context.Context) service.AppUserGroupRepository {
	return NewAppUserGroupRepository(ctx, f.db)
}

func (f *repositoryFactory) NewGroupUserRepository(ctx context.Context) service.GroupUserRepository {
	return NewGroupUserRepository(ctx, f.db)
}

func (f *repositoryFactory) NewUserSpaceRepository(ctx context.Context) service.UserSpaceRepository {
	return NewUserSpaceRepository(ctx, f, f.db)
}

func (f *repositoryFactory) NewRBACRepository(ctx context.Context) service.RBACRepository {
	return NewRBACRepository(ctx, f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
