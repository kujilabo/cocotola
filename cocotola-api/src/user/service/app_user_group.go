package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
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

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
