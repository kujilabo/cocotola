package service

import (
	"context"
	"errors"
	"io"
	"strconv"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	pluginS "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

var (
	EnglishSentenceProblemQuotaSizeUnit                     = appS.QuotaUnitPersitance
	EnglishSentenceProblemQuotaSizeLimit                    = 5000
	EnglishSentenceProblemQuotaUpdateUnit                   = appS.QuotaUnitDay
	EnglishSentenceProblemQuotaUpdateLimit                  = 100
	EnglishSentenceProblemAddPropertyAudioID                = "audioId"
	EnglishSentenceProblemAddPropertyLang2                  = "lang2"
	EnglishSentenceProblemAddPropertyText                   = "text"
	EnglishSentenceProblemAddPropertyTranslated             = "translated"
	EnglishSentenceProblemAddPropertyProvider               = "provider"
	EnglishSentenceProblemAddPropertyTatoebaSentenceNumber1 = "tatoebaSentenceNumber1"
	EnglishSentenceProblemAddPropertyTatoebaSentenceNumber2 = "tatoebaSentenceNumber2"
	EnglishSentenceProblemAddPropertyTatoebaAuthor1         = "tatoebaAuthor1"
	EnglishSentenceProblemAddPropertyTatoebaAuthor2         = "tatoebaAuthor2"
)

type englishSentenceProblemAddParemeter struct {
	Lang2                  appD.Lang2 `validate:"required"`
	Text                   string     `validate:"required"`
	Translated             string
	Provider               string
	TatoebaSentenceNumber1 int
	TatoebaSentenceNumber2 int
	TatoebaAuthor1         string
	TatoebaAuthor2         string
}

func (p *englishSentenceProblemAddParemeter) toProperties(audioID appD.AudioID) map[string]string {
	return map[string]string{
		EnglishSentenceProblemAddPropertyAudioID:                strconv.Itoa(int(uint(audioID))),
		EnglishSentenceProblemAddPropertyLang2:                  p.Lang2.String(),
		EnglishSentenceProblemAddPropertyText:                   p.Text,
		EnglishSentenceProblemAddPropertyTranslated:             p.Translated,
		EnglishSentenceProblemAddPropertyProvider:               p.Provider,
		EnglishSentenceProblemAddPropertyTatoebaSentenceNumber1: strconv.Itoa(p.TatoebaSentenceNumber1),
		EnglishSentenceProblemAddPropertyTatoebaSentenceNumber2: strconv.Itoa(p.TatoebaSentenceNumber2),
		EnglishSentenceProblemAddPropertyTatoebaAuthor1:         p.TatoebaAuthor1,
		EnglishSentenceProblemAddPropertyTatoebaAuthor2:         p.TatoebaAuthor2,
	}
}

func toEnglishSentenceProblemAddParemeter(param appD.ProblemAddParameter) (*englishSentenceProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["translated"]; !ok {
		return nil, liberrors.Errorf("translated is not defined. err: %w", libD.ErrInvalidArgument)
	}
	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, liberrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	lang2, err := appD.NewLang2(param.GetProperties()["lang2"])
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	m := &englishSentenceProblemAddParemeter{
		Lang2:      lang2,
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

type EnglishSentenceProblemProcessor interface {
	appS.ProblemAddProcessor
	appS.ProblemUpdateProcessor
	appS.ProblemRemoveProcessor
	appS.ProblemImportProcessor
	appS.ProblemQuotaProcessor
}

type englishSentenceProblemProcessor struct {
	synthesizerClient               appS.SynthesizerClient
	translatorClient                pluginS.TranslatorClient
	newProblemAddParameterCSVReader func(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator
}

func NewEnglishSentenceProblemProcessor(synthesizerClient appS.SynthesizerClient, translatorClient pluginS.TranslatorClient, newProblemAddParameterCSVReader func(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator) EnglishSentenceProblemProcessor {
	return &englishSentenceProblemProcessor{
		synthesizerClient:               synthesizerClient,
		translatorClient:                translatorClient,
		newProblemAddParameterCSVReader: newProblemAddParameterCSVReader,
	}
}

func (p *englishSentenceProblemProcessor) AddProblem(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, param appD.ProblemAddParameter) ([]appD.ProblemID, []appD.ProblemID, []appD.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Info("englishSentenceProblemProcessor.AddProblem")

	problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	extractedParam, err := toEnglishSentenceProblemAddParemeter(param)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to toNewEnglishSentenceProblemParemeter. err: %w", err)
	}

	audioID := appD.AudioID(0)
	if workbook.GetProperties()["audioEnabled"] == "true" {
		logger.Infof("audioEnabled is true")
		audio, err := p.synthesizerClient.Synthesize(ctx, appD.Lang2EN, extractedParam.Text)
		if err != nil {
			return nil, nil, nil, liberrors.Errorf("p.synthesizerClient.Synthesize. err: %w", err)
		}

		audioID = appD.AudioID(audio.GetAudioModel().GetID())
	} else {
		logger.Infof("audioEnabled is false")
	}

	problemID, err := p.addSingleProblem(ctx, operator, problemRepo, param, extractedParam, audioID)
	if err != nil {
		return nil, nil, nil, liberrors.Errorf("failed to addSingleProblem: extractedParam: %+v, err: %w", extractedParam, err)
	}

	return []appD.ProblemID{problemID}, nil, nil, nil
}

func (p *englishSentenceProblemProcessor) addSingleProblem(ctx context.Context, operator appD.StudentModel, problemRepo appS.ProblemRepository, param appD.ProblemAddParameter, extractedParam *englishSentenceProblemAddParemeter, audioID appD.AudioID) (appD.ProblemID, error) {
	logger := log.FromContext(ctx)

	logger.Infof("text: %s, audio ID: %d", extractedParam.Text, audioID)

	properties := extractedParam.toProperties(audioID)
	newParam, err := appD.NewProblemAddParameter(param.GetWorkbookID() /*param.GetNumber(),*/, properties)
	if err != nil {
		return 0, liberrors.Errorf("failed to NewParameter. err: %w", err)
	}

	problemID, err := problemRepo.AddProblem(ctx, operator, newParam)
	if err != nil {
		return 0, liberrors.Errorf("problemRepo.AddProblem. err: %w, newParam: %+v", err, newParam.GetProperties())
	}

	return problemID, nil
}

func (p *englishSentenceProblemProcessor) UpdateProblem(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, id appD.ProblemSelectParameter2, param appD.ProblemUpdateParameter) ([]appD.ProblemID, []appD.ProblemID, []appD.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("englishSentenceProblemProcessor.UpdateProblem, param: %+v", param)

	return nil, nil, nil, errors.New("NotImplemented")
}

func (p *englishSentenceProblemProcessor) UpdateProblemProperty(ctx context.Context, rf appS.RepositoryFactory, operator appD.StudentModel, workbook appD.WorkbookModel, id appD.ProblemSelectParameter2, param appD.ProblemUpdateParameter) ([]appD.ProblemID, []appD.ProblemID, []appD.ProblemID, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("englishSentenceProblemProcessor.UpdateProblem, param: %+v", param)

	fn := func() error {
		problemRepo, err := rf.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
		if err != nil {
			return liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
		}
		if err := problemRepo.UpdateProblemProperty(ctx, operator, id, param); err != nil {
			return liberrors.Errorf(". err: %w", err)
		}
		return nil
	}
	if err := fn(); err != nil {
		return nil, nil, nil, err
	}

	return nil, []appD.ProblemID{id.GetProblemID()}, nil, nil
}

func (p *englishSentenceProblemProcessor) RemoveProblem(ctx context.Context, repo appS.RepositoryFactory, operator appD.StudentModel, id appD.ProblemSelectParameter2) ([]appD.ProblemID, []appD.ProblemID, []appD.ProblemID, error) {
	fn := func() error {
		problemRepo, err := repo.NewProblemRepository(ctx, domain.EnglishSentenceProblemType)
		if err != nil {
			return liberrors.Errorf("failed to NewProblemRepository. err: %w", err)
		}

		if err := problemRepo.RemoveProblem(ctx, operator, id); err != nil {
			return liberrors.Errorf("failed to RemoveProblem. err: %w", err)
		}

		return nil
	}
	if err := fn(); err != nil {
		return nil, nil, nil, err
	}

	return []appD.ProblemID{id.GetProblemID()}, nil, nil, nil
}

func (p *englishSentenceProblemProcessor) CreateCSVReader(ctx context.Context, workbookID appD.WorkbookID, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
	return p.newProblemAddParameterCSVReader(workbookID, reader), nil
}

func (p *englishSentenceProblemProcessor) GetUnitForSizeQuota() appS.QuotaUnit {
	return EnglishSentenceProblemQuotaSizeUnit
}

func (p *englishSentenceProblemProcessor) GetLimitForSizeQuota() int {
	return EnglishSentenceProblemQuotaSizeLimit
}

func (p *englishSentenceProblemProcessor) GetUnitForUpdateQuota() appS.QuotaUnit {
	return EnglishSentenceProblemQuotaUpdateUnit
}

func (p *englishSentenceProblemProcessor) GetLimitForUpdateQuota() int {
	return EnglishSentenceProblemQuotaUpdateLimit
}
