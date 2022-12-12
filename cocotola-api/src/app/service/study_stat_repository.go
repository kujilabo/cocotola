package service

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type StudyStatRepository interface {
	AggregateResults(ctx context.Context, operator domain.SystemOwnerModel, targetDate time.Time, userID domain.AppUserID) error
	CleanStudyStats(ctx context.Context, operator domain.SystemOwnerModel, expirationDate time.Time) error
}
