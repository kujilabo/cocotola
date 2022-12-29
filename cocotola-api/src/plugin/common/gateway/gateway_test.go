package gateway_test

import (
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

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
