package service

import (
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
)

type JobStatus interface {
	GetJobStatusID() domain.JobStatusID
	GetJobName() domain.JobName
	GetJobParameter() string
}

type jobStatus struct {
	JobStatusID        domain.JobStatusID
	JobName            domain.JobName
	JobParameter       string
	ExpirationDatetime *time.Time
	CreatedAt          time.Time
}

func NewJobStatus(jobStatusID domain.JobStatusID, jobName domain.JobName, jobParameter string, expirationDatetime *time.Time, createdAt time.Time) (JobStatus, error) {
	return &jobStatus{
		JobStatusID:        jobStatusID,
		JobName:            jobName,
		JobParameter:       jobParameter,
		ExpirationDatetime: expirationDatetime,
		CreatedAt:          createdAt,
	}, nil
}

func (m *jobStatus) GetJobStatusID() domain.JobStatusID {
	return m.JobStatusID
}

func (m *jobStatus) GetJobName() domain.JobName {
	return m.JobName
}

func (m *jobStatus) GetJobParameter() string {
	return m.JobParameter
}
