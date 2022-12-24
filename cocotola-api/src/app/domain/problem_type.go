package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type ProblemTypeName string

type ProblemType interface {
	GetID() uint
	GetName() ProblemTypeName
}

type problemType struct {
	ID   uint   `validate:"required,gte=1"`
	Name string `validate:"required"`
}

func NewProblemType(id uint, name string) (ProblemType, error) {
	m := &problemType{
		ID:   id,
		Name: name,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *problemType) GetID() uint {
	return m.ID
}

func (m *problemType) GetName() ProblemTypeName {
	return ProblemTypeName(m.Name)
}
