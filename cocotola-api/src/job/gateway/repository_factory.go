package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type repositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	return &repositoryFactory{
		db: db,
	}, nil
}

func (f *repositoryFactory) NewJobStatusRepository(ctx context.Context) service.JobStatusRepository {
	return newJobStatusRepository(ctx, f.db)
}

func (f *repositoryFactory) NewJobHistoryRepository(ctx context.Context) service.JobHistoryRepository {
	return newJobHistoryRepository(ctx, f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)

type transaction struct {
	db  *gorm.DB
	rff RepositoryFactoryFunc
}

func NewTransaction(db *gorm.DB, rff RepositoryFactoryFunc) service.Transaction {
	return &transaction{
		db:  db,
		rff: rff,
	}
}

func (t *transaction) Do(ctx context.Context, fn func(rf service.RepositoryFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error { // nolint:wrapcheck
		rf, err := t.rff(ctx, tx)
		if err != nil {
			return liberrors.Errorf("rff. err: %w", err)
		}
		return fn(rf)
	})
}
