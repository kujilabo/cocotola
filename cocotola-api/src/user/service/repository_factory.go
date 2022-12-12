//go:generate mockery --output mock --name RepositoryFactory
package service

import "context"

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) (OrganizationRepository, error)
	NewSpaceRepository(ctx context.Context) (SpaceRepository, error)
	NewAppUserRepository(ctx context.Context) (AppUserRepository, error)
	NewAppUserGroupRepository(ctx context.Context) (AppUserGroupRepository, error)

	NewGroupUserRepository(ctx context.Context) (GroupUserRepository, error)
	NewUserSpaceRepository(ctx context.Context) (UserSpaceRepository, error)
	NewRBACRepository(ctx context.Context) (RBACRepository, error)
}
