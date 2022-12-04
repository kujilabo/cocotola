package gateway_test

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

var userRfFunc func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	testlibG.InitMySQL(sqls.SQL, "127.0.0.1", 3307)
	testlibG.InitSQLite(sqls.SQL)
	rand.Seed(time.Now().UnixNano())

	userRfFunc = func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userS.InitSystemAdmin(userRfFunc)
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
