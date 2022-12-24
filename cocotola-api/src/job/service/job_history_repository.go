//go:generate mockery --output mock --name JobHistoryRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

var ErrJobHistoryNotFound = errors.New("JobHistory not found")

type JobHistory interface {
}

type jobHistory struct {
}

func NewJobHistory() (JobHistory, error) {
	return &jobHistory{}, nil
}

type JobHistoryAddParameter interface {
	GetJobStatusID() domain.JobStatusID
	GetJobName() domain.JobName
	GetJobParamter() string
	GetStatus() string
}

type jobHistoryAddParameter struct {
	JobStatusID  domain.JobStatusID
	JobName      domain.JobName
	JobParameter string
	Status       string
}

func NewJobHistoryAddParameter(jobStatusID domain.JobStatusID, jobName domain.JobName, jobParameter string, status string) (JobHistoryAddParameter, error) {
	m := &jobHistoryAddParameter{
		JobStatusID:  jobStatusID,
		JobName:      jobName,
		JobParameter: jobParameter,
		Status:       status,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *jobHistoryAddParameter) GetJobStatusID() domain.JobStatusID {
	return m.JobStatusID
}

func (m *jobHistoryAddParameter) GetJobName() domain.JobName {
	return m.JobName
}

func (m *jobHistoryAddParameter) GetJobParamter() string {
	return m.JobParameter
}

func (m *jobHistoryAddParameter) GetStatus() string {
	return m.Status
}

type JobHistoryRepository interface {
	AddJobHistory(ctx context.Context, param JobHistoryAddParameter) error
	FindJobHistoryByJobName(ctx context.Context, jobName domain.JobName) (JobHistory, error)
}
