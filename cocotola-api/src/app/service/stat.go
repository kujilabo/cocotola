//go:generate mockery --output mock --name Stat
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type Stat interface {
	domain.StatModel
}

type stat struct {
	domain.StatModel
}

func NewStat(statModel domain.StatModel) (Stat, error) {
	m := &stat{
		StatModel: statModel,
	}

	return m, libD.Validator.Struct(m)
}
