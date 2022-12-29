package gateway_test

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	fns := []func() (*gorm.DB, error){
		func() (*gorm.DB, error) {
			return testlibG.InitMySQL(sqls.SQL, "127.0.0.1", 3307)
		},
		func() (*gorm.DB, error) {
			return testlibG.InitSQLiteInFile(sqls.SQL)
		},
	}

	for _, fn := range fns {
		db, err := fn()
		if err != nil {
			panic(err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.Close()
	}

	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		logrus.Debugf("%s\n", driverName)
		rbacRepo := userG.NewRBACRepository(ctx, db)
		err := rbacRepo.Init()
		if err != nil {
			panic(err)
		}
	}
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
