//go:generate mockery --output mock --name Space
package service

import (
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
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

	return m, libD.Validator.Struct(m)
}
