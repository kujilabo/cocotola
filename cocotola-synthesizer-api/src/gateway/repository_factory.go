package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-synthesizer-api/src/service"
)

type repositoryFactory struct {
	db         *gorm.DB
	driverName string
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string) (service.RepositoryFactory, error) {
	return &repositoryFactory{
		db:         db,
		driverName: driverName,
	}, nil
}

func (f *repositoryFactory) NewAudioRepository(ctx context.Context) service.AudioRepository {
	return NewAudioRepository(f.db)
}
