package gateway

import (
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type SinmpleModelEntity struct {
	ID        uint
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint
	UpdatedBy uint
}

func (e *SinmpleModelEntity) toModel() (domain.Model, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return model, nil
}

type JunctionModelEntity struct {
	CreatedAt time.Time
	CreatedBy uint
}

// func (e *junctionModelEntity) toModel() (domain.Model, error) {
// 	return domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
// }
