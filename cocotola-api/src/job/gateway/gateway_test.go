package gateway_test

import (
	"crypto/rand"
	"math/big"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

// var userRff func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	testlibG.InitMySQL(sqls.SQL, "127.0.0.1", 3307)
	testlibG.InitSQLite(sqls.SQL)

	// userRff = func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
	// 	return gateway.NewRepositoryFactory(ctx, db)
	// }

	// service.InitSystemAdmin(userRff)
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
