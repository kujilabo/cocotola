//go:generate mockery --output mock --name Stat
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
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

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
