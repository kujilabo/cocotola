package gateway_test

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

func init() {
	testlibG.InitMySQL(sqls.SQL)
	testlibG.InitSQLite(sqls.SQL)
}
