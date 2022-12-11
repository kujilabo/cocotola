package gateway

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	"gorm.io/gorm"
)

type jobHistoryEntity struct {
}

func (e *jobHistoryEntity) TableName() string {
	return "job_history"
}

type jobHistoryRepository struct {
	db *gorm.DB
}

func NewJobHistoryRepository(ctx context.Context, db *gorm.DB) (service.JobHistoryRepository, error) {
	return &jobHistoryRepository{
		db: db,
	}, nil
}

func (r *jobHistoryRepository) AddJobHistory(ctx context.Context, jobHistory service.JobHistory) error {
	return nil
}
