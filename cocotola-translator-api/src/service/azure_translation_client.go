//go:generate mockery --output mock --name AzureTranslationClient
package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
)

type AzureTranslationClient interface {
	DictionaryLookup(ctx context.Context, text string, fromLang, toLang domain.Lang2) ([]AzureTranslation, error)
}
