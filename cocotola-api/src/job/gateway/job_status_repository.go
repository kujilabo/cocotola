package gateway

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
)

type jobStatusEntity struct {
	ID                 string
	JobName            string
	JobParameter       string
	ConcurrencyKey     *string
	ExpirationDatetime *time.Time
	CreatedAt          time.Time
}

func (e *jobStatusEntity) TableName() string {
	return "job_status"
}

func (e *jobStatusEntity) toJobStatus() (service.JobStatus, error) {
	return service.NewJobStatus(domain.JobStatusID(e.ID), domain.JobName(e.JobName), e.JobParameter, e.ExpirationDatetime, e.CreatedAt)
}

type jobStatusRepository struct {
	db *gorm.DB
}

func newJobStatusRepository(ctx context.Context, db *gorm.DB) service.JobStatusRepository {
	return &jobStatusRepository{
		db: db,
	}
}

func (r *jobStatusRepository) AddJobStatus(ctx context.Context, job service.Job) (domain.JobStatusID, error) {
	_, span := tracer.Start(ctx, "jobStatusRepository.AddJobStatus")
	defer span.End()

	// jobStatusID
	id := libD.NewULID()

	// concurrencyKey
	var concurrencyKey *string
	jobName := (string)(job.GetName())
	if !job.IsAllowedConcurrentExecution() {
		concurrencyKey = &jobName
	}

	// expirationDatetime
	expirationDatetime := time.Now().Add(job.GetTimeout())

	entity := jobStatusEntity{
		ID:                 id,
		JobName:            (string)(job.GetName()),
		JobParameter:       job.GetJobParameter(),
		ConcurrencyKey:     concurrencyKey,
		ExpirationDatetime: &expirationDatetime,
	}

	if result := r.db.Create(&entity); result.Error != nil {
		return "", libG.ConvertDuplicatedError(result.Error, service.ErrJobStatusAlreadyExists)
	}

	return domain.JobStatusID(id), nil
}

func (r *jobStatusRepository) RemoveJobStatus(ctx context.Context, jobStatusID domain.JobStatusID) error {
	_, span := tracer.Start(ctx, "jobStatusRepository.RemoveJobStatus")
	defer span.End()

	result := r.db.
		Where("id = ?", (string)(jobStatusID)).
		Delete(&jobStatusEntity{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return service.ErrJobStatusNotFound
		}
		return result.Error
	}
	return nil
}

func (r *jobStatusRepository) RemoveExpiredJobStatus(ctx context.Context) (int, error) {
	_, span := tracer.Start(ctx, "jobStatusRepository.RemoveJobStatus")
	defer span.End()

	result := r.db.
		Where("expiration_datetime <= ?", time.Now()).
		Delete(&jobStatusEntity{})
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (r *jobStatusRepository) FindJobStatusByJobName(ctx context.Context, jobName domain.JobName) ([]service.JobStatus, error) {
	_, span := tracer.Start(ctx, "jobStatusRepository.FindJobStatusByJobName")
	defer span.End()

	var jobStatusEntities []jobStatusEntity
	if result := r.db.
		Where("job_name = ?", (string)(jobName)).
		Order("created_at desc").
		Find(&jobStatusEntities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]service.JobStatus, len(jobStatusEntities))
	for i, e := range jobStatusEntities {
		w, err := e.toJobStatus()
		if err != nil {
			return nil, liberrors.Errorf("toJobStatus. err: %w", err)
		}
		results[i] = w
	}

	return results, nil
}
