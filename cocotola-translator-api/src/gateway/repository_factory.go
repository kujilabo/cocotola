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

type transaction struct {
	db  *gorm.DB
	rff RepositoryFactoryFunc
}

func NewTransaction(db *gorm.DB, rff RepositoryFactoryFunc) (service.Transaction, error) {
	return &transaction{
		db:  db,
		rff: rff,
	}, nil
}

func (t *transaction) Do(ctx context.Context, fn func(rf service.RepositoryFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error { // nolint:wrapcheck
		rf, err := t.rff(ctx, tx)
		if err != nil {
			return err // nolint:wrapcheck
		}
		return fn(rf)
	})
}
