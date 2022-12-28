//go:generate mockery --output mock --name JobService
package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type JobService interface {
	StartJob(ctx context.Context, job Job) error
	CleanOldJobs(ctx context.Context) error
}

type jobService struct {
	transaction Transaction
}

func NewJobService(ctx context.Context, transaction Transaction) (JobService, error) {
	return &jobService{
		transaction: transaction,
	}, nil
}

func (s *jobService) registerStartedRecord(ctx context.Context, job Job) (domain.JobStatusID, error) {
	var jobStatusID domain.JobStatusID
	if err := s.transaction.Do(ctx, func(rf RepositoryFactory) error {
		jobStatusRepo := rf.NewJobStatusRepository(ctx)
		jobHistoryRepo := rf.NewJobHistoryRepository(ctx)
		tmpJobStatusID, err := jobStatusRepo.AddJobStatus(ctx, job)
		if err != nil {
			return liberrors.Errorf("jobStatusRepo.AddJobStatus. err: %w", err)
		}

		param, err := NewJobHistoryAddParameter(tmpJobStatusID, job.GetName(), job.GetJobParameter(), "started")
		if err != nil {
			return liberrors.Errorf("NewJobHistoryAddParameter. err: %w", err)
		}

		if err := jobHistoryRepo.AddJobHistory(ctx, param); err != nil {
			return liberrors.Errorf("jobHistoryRepo.AddJobHistory. err: %w", err)
		}

		jobStatusID = tmpJobStatusID
		return nil
	}); err != nil {
		return "", liberrors.Errorf("registerStartedRecord. err: %w", err)
	}

	return jobStatusID, nil
}

func (s *jobService) registerStoppedRecord(ctx context.Context, job Job, jobStatusID domain.JobStatusID, status string) error {
	if err := s.transaction.Do(ctx, func(rf RepositoryFactory) error {
		jobStatusRepo := rf.NewJobStatusRepository(ctx)
		jobHistoryRepo := rf.NewJobHistoryRepository(ctx)

		if err := jobStatusRepo.RemoveJobStatus(ctx, jobStatusID); err != nil {
			return liberrors.Errorf("jobStatusRepo.RemoveJobStatus. err: %w", err)
		}

		param, err := NewJobHistoryAddParameter(jobStatusID, job.GetName(), job.GetJobParameter(), status)
		if err != nil {
			return liberrors.Errorf("NewJobHistoryAddParameter. err: %w", err)
		}

		if err := jobHistoryRepo.AddJobHistory(ctx, param); err != nil {
			return liberrors.Errorf("jobHistoryRepo.AddJobHistory. err: %w", err)
		}

		return nil
	}); err != nil {
		return liberrors.Errorf("registerStoppedRecord. err: %w", err)
	}

	return nil
}

func (s *jobService) registerCompletedRecord(ctx context.Context, job Job, jobStatusID domain.JobStatusID) error {
	return s.registerStoppedRecord(ctx, job, jobStatusID, "completed")
}

func (s *jobService) registerTimedoutRecord(ctx context.Context, job Job, jobStatusID domain.JobStatusID) error {
	return s.registerStoppedRecord(ctx, job, jobStatusID, "timedout")
}

func (s *jobService) registerFailedRecord(ctx context.Context, job Job, jobStatusID domain.JobStatusID) error {
	return s.registerStoppedRecord(ctx, job, jobStatusID, "failed")
}

func (s *jobService) StartJob(ctx context.Context, job Job) error {
	jobStatusID, err := s.registerStartedRecord(ctx, job)
	if err != nil {
		return liberrors.Errorf("registerStartedRecord. err: %w", err)
	}

	go func() {
		bg := context.Background()
		logger := log.FromContext(bg)
		timeoutCtx, cancelFunc := context.WithTimeout(bg, job.GetTimeout())
		defer cancelFunc()

		logger.Infof("job started. jobStatusId: %s", jobStatusID)

		errCh := make(chan error)
		go func(ctx context.Context) {
			errCh <- job.Run(ctx)
		}(timeoutCtx)

		select {
		case <-timeoutCtx.Done():
			logger.Infof("job timedout. jobStatusId: %s", jobStatusID)
			if err := s.registerTimedoutRecord(bg, job, jobStatusID); err != nil {
				logger.Errorf("registerTimeoutRecord. err: %v", err)
			}
		case err := <-errCh:
			if err != nil {
				logger.Infof("job failed. jobStatusId: %s", jobStatusID)
				if err := s.registerFailedRecord(bg, job, jobStatusID); err != nil {
					logger.Errorf("registerTimeoutRecord. err: %v", err)
				}
			} else {
				logger.Infof("job completed. jobStatusId: %s", jobStatusID)
				if err := s.registerCompletedRecord(bg, job, jobStatusID); err != nil {
					logger.Errorf("registerTimeoutRecord. err: %v", err)
				}
			}
		}
	}()

	return nil
}

func (s *jobService) CleanOldJobs(ctx context.Context) error {

	return nil
}
