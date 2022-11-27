package gateway

import (
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type repositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(db *gorm.DB) (service.RepositoryFactory, error) {
	return &repositoryFactory{
		db: db,
	}, nil
}

func (f *repositoryFactory) NewOrganizationRepository() (service.OrganizationRepository, error) {
	return NewOrganizationRepository(f.db)
}

func (f *repositoryFactory) NewSpaceRepository() (service.SpaceRepository, error) {
	return NewSpaceRepository(f.db)
}

func (f *repositoryFactory) NewAppUserRepository() (service.AppUserRepository, error) {
	return NewAppUserRepository(f, f.db)
}

func (f *repositoryFactory) NewAppUserGroupRepository() (service.AppUserGroupRepository, error) {
	return NewAppUserGroupRepository(f.db)
}

func (f *repositoryFactory) NewGroupUserRepository() (service.GroupUserRepository, error) {
	return NewGroupUserRepository(f.db)
}

func (f *repositoryFactory) NewUserSpaceRepository() (service.UserSpaceRepository, error) {
	return NewUserSpaceRepository(f, f.db)
}

func (f *repositoryFactory) NewRBACRepository() (service.RBACRepository, error) {
	return NewRBACRepository(f.db)
}
