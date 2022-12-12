package gateway

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	"gorm.io/gorm"
)

type repositoryFactory struct {
	db *gorm.DB
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
	if db == nil {
		return nil, libD.ErrInvalidArgument
	}

	return &repositoryFactory{
		db: db,
	}, nil
}

func (f *repositoryFactory) NewJobStatusRepository(ctx context.Context) (service.JobStatusRepository, error) {
	return NewJobStatusRepository(ctx, f.db)
}

func (f *repositoryFactory) NewJobHistoryRepository(ctx context.Context) (service.JobHistoryRepository, error) {
	return NewJobHistoryRepository(ctx, f.db)
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
	return t.db.Transaction(func(tx *gorm.DB) error {
		rf, err := t.rff(ctx, tx)
		if err != nil {
			return err
		}
		return fn(rf)
	})
}
