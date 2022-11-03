package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type AppUserGroup interface {
	domain.AppUserGroupModel
}

type appUserGroup struct {
	domain.AppUserGroupModel
}

// NewAppUserGroup returns a new AppUserGroup
func NewAppUserGroup(appUserGroupModel domain.AppUserGroupModel) (AppUserGroup, error) {
	m := &appUserGroup{
		appUserGroupModel,
	}

	return m, libD.Validator.Struct(m)
}
