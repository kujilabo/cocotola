package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/usecase"
)

func Test_userUsecase_DictionaryLookup_azureRepo(t *testing.T) {

	fn := func(t *testing.T, ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		azureRepo := ts.rf.NewAzureTranslationRepository(ctx)
		userUsecase := usecase.NewUserUsecase(ctx, ts.transaction, ts.azureTranslationClient)
		// given
		// - azureRepo has one data
		azureRepoResults := []service.AzureTranslation{{
			Pos:        domain.PosNoun,
			Target:     "本ar",
			Confidence: 1,
		}}
		err := azureRepo.Add(ctx, domain.Lang2JA, "book", azureRepoResults)
		assert.NoError(t, err)

		// when
		actual, err := userUsecase.DictionaryLookup(ctx, domain.Lang2EN, domain.Lang2JA, "book")
		assert.NoError(t, err)
		// then
		assert.Equal(t, len(actual), 1)
		assert.Equal(t, actual[0].GetTranslated(), "本ar")
	}
	test(t, fn)
}

func Test_userUsecase_DictionaryLookup_azureClient(t *testing.T) {
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		userUsecase := usecase.NewUserUsecase(ctx, ts.transaction, ts.azureTranslationClient)
		// given
		// - azureClient has one data
		azureClientResults := []service.AzureTranslation{{
			Pos:        domain.PosNoun,
			Target:     "本ar",
			Confidence: 1,
		}}
		ts.azureTranslationClient.On("DictionaryLookup", ctx, "book", domain.Lang2EN, domain.Lang2JA).Return(azureClientResults, nil)

		// when
		actual, err := userUsecase.DictionaryLookup(ctx, domain.Lang2EN, domain.Lang2JA, "book")
		assert.NoError(t, err)

		// then
		assert.Equal(t, len(actual), 1)
		assert.Equal(t, actual[0].GetTranslated(), "本ar")
	}
	test(t, fn)
}

func Test_userUsecase_DictionaryLookup_azureRepo_azureClient(t *testing.T) {
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		azureRepo := ts.rf.NewAzureTranslationRepository(ctx)
		userUsecase := usecase.NewUserUsecase(ctx, ts.transaction, ts.azureTranslationClient)

		// given
		// - azureRepo has one data
		azureRepoResults := []service.AzureTranslation{{
			Pos:        domain.PosNoun,
			Target:     "本ar",
			Confidence: 1,
		}}
		err := azureRepo.Add(ctx, domain.Lang2JA, "book", azureRepoResults)
		assert.NoError(t, err)

		// - azureClient has onedata
		azureClientResults := []service.AzureTranslation{{
			Pos:        domain.PosNoun,
			Target:     "本ac",
			Confidence: 1,
		}}
		ts.azureTranslationClient.On("DictionaryLookup", ctx, "book", domain.Lang2EN, domain.Lang2JA).Return(azureClientResults, nil)

		// when
		actual, err := userUsecase.DictionaryLookup(ctx, domain.Lang2EN, domain.Lang2JA, "book")
		assert.NoError(t, err)

		// then
		// - the translation registered in auzreRepo is selected
		assert.Equal(t, len(actual), 1)
		assert.Equal(t, actual[0].GetTranslated(), "本ar")
	}
	test(t, fn)
}

func Test_userUsecase_DictionaryLookup_custom_azureRepo(t *testing.T) {
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		azureRepo := ts.rf.NewAzureTranslationRepository(ctx)
		customRepo := ts.rf.NewCustomTranslationRepository(ctx)
		userUsecase := usecase.NewUserUsecase(ctx, ts.transaction, ts.azureTranslationClient)

		// given
		// - customRepo has one data
		param, err := service.NewTransalationAddParameter("book", domain.PosNoun, domain.Lang2JA, "本c")
		assert.NoError(t, err)
		err = customRepo.Add(ctx, param)
		assert.NoError(t, err)
		// - azureRepo has two data. One is a noun word and the other is a verb word.
		azureRepoResults := []service.AzureTranslation{
			{
				Pos:        domain.PosNoun,
				Target:     "本ar",
				Confidence: 1,
			},
			{
				Pos:        domain.PosVerb,
				Target:     "予約するar",
				Confidence: 1,
			},
		}
		err = azureRepo.Add(ctx, domain.Lang2JA, "book", azureRepoResults)
		assert.NoError(t, err)

		// when
		actual, err := userUsecase.DictionaryLookup(ctx, domain.Lang2EN, domain.Lang2JA, "book")
		assert.NoError(t, err)

		// then
		// - Noun: the translation registered in customRepo is selected because customRepo has higher priority than azureRepo.
		// - Verb: the translation registered in azureRepo is selected because customRepo does not have translations for verb.
		assert.Equal(t, len(actual), 2)
		assert.Equal(t, actual[0].GetTranslated(), "本c")
		assert.Equal(t, actual[1].GetTranslated(), "予約するar")
	}
	test(t, fn)
}
