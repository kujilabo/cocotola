package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
)

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
