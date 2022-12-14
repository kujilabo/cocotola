package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
)

type repositoryFactory struct {
	db         *gorm.DB
	driverName string
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string) (service.RepositoryFactory, error) {
	return &repositoryFactory{
		db:         db,
		driverName: driverName,
	}, nil
}

func (f *repositoryFactory) NewAzureTranslationRepository(ctx context.Context) service.AzureTranslationRepository {
	return newAzureTranslationRepository(f.db)
}

func (f *repositoryFactory) NewCustomTranslationRepository(ctx context.Context) service.CustomTranslationRepository {
	return newCustomTranslationRepository(f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
