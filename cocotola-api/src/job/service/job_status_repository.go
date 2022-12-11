//go:generate mockery --output mock --name JobStatusRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
)

var ErrJobStatusAlreadyExists = errors.New("JobStatus already exists")
var ErrJobStatusNotFound = errors.New("JobStatus not found")

type JobStatusRepository interface {
	AddJobStatus(ctx context.Context, job Job) (domain.JobStatusID, error)
	RemoveJobStatus(ctx context.Context, jobStatusID domain.JobStatusID) error
	RemoveExpiredJobStatus(ctx context.Context) (int, error)
	FindJobStatusByJobName(ctx context.Context, jobName domain.JobName) ([]JobStatus, error)
}
