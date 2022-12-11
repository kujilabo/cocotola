package gateway

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	"gorm.io/gorm"
)

type transaction struct {
	db     *gorm.DB
	rfFunc service.RepositoryFactoryFunc
}

func NewTransaction(db *gorm.DB, rfFunc service.RepositoryFactoryFunc) (service.Transaction, error) {
	return &transaction{
		db:     db,
		rfFunc: rfFunc,
	}, nil
}

func (t *transaction) Do(ctx context.Context, fn func(rf service.RepositoryFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		rf, err := t.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		return fn(rf)
	})
}
