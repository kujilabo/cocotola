package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
)

const SystemOwnerID = 2

type SystemOwnerModel interface {
	OwnerModel
	IsSystemOwnerModel() bool
}

type systemOwnerModel struct {
	OwnerModel
}

func NewSystemOwnerModel(appUser OwnerModel) (SystemOwnerModel, error) {
	m := &systemOwnerModel{
		OwnerModel: appUser,
	}

	return m, libD.Validator.Struct(m)
}

func (s *systemOwnerModel) IsSystemOwnerModel() bool {
	return true
}
