package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

const EnglishSentenceProblemType = "english_sentence"

type EnglishSentenceProblem interface {
	domain.EnglishSentenceProblemModel
	service.ProblemFeature
}

type englishSentenceProblem struct {
	domain.EnglishSentenceProblemModel
	service.ProblemFeature
}

func NewEnglishSentenceProblem(problemModel domain.EnglishSentenceProblemModel, problem service.ProblemFeature) (EnglishSentenceProblem, error) {
	m := &englishSentenceProblem{
		EnglishSentenceProblemModel: problemModel,
		ProblemFeature:              problem,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
