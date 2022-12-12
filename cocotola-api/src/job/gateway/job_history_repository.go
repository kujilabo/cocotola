package gateway

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	"gorm.io/gorm"
)

type jobHistoryEntity struct {
	JobStatusID  string
	JobName      string
	JobParameter string
	Status       string
	CreatedAt    time.Time
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

func (r *jobHistoryRepository) AddJobHistory(ctx context.Context, param service.JobHistoryAddParameter) error {
	jobHistoryEntity := jobHistoryEntity{
		JobStatusID:  (string)(param.GetJobStatusID()),
		JobName:      (string)(param.GetJobName()),
		JobParameter: param.GetJobParamter(),
		Status:       param.GetStatus(),
	}
	if result := r.db.Create(&jobHistoryEntity); result.Error != nil {
		return result.Error
	}

	return nil
}
