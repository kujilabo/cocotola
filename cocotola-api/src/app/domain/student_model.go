//go:generate mockery --output mock --name StudentModel
package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type StudentModel interface {
	userD.AppUserModel
	GetAppUserID() userD.AppUserID
	IsStudentModel() bool
}

type studentModel struct {
	userD.AppUserModel
}

func NewStudentModel(appUserModel userD.AppUserModel) (StudentModel, error) {
	m := &studentModel{
		AppUserModel: appUserModel,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *studentModel) GetAppUserID() userD.AppUserID {
	return userD.AppUserID(m.GetID())
}

func (m *studentModel) IsStudentModel() bool {
	return true
}
