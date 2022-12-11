package service

import (
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
)

type JobStatus interface {
}

type jobStatus struct {
	JobName            domain.JobName
	JobParameter       string
	ExpirationDatetime *time.Time
	CreatedAt          time.Time
}

func NewJobStatus(jobName domain.JobName, jobParameter string, expirationDatetime *time.Time, createdAt time.Time) (JobStatus, error) {
	return &jobStatus{
		JobName:            jobName,
		JobParameter:       jobParameter,
		ExpirationDatetime: expirationDatetime,
		CreatedAt:          createdAt,
	}, nil
}
