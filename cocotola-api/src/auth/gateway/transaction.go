package gateway

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	"gorm.io/gorm"
)

type transaction struct {
	db         *gorm.DB
	userRfFunc userS.RepositoryFactoryFunc
}

func NewTransaction(db *gorm.DB, userRfFunc userS.RepositoryFactoryFunc) (service.Transaction, error) {
	return &transaction{
		db:         db,
		userRfFunc: userRfFunc,
	}, nil
}

func (t *transaction) Do(ctx context.Context, fn func(userRf userS.RepositoryFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		userRf, err := t.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		return fn(userRf)
	})
}
