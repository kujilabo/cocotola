package service

import (
	"context"
	"time"

	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type StudyStatRepository interface {
	AggregateResults(ctx context.Context, operator userD.SystemOwnerModel, targetDate time.Time, userID userD.AppUserID) error
	CleanStudyStats(ctx context.Context, operator userD.SystemOwnerModel, expirationDate time.Time) error
}
