//go:generate mockery --output mock --name StudentModel
package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type StudentModel interface {
	userD.AppUserModel
	IsStudentModel() bool
}

type studentModel struct {
	userD.AppUserModel
}

func NewStudentModel(appUserModel userD.AppUserModel) (StudentModel, error) {
	m := &studentModel{
		AppUserModel: appUserModel,
	}

	return m, libD.Validator.Struct(m)
}

func (m *studentModel) IsStudentModel() bool {
	return true
}
