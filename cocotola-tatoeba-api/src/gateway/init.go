package gateway

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/gateway")
