package gateway_test

import (
	"context"
	"testing"
)

func Test_a(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(ctx context.Context, ts testService) {
	}
	testDB(t, fn)
}
