//go:generate mockery --output mock --name UserQuotaRepository
package service

import (
	"context"
	"errors"

	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type QuotaUnit string
type QuotaName string

var (
	ErrQuotaExceeded              = errors.New("quota exceeded")
	QuotaUnitPersitance QuotaUnit = "persitance"
	QuotaUnitMonth      QuotaUnit = "month"
	QuotaUnitDay        QuotaUnit = "day"
	QuotaNameSize       QuotaName = "Size"
	QuotaNameUpdate     QuotaName = "Update"
)

type UserQuotaRepository interface {
	IsExceeded(ctx context.Context, organizationID userD.OrganizationID, appUserID userD.AppUserID, name string, unit QuotaUnit, limit int) (bool, error)

	Increment(ctx context.Context, organizationID userD.OrganizationID, appUserID userD.AppUserID, name string, unit QuotaUnit, limit int, count int) (bool, error)
}
