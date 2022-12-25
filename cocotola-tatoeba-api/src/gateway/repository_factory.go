package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/service"
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

func (f *repositoryFactory) NewTatoebaSentenceRepository(ctx context.Context) service.TatoebaSentenceRepository {
	return newTatoebaSentenceRepository(f.db)
}

func (f *repositoryFactory) NewTatoebaLinkRepository(ctx context.Context) service.TatoebaLinkRepository {
	return newTatoebaLinkRepository(f.db, f)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
