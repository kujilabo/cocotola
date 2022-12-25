//go:generate mockery --output mock --name UserUsecase
package usecase

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type UserUsecase interface {
	DictionaryLookup(ctx context.Context, fromLang, toLang domain.Lang2, text string) ([]domain.Translation, error)

	DictionaryLookupWithPos(ctx context.Context, fromLang, toLang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error)
}

type userUsecase struct {
	transaction            service.Transaction
	azureTranslationClient service.AzureTranslationClient
}

type UserPresenter interface {
	WriteTranslations(ctx context.Context, translations []domain.Translation) error
	WriteTranslation(ctx context.Context, translation domain.Translation) error
}

func NewUserUsecase(ctx context.Context, transaction service.Transaction, azureTranslationClient service.AzureTranslationClient) UserUsecase {
	return &userUsecase{
		transaction:            transaction,
		azureTranslationClient: azureTranslationClient,
	}
}

func (u *userUsecase) selectMaxConfidenceTranslations(ctx context.Context, in []service.AzureTranslation) map[domain.WordPos]service.AzureTranslation {
	results := make(map[domain.WordPos]service.AzureTranslation)
	for _, i := range in {
		if _, ok := results[i.Pos]; !ok {
			results[i.Pos] = i
		} else if i.Confidence > results[i.Pos].Confidence {
			results[i.Pos] = i
		}
	}
	return results
}

func (u *userUsecase) customDictionaryLookup(ctx context.Context, customRepo service.CustomTranslationRepository, text string, fromLang, toLang domain.Lang2) ([]domain.Translation, error) {
	customContained, err := customRepo.Contain(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	if !customContained {
		return nil, service.ErrTranslationNotFound
	}

	customResults, err := customRepo.FindByText(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	return customResults, nil
}

func (u *userUsecase) azureDictionaryLookup(ctx context.Context, azureRepo service.AzureTranslationRepository, fromLang, toLang domain.Lang2, text string) ([]service.AzureTranslation, error) {
	azureContained, err := azureRepo.Contain(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	if azureContained {
		azureResults, err := azureRepo.Find(ctx, toLang, text)
		if err != nil {
			return nil, err
		}
		return azureResults, nil
	}

	azureResults, err := u.azureTranslationClient.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}

	if len(azureResults) == 0 {
		return azureResults, nil
	}

	if err := azureRepo.Add(ctx, toLang, text, azureResults); err != nil {
		return nil, liberrors.Errorf("failed to add auzre_translation. err: %w", err)
	}

	return azureResults, nil
}

func (u *userUsecase) DictionaryLookup(ctx context.Context, fromLang, toLang domain.Lang2, text string) ([]domain.Translation, error) {
	results := make([]domain.Translation, 0)
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		customRepo := rf.NewCustomTranslationRepository(ctx)

		// find translations from custom reopository
		customResults, err := u.customDictionaryLookup(ctx, customRepo, text, fromLang, toLang)
		if err != nil && !errors.Is(err, service.ErrTranslationNotFound) {
			return err
		}
		// if !errors.Is(err, service.ErrTranslationNotFound) {
		// 	return customResults, err
		// }

		azureRepo := rf.NewAzureTranslationRepository(ctx)

		// find translations from azure
		azureResults, err := u.azureDictionaryLookup(ctx, azureRepo, fromLang, toLang, text)
		if err != nil {
			return err
		}
		azureResultMap := u.selectMaxConfidenceTranslations(ctx, azureResults)
		makeKey := func(text string, pos domain.WordPos) string {
			return text + "_" + strconv.Itoa(int(pos))
		}
		resultMap := make(map[string]domain.Translation)

		// insert customResults into resultMap
		for _, c := range customResults {
			key := makeKey(c.GetText(), c.GetPos())
			resultMap[key] = c
		}

		// insert azureResultMap into resultMap
		for _, a := range azureResultMap {
			key := makeKey(text, a.Pos)
			if _, ok := resultMap[key]; ok {
				continue
			}

			result, err := domain.NewTranslation(1, time.Now(), time.Now(), text, a.Pos, fromLang, a.Target, "azure")
			if err != nil {
				return err
			}
			resultMap[key] = result
		}

		// convert map to list
		tmpResults := make([]domain.Translation, 0)
		for _, v := range resultMap {
			tmpResults = append(tmpResults, v)
		}

		sort.Slice(tmpResults, func(i, j int) bool { return tmpResults[i].GetPos() < tmpResults[j].GetPos() })

		results = tmpResults
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (u *userUsecase) DictionaryLookupWithPos(ctx context.Context, fromLang, toLang domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	results, err := u.DictionaryLookup(ctx, fromLang, toLang, text)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.GetPos() == pos {
			return r, nil
		}
	}
	return nil, service.ErrTranslationNotFound
}
