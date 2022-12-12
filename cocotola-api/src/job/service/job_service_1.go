package service

// import (
// 	"context"
// 	"errors"

// 	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
// )

// var ErrJobAlreadyExists = errors.New("Job already exists")

// type JobService interface {
// 	RegisterJob(ctx context.Context, job Job) error
// 	Run(ctx context.Context, jobName domain.JobName) (domain.JobStatusID, error)
// 	GetJobStatus(ctx context.Context, jobID domain.JobStatusID) (JobStatus, error)
// }

// type jobService struct {
// 	jobList map[domain.JobName]Job
// }

// func NewJobService() JobService {
// 	return &jobService{
// 		jobList: map[domain.JobName]Job{},
// 	}
// }

// func (s *jobService) RegisterJob(ctx context.Context, job Job) error {
// 	if _, ok := s.jobList[job.GetName()]; ok {
// 		return ErrJobAlreadyExists
// 	}

// 	s.jobList[job.GetName()] = job

// 	return nil
// }

// func (s *jobService) Run(ctx context.Context, jobName domain.JobName) (domain.JobStatusID, error) {
// 	return "", nil
// }

// func (s *jobService) GetJobStatus(ctx context.Context, jobID domain.JobStatusID) (JobStatus, error) {
// 	return nil, nil
// }
