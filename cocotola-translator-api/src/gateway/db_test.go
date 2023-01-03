package gateway_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/gateway"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

type testService struct {
	driverName string
	db         *gorm.DB
	rf         service.RepositoryFactory
}

func testDB(t *testing.T, fn func(t *testing.T, ctx context.Context, ts testService)) {
	logrus.SetLevel(logrus.WarnLevel)

	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		driverName := driverName
		db := db
		t.Run(driverName, func(t *testing.T) {
			// t.Parallel()
			sqlDB, err := db.DB()
			require.NoError(t, err)
			defer sqlDB.Close()

			rf, err := gateway.NewRepositoryFactory(ctx, db, driverName)
			require.NoError(t, err)
			testService := testService{driverName: driverName, db: db, rf: rf}

			fn(t, ctx, testService)
		})
	}
}

func teardownDB(t *testing.T, ts testService) {
	result := ts.db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from azure_translation")
	assert.NoError(t, result.Error)
	result = ts.db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from custom_translation")
	assert.NoError(t, result.Error)
}
