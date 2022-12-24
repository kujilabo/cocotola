//go:generate mockery --output mock --name Space
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type Space interface {
	domain.SpaceModel
}

type space struct {
	domain.SpaceModel
}

func NewSpace(spaceModel domain.SpaceModel) (Space, error) {
	m := &space{
		spaceModel,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}
