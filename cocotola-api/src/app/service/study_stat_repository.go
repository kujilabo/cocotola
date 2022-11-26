package service

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type StudyStatRepository interface {
	AggregateResultsOfAllUsers(ctx context.Context, operator domain.SystemOwnerModel, targetDate time.Time) error
}
