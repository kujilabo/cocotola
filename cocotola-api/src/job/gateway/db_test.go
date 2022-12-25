package gateway_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

type testService struct {
	driverName string
	db         *gorm.DB
}

func testDB(t *testing.T, fn func(ctx context.Context, ts testService)) {
	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		logrus.Debugf("%s\n", driverName)
		sqlDB, err := db.DB()
		require.NoError(t, err)
		defer sqlDB.Close()

		testService := testService{driverName: driverName, db: db}

		fn(ctx, testService)
	}
}

func setupJob(t *testing.T, ts testService) {
}

func teardownJob(t *testing.T, ts testService) {
	ts.db.Exec("delete from job_status")
}
