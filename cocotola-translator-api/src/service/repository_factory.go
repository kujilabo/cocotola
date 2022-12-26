//go:generate mockery --output mock --name RepositoryFactory
//go:generate mockery --output mock --name Transaction
package service

import (
	"context"
)

type RepositoryFactory interface {
	NewAzureTranslationRepository(ctx context.Context) AzureTranslationRepository

	NewCustomTranslationRepository(ctx context.Context) CustomTranslationRepository
}
