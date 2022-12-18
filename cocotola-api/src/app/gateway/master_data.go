package gateway

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type ProblemTypes interface {
	ToProblemTypeID(problemType domain.ProblemTypeName) (uint, error)
	ToProblemType(studyTypeID uint) (domain.ProblemTypeName, error)
}

type problemTypes struct {
	data []domain.ProblemType
}

func NewProblemTypes(data []domain.ProblemType) ProblemTypes {
	return &problemTypes{data: data}
}

func (r *problemTypes) ToProblemTypeID(problemType domain.ProblemTypeName) (uint, error) {
	for _, m := range r.data {
		if m.GetName() == problemType {
			return m.GetID(), nil
		}
	}
	return 0, libD.ErrInvalidArgument
}

func (r *problemTypes) ToProblemType(problemTypeID uint) (domain.ProblemTypeName, error) {
	for _, m := range r.data {
		if m.GetID() == problemTypeID {
			return m.GetName(), nil
		}
	}
	return "", libD.ErrInvalidArgument
}

type StudyTypes interface {
	ToStudyTypeID(studyType domain.StudyTypeName) (uint, error)
	ToStudyType(studyTypeID uint) (domain.StudyTypeName, error)
	Values() []domain.StudyType
}

type studyTypes struct {
	data []domain.StudyType
}

func NewStudyTypes(data []domain.StudyType) StudyTypes {
	return &studyTypes{data: data}
}

func (r *studyTypes) ToStudyTypeID(studyType domain.StudyTypeName) (uint, error) {
	for _, m := range r.data {
		if m.GetName() == studyType {
			return m.GetID(), nil
		}
	}
	return 0, libD.ErrInvalidArgument
}

func (r *studyTypes) ToStudyType(studyTypeID uint) (domain.StudyTypeName, error) {
	for _, m := range r.data {
		if m.GetID() == studyTypeID {
			return m.GetName(), nil
		}
	}

	return "", libD.ErrInvalidArgument
}

func (r *studyTypes) Values() []domain.StudyType {
	return r.data
}
