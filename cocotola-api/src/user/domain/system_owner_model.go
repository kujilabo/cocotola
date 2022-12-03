package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
)

const SystemOwnerID = 2

type SystemOwnerModel interface {
	AppUserModel
	IsSystemOwnerModel() bool
}

type systemOwnerModel struct {
	AppUserModel
}

func NewSystemOwnerModel(appUser AppUserModel) (SystemOwnerModel, error) {
	m := &systemOwnerModel{
		AppUserModel: appUser,
	}

	return m, libD.Validator.Struct(m)
}

func (s *systemOwnerModel) IsSystemOwnerModel() bool {
	return true
}
