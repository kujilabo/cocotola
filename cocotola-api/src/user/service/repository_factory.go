//go:generate mockery --output mock --name RepositoryFactory
package service

type RepositoryFactory interface {
	NewOrganizationRepository() (OrganizationRepository, error)
	NewSpaceRepository() (SpaceRepository, error)
	NewAppUserRepository() (AppUserRepository, error)
	NewAppUserGroupRepository() (AppUserGroupRepository, error)

	NewGroupUserRepository() (GroupUserRepository, error)
	NewUserSpaceRepository() (UserSpaceRepository, error)
	NewRBACRepository() (RBACRepository, error)
}
