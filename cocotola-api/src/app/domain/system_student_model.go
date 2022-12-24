//go:generate mockery --output mock --name SystemStudentModel
package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type SystemStudentModel interface {
	userD.AppUserModel
	IsSystemStudentModel() bool
}

type systemStudentModel struct {
	userD.AppUserModel
}

func NewSystemStudentModel(appUser userD.AppUserModel) (SystemStudentModel, error) {
	m := &systemStudentModel{
		AppUserModel: appUser,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *systemStudentModel) IsSystemStudentModel() bool {
	return true
}
