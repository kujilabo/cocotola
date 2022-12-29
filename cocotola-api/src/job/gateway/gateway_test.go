package gateway_test

import (
	"crypto/rand"
	"math/big"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

// var userRff func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
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
