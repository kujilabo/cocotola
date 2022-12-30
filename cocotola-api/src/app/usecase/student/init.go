package student

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/kujilabo/cocotola/cocotola-api/src/app/usecase/student")

type FindStudentFunc func(context.Context, service.RepositoryFactory, userD.OrganizationID, userD.AppUserID) (service.Student, error)
