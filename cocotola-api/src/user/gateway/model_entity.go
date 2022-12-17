package gateway

import (
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
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
	return domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
}

type JunctionModelEntity struct {
	CreatedAt time.Time
	CreatedBy uint
}

// func (e *junctionModelEntity) toModel() (domain.Model, error) {
// 	return domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
// }
