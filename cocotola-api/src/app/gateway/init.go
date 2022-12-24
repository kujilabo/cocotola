package gateway

import (
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/kujilabo/cocotola/cocotola-api/src/app/gateway")

var jst *time.Location

func init() {
	jst = time.Now().Local().Location()
}
