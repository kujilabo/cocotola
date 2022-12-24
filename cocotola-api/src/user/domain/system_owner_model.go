package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
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

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *systemOwnerModel) IsSystemOwnerModel() bool {
	return true
}
