package service

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola/cocotola-api/src/plugin/english/service")
