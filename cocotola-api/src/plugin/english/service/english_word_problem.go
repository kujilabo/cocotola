package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

const EnglishWordProblemType = "english_word"

type EnglishWordProblem interface {
	domain.EnglishWordProblemModel
	service.ProblemFeature
}

type englishWordProblem struct {
	domain.EnglishWordProblemModel
	service.ProblemFeature
}

func NewEnglishWordProblem(problemModel domain.EnglishWordProblemModel, problem service.ProblemFeature) (EnglishWordProblem, error) {
	m := &englishWordProblem{
		EnglishWordProblemModel: problemModel,
		ProblemFeature:          problem,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
