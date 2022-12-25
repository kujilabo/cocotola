package usecase_test

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
)

type transaction struct {
	rf service.RepositoryFactory
}

func newTransaction(rf service.RepositoryFactory) service.Transaction {
	return &transaction{
		rf: rf,
	}
}

func (t *transaction) Do(ctx context.Context, fn func(rf service.RepositoryFactory) error) error {
	return fn(t.rf)
}
