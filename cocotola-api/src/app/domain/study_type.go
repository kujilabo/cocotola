//go:generate mockery --output mock --name StudyType
package domain

import libD "github.com/kujilabo/cocotola/lib/domain"

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

	return m, libD.Validator.Struct(m)
}

func (m *studyType) GetID() uint {
	return m.ID
}

func (m *studyType) GetName() StudyTypeName {
	return StudyTypeName(m.Name)
}
