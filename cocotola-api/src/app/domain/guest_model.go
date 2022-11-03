package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type GuestModel interface {
	userD.AppUserModel
}

type guestModel struct {
	userD.AppUserModel
}

func NewGuestModel(appUser userD.AppUserModel) (GuestModel, error) {
	m := &guestModel{
		AppUserModel: appUser,
	}

	return m, libD.Validator.Struct(m)
}
