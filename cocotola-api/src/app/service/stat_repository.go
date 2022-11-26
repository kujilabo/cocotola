package service

import (
	"context"

	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type StatRepository interface {
	FindStat(ctx context.Context, operatorID userD.AppUserID) (Stat, error)
}
