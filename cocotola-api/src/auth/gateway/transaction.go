package gateway

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	"gorm.io/gorm"
)

type transaction struct {
	db      *gorm.DB
	userRff userG.RepositoryFactoryFunc
}

func NewTransaction(db *gorm.DB, userRff userG.RepositoryFactoryFunc) (service.Transaction, error) {
	return &transaction{
		db:      db,
		userRff: userRff,
	}, nil
}

func (t *transaction) Do(ctx context.Context, fn func(userRf userS.RepositoryFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error { // nolint:wrapcheck
		userRf, err := t.userRff(ctx, tx)
		if err != nil {
			return err
		}
		return fn(userRf)
	})
}
