package service

import (
	"context"
	"strconv"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	pluginS "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type englishPhraseProblemAddParemeter struct {
	Text       string `validate:"required"`
	Lang2      string `validate:"required"`
	Translated string `validate:"required"`
}

func toEnglishPhraseProblemAddParemeter(param appD.ProblemAddParameter) (*englishPhraseProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, liberrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, liberrors.Errorf("translated is not defined. err: %w", libD.ErrInvalidArgument)
	}

	m := &englishPhraseProblemAddParemeter{
		Lang2:      param.GetProperties()["lang2"],
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

type EnglishPhraseProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemRemoveProcessor
}

type englishPhraseProblemProcessor struct {
	synthesizerClient appS.SynthesizerClient
	translatorClient  pluginS.TranslatorClient
}

func NewEnglishPhraseProblemProcessor(synthesizerClient appS.SynthesizerClient, translatorClient pluginS.TranslatorClient) EnglishPhraseProblemProcessor {
	return &englishPhraseProblemProcessor{
		synthesizerClient: synthesizerClient,
		translatorClient:  translatorClient,
	}
}

func (p *englishPhraseProblemProcessor) AddProblem(ctx context.Context, repo appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, param appD.ProblemAddParameter) ([]appD.ProblemID, []appD.ProblemID, []appD.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishPhraseProblemType)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishPhraseProblemAddParemeter(param)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to toNewEnglishPhraseProblemParemeter. err: %w", err)
	}

	audio, err := p.synthesizerClient.Synthesize(ctx, appD.Lang2EN, extractedParam.Text)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("p.synthesizerClient.Synthesize. err: %w", err)
	}

	problemID, err := p.addSingleProblem(ctx, operator, problemRepo, param, extractedParam, appD.AudioID(audio.GetAudioModel().GetID()))
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to addSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return []appD.ProblemID{problemID}, nil, nil, nil
}

func (p *englishPhraseProblemProcessor) addSingleProblem(ctx context.Context, operator appD.StudentModel, problemRepo appS.ProblemRepository, param appD.ProblemAddParameter, extractedParam *englishPhraseProblemAddParemeter, audioID appD.AudioID) (appD.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem1")

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := map[string]string{
		"text":       extractedParam.Text,
		"translated": extractedParam.Translated,
		"lang2":      extractedParam.Lang2,
		"audioId":    strconv.Itoa(int(audioID)),
	}
	newParam, err := appD.NewProblemAddParameter(param.GetWorkbookID() /*param.GetNumber(),*/, properties)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, liberrors.Errorf("failed to problemRepo.AddProblem. err: %w", err)
	}

	return problemID, nil

}

func (p *englishPhraseProblemProcessor) RemoveProblem(ctx context.Context, repo appS.RepositoryFactory, operator appD.StudentModel, id appD.ProblemSelectParameter2) ([]appD.ProblemID, []appD.ProblemID, []appD.ProblemID, error) {
	problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishPhraseProblemType)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
		return nil, nil, nil, liberrors.Errorf("problemRepo.RemoveProblem. err: %w", err)
	}

	return nil, nil, []appD.ProblemID{id.GetProblemID()}, nil
}
