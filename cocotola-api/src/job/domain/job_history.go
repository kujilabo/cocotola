package domain

import "time"

type JobHistory interface {
}

type jobHistory struct {
	JobStatusID  JobStatusID
	JobName      JobName
	JobParameter string
	Status       string
	CreatedAt    time.Time
}

func NewJobHistory(jobStatusID JobStatusID, jobName JobName, jobParameter string, status string, createdAt time.Time) (JobHistory, error) {
	return &jobHistory{
		JobStatusID:  jobStatusID,
		JobName:      jobName,
		JobParameter: jobParameter,
		Status:       status,
		CreatedAt:    createdAt,
	}, nil
}

func (m *jobHistory) GetJobStatusID() JobStatusID {
	return m.JobStatusID
}

func (m *jobHistory) GetJobName() JobName {
	return m.JobName
}

func (m *jobHistory) GetJobParameter() string {
	return m.JobParameter
}

func (m *jobHistory) GetStatus() string {
	return m.Status
}

func (m *jobHistory) GteCreatedAt() time.Time {
	return m.CreatedAt
}
