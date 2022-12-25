package gateway_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

func Test_xxx_FindByText(t *testing.T) {
	t.Parallel()
	// logrus.SetLevel(logrus.DebugLevel)

	for driverName, db := range testlibG.ListDB() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
	}
}
