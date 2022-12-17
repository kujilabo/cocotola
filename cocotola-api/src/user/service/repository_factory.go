//go:generate mockery --output mock --name RepositoryFactory
package service

import "context"

type RepositoryFactory interface {
	NewOrganizationRepository(ctx context.Context) OrganizationRepository
	NewSpaceRepository(ctx context.Context) SpaceRepository
	NewAppUserRepository(ctx context.Context) AppUserRepository
	NewAppUserGroupRepository(ctx context.Context) AppUserGroupRepository

	NewGroupUserRepository(ctx context.Context) GroupUserRepository
	NewUserSpaceRepository(ctx context.Context) UserSpaceRepository
	NewRBACRepository(ctx context.Context) RBACRepository
}
