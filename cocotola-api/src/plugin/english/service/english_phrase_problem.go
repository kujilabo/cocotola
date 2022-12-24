package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type EnglishPhraseProblem interface {
	domain.EnglishPhraseProblemModel
	service.ProblemFeature
}

type englishPhraseProblem struct {
	domain.EnglishPhraseProblemModel
	service.ProblemFeature
}

func NewEnglishPhraseProblem(problemModel domain.EnglishPhraseProblemModel, problem service.ProblemFeature) (EnglishPhraseProblem, error) {
	m := &englishPhraseProblem{
		EnglishPhraseProblemModel: problemModel,
		ProblemFeature:            problem,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
