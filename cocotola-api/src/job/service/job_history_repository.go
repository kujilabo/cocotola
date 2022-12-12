//go:generate mockery --output mock --name JobHistoryRepository
package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type JobHistory interface {
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

	return m, libD.Validator.Struct(m)
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
}
