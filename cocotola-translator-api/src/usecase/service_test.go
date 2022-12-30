package usecase_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-translator-api/src/gateway"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	service_mock "github.com/kujilabo/cocotola/cocotola-translator-api/src/service/mock"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/sqls"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

type testService struct {
	db                     *gorm.DB
	rf                     service.RepositoryFactory
	transaction            service.Transaction
	azureTranslationClient *service_mock.AzureTranslationClient
}

func test(t *testing.T, fn func(t *testing.T, ctx context.Context, ts testService)) {
	logrus.SetLevel(logrus.WarnLevel)

	ctx := context.Background()

	db, err := testlibG.OpenSQLiteInMemory(sqls.SQL)
	require.NoError(t, err)
	rf, err := gateway.NewRepositoryFactory(ctx, db, "sqlite3")
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()
	rff := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, db, "sqlite3")
	}
	transaction, err := gateway.NewTransaction(db, rff)
	require.NoError(t, err)
	azureTranslationClient := new(service_mock.AzureTranslationClient)
	testService := testService{
		db:                     db,
		rf:                     rf,
		transaction:            transaction,
		azureTranslationClient: azureTranslationClient,
	}

	fn(t, ctx, testService)
}
