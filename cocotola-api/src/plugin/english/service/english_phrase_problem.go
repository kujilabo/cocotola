package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
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

	return m, libD.Validator.Struct(m)
}
