//go:generate mockery --output mock --name StudyType
package domain

import (
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type StudyTypeName string

type StudyType interface {
	GetID() uint
	GetName() StudyTypeName
}

type studyType struct {
	ID   uint   `validate:"required,gte=1"`
	Name string `validate:"required"`
}

func NewStudyType(id uint, name string) (StudyType, error) {
	m := &studyType{
		ID:   id,
		Name: name,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *studyType) GetID() uint {
	return m.ID
}

func (m *studyType) GetName() StudyTypeName {
	return StudyTypeName(m.Name)
}
