package gateway_test

import (
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

func init() {
	testlibG.InitMySQL(sqls.SQL, "127.0.0.1", 3327)
	testlibG.InitSQLite(sqls.SQL)
}
