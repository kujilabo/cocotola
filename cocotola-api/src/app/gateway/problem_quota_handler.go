package gateway

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
)

type problemQuotaHandler struct {
	rf service.RepositoryFactory
	pf service.ProcessorFactory
}

func NewProblemQuotaHandler(pf service.ProcessorFactory, rf service.RepositoryFactory) service.ProblemQuotaHandler {
	return &problemQuotaHandler{
		pf: pf,
		rf: rf,
	}
}

func (p *problemQuotaHandler) Update(ctx context.Context, event service.ProblemEvent) error {
	organizationID := event.GetOrganizationID()
	appUserID := event.GetAppUserID()
	problemType := event.GetProblemType()
	problemEventType := event.GetProblemEventType()
	value := len(event.GetProblemIDs())
	switch problemEventType {
	case service.ProblemEventTypeAdd:
		processor, err := p.pf.NewProblemQuotaProcessor(problemType)
		userQuotaRepo := p.rf.NewUserQuotaRepository(ctx)
		if err != nil {
			return err
		}

		{
			unit := processor.GetUnitForSizeQuota()
			limit := processor.GetLimitForSizeQuota()
			isExceeded, err := userQuotaRepo.Increment(ctx, organizationID, appUserID, problemType+"_size", unit, limit, value)
			if err != nil {
				return err
			}

			if isExceeded {
				return service.ErrQuotaExceeded
			}
		}
		{
			unit := processor.GetUnitForUpdateQuota()
			limit := processor.GetLimitForUpdateQuota()
			isExceeded, err := userQuotaRepo.Increment(ctx, organizationID, appUserID, problemType+"_update", unit, limit, value)
			if err != nil {
				return err
			}

			if isExceeded {
				return service.ErrQuotaExceeded
			}

		}
	default:
		return errors.New("invalid problem event type")
	}
	return nil
}
