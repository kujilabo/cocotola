//go:generate mockery --output mock --name ToEnglishWordProblemAddParameter
package service

import (
	"context"
	"errors"
	"strconv"

	appD "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	pluginS "github.com/kujilabo/cocotola/cocotola-api/src/plugin/common/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type ToEnglishWordProblemAddParameter interface {
	Run(ctx context.Context) ([]appD.ProblemAddParameter, error)
}

type ToEnglishWordProblemUpdateParameter interface {
	Run(ctx context.Context) ([]appD.ProblemUpdateParameter, error)
}

type toSingleEnglishWordProblemAddParameter struct {
	translatorClient pluginS.TranslatorClient
	workbookID       appD.WorkbookID
	// number           int
	param   *EnglishWordProblemAddParemeter
	audioID appD.AudioID
}

func NewToSingleEnglishWordProblemAddParameter(translatorClient pluginS.TranslatorClient, workbookID appD.WorkbookID, param *EnglishWordProblemAddParemeter, audioID appD.AudioID) ToEnglishWordProblemAddParameter {
	return &toSingleEnglishWordProblemAddParameter{
		translatorClient: translatorClient,
		workbookID:       workbookID,
		// number:           number,
		param:   param,
		audioID: audioID,
	}
}

func (c *toSingleEnglishWordProblemAddParameter) Run(ctx context.Context) ([]appD.ProblemAddParameter, error) {
	translated := c.param.Translated
	if translated == "" {
		translation, err := c.translatorClient.DictionaryLookupWithPos(ctx, appD.Lang2EN, c.param.Lang2, c.param.Text, c.param.Pos)
		if err != nil {
			if !errors.Is(err, pluginS.ErrTranslationNotFound) {
				return nil, liberrors.Errorf(". err: %w", err)
			}
		} else {
			translated = translation.GetTranslated()
		}
	}

	properties := c.param.toProperties()
	properties[EnglishSentenceProblemAddPropertyTranslated] = translated
	properties[EnglishWordProblemAddPropertyAudioID] = strconv.Itoa(int(uint(c.audioID)))

	param, err := appD.NewProblemAddParameter(c.workbookID /*c.number,*/, properties)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
	}

	return []appD.ProblemAddParameter{param}, nil
}

type toMultipleEnglishWordProblemAddParameter struct {
	translatorClient pluginS.TranslatorClient
	workbookID       appD.WorkbookID
	// number           int
	param   *EnglishWordProblemAddParemeter
	audioID appD.AudioID
}

func NewToMultipleEnglishWordProblemAddParameter(translatorClient pluginS.TranslatorClient, workbookID appD.WorkbookID, param *EnglishWordProblemAddParemeter, audioID appD.AudioID) ToEnglishWordProblemAddParameter {
	return &toMultipleEnglishWordProblemAddParameter{
		translatorClient: translatorClient,
		workbookID:       workbookID,
		// number:           number,
		param:   param,
		audioID: audioID,
	}
}

func (c *toMultipleEnglishWordProblemAddParameter) Run(ctx context.Context) ([]appD.ProblemAddParameter, error) {
	logger := log.FromContext(ctx)

	translated, err := c.translatorClient.DictionaryLookup(ctx, appD.Lang2EN, c.param.Lang2, c.param.Text)
	if errors.Is(err, pluginS.ErrTranslationNotFound) {
		return nil, liberrors.Errorf("c.translatorClient.DictionaryLookup. err: %w", err)
	}

	if len(translated) == 0 || err != nil {
		logger.Errorf("translate err: %v", err)

		properties := c.param.toProperties()
		properties[EnglishWordProblemAddPropertyAudioID] = strconv.Itoa(int(uint(c.audioID)))

		param, err := appD.NewProblemAddParameter(c.workbookID /*c.number,*/, properties)
		if err != nil {
			return nil, liberrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		return []appD.ProblemAddParameter{param}, nil
	}

	logger.Infof("translated: %v", translated)

	params := make([]appD.ProblemAddParameter, len(translated))
	for i, t := range translated {

		properties := c.param.toProperties()
		properties[EnglishWordProblemAddPropertyAudioID] = strconv.Itoa(int(uint(c.audioID)))
		properties[EnglishWordProblemAddPropertyTranslated] = t.GetTranslated()
		properties[EnglishWordProblemAddPropertyPos] = strconv.Itoa(int(t.GetPos()))

		param, err := appD.NewProblemAddParameter(c.workbookID /*c.number,*/, properties)
		if err != nil {
			return nil, liberrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		params[i] = param
	}

	return params, nil
}

type toSingleEnglishWordProblemUpdateParameter struct {
	translatorClient pluginS.TranslatorClient
	number           int
	param            *EnglishWordProblemUpdateParemeter
	audioID          appD.AudioID
	sentenceID1      appD.ProblemID
}

func NewToSingleEnglishWordProblemUpdateParameter(translatorClient pluginS.TranslatorClient, number int, param *EnglishWordProblemUpdateParemeter, audioID appD.AudioID, sentenceID1 appD.ProblemID) ToEnglishWordProblemUpdateParameter {
	return &toSingleEnglishWordProblemUpdateParameter{
		translatorClient: translatorClient,
		number:           number,
		param:            param,
		audioID:          audioID,
		sentenceID1:      sentenceID1,
	}
}

func (c *toSingleEnglishWordProblemUpdateParameter) Run(ctx context.Context) ([]appD.ProblemUpdateParameter, error) {
	// translated := c.param.Translated
	// if translated == "" {
	// 	translation, err := c.translatorClient.DictionaryLookupWithPos(ctx, appD.Lang2EN, appD.Lang2JA, c.param.Text, c.param.Pos)
	// 	if err != nil {
	// 		if !errors.Is(err, pluginS.ErrTranslationNotFound) {
	// 			return nil, err
	// 		}
	// 	} else {
	// 		translated = translation.GetTranslated()
	// 	}
	// }

	properties := map[string]string{
		EnglishWordProblemUpdatePropertyText:        c.param.Text,
		EnglishWordProblemUpdatePropertyTranslated:  c.param.Translated,
		EnglishWordProblemUpdatePropertyAudioID:     strconv.Itoa(int(c.audioID)),
		EnglishWordProblemUpdatePropertySentenceID1: strconv.Itoa(int(c.sentenceID1)),
	}

	param, err := appD.NewProblemUpdateParameter( /*c.number,*/ properties)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
	}

	return []appD.ProblemUpdateParameter{param}, nil
}
