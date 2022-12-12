package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
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

func (e *jobHistoryEntity) toModel() (service.JobHistory, error) {
	return service.NewJobHistory()
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

func (r *jobHistoryRepository) FindJobHistoryByJobName(ctx context.Context, jobName domain.JobName) (service.JobHistory, error) {
	jobHistoryEntity := jobHistoryEntity{}

	if result := r.db.Where("job_name = ?", (string)(jobName)).
		First(&jobHistoryEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrJobHistoryNotFound
		}
		return nil, result.Error
	}

	return jobHistoryEntity.toModel()
}
