package gateway_test

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
	"github.com/sirupsen/logrus"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	testlibG.InitMySQL(sqls.SQL, "127.0.0.1", 3307)
	testlibG.InitSQLite(sqls.SQL)

	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		logrus.Debugf("%s\n", driverName)
		rbacRepo := userG.NewRBACRepository(ctx, db)
		err := rbacRepo.Init()
		if err != nil {
			panic(err)
		}
	}
	// userS.InitSystemAdmin(userRfFunc)
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			panic(err)
		}
		b[i] = letterRunes[val.Int64()]
	}
	return string(b)
}
