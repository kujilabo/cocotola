package domain

import libD "github.com/kujilabo/cocotola/lib/domain"

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

	return m, libD.Validator.Struct(m)
}

func (m *problemType) GetID() uint {
	return m.ID
}

func (m *problemType) GetName() ProblemTypeName {
	return ProblemTypeName(m.Name)
}
