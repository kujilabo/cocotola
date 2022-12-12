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

func (f *repositoryFactory) NewOrganizationRepository(ctx context.Context) (service.OrganizationRepository, error) {
	return NewOrganizationRepository(ctx, f.db)
}

func (f *repositoryFactory) NewSpaceRepository(ctx context.Context) (service.SpaceRepository, error) {
	return NewSpaceRepository(ctx, f.db)
}

func (f *repositoryFactory) NewAppUserRepository(ctx context.Context) (service.AppUserRepository, error) {
	return NewAppUserRepository(ctx, f, f.db)
}

func (f *repositoryFactory) NewAppUserGroupRepository(ctx context.Context) (service.AppUserGroupRepository, error) {
	return NewAppUserGroupRepository(ctx, f.db)
}

func (f *repositoryFactory) NewGroupUserRepository(ctx context.Context) (service.GroupUserRepository, error) {
	return NewGroupUserRepository(ctx, f.db)
}

func (f *repositoryFactory) NewUserSpaceRepository(ctx context.Context) (service.UserSpaceRepository, error) {
	return NewUserSpaceRepository(ctx, f, f.db)
}

func (f *repositoryFactory) NewRBACRepository(ctx context.Context) (service.RBACRepository, error) {
	return NewRBACRepository(ctx, f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
