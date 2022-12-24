package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type OrganizationID uint

type OrganizationModel interface {
	Model
	GetName() string
}

type organizationModel struct {
	Model
	Name string `validate:"required"`
}

func NewOrganizationModel(model Model, name string) (OrganizationModel, error) {
	m := &organizationModel{
		Model: model,
		Name:  name,
	}
	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *organizationModel) GetName() string {
	return m.Name
}

func (m *organizationModel) String() string {
	return m.Name
}
