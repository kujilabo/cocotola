package gateway_test

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

func init() {
	testlibG.InitMySQL(sqls.SQL, "127.0.0.1", 3317)
	testlibG.InitSQLite(sqls.SQL)
}

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
