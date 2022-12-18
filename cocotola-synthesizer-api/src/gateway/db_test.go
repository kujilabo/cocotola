package gateway_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/gateway"
	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/service"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

type testService struct {
	driverName string
	db         *gorm.DB
	rf         service.RepositoryFactory
}

func testDB(t *testing.T, fn func(ctx context.Context, ts testService)) {
	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		logrus.Debugf("%s\n", driverName)
		sqlDB, err := db.DB()
		require.NoError(t, err)
		defer sqlDB.Close()

		rf, err := gateway.NewRepositoryFactory(ctx, db, driverName)
		require.NoError(t, err)

		testService := testService{driverName: driverName, db: db, rf: rf}

		fn(ctx, testService)
	}
}
