package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type GuestModel interface {
	userD.AppUserModel
	IsGuestModel() bool
}

type guestModel struct {
	userD.AppUserModel
}

func NewGuestModel(appUser userD.AppUserModel) (GuestModel, error) {
	m := &guestModel{
		AppUserModel: appUser,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *guestModel) IsGuestModel() bool {
	return true
}
