//go:generate mockery --output mock --name SystemStudentModel
package domain

import (
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type SystemStudentModel interface {
	userD.AppUserModel
}

type systemStudentModel struct {
	userD.AppUserModel
}

func NewSystemStudentModel(appUser userD.AppUserModel) (SystemStudentModel, error) {
	m := &systemStudentModel{
		AppUserModel: appUser,
	}

	return m, libD.Validator.Struct(m)
}
