package gateway

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
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
	return 0, liberrors.Errorf("unsupported problemType. problemType: %s, err: %w", problemType, libD.ErrInvalidArgument)
}

func (r *problemTypes) ToProblemType(problemTypeID uint) (domain.ProblemTypeName, error) {
	for _, m := range r.data {
		if m.GetID() == problemTypeID {
			return m.GetName(), nil
		}
	}
	return "", liberrors.Errorf("unsupported problemTypeID. problemTypeID: %d, err: %w", problemTypeID, libD.ErrInvalidArgument)
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
	names := make([]string, 0)
	for _, data := range r.data {
		names = append(names, string(data.GetName()))
	}
	return 0, liberrors.Errorf("unsupported studyType. studyType: %s, studyTypes: %v, err: %w", studyType, names, libD.ErrInvalidArgument)
}

func (r *studyTypes) ToStudyType(studyTypeID uint) (domain.StudyTypeName, error) {
	for _, m := range r.data {
		if m.GetID() == studyTypeID {
			return m.GetName(), nil
		}
	}

	return "", liberrors.Errorf("unsupported studyTypeID. studyTypeID: %d, err: %w", studyTypeID, libD.ErrInvalidArgument)
}

func (r *studyTypes) Values() []domain.StudyType {
	return r.data
}
