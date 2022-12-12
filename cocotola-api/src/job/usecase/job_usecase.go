package usecase

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	"github.com/kujilabo/cocotola/lib/log"
)

type JobUsecase interface {
	StartJob(ctx context.Context, job service.Job) error
	CleanOldJobs(ctx context.Context) error
}

type jobUsecase struct {
	transaction service.Transaction
}

func NewJobUsecase(ctx context.Context, transaction service.Transaction) (JobUsecase, error) {
	return &jobUsecase{
		transaction: transaction,
	}, nil
}

func (u *jobUsecase) registerStartedRecord(ctx context.Context, job service.Job) (domain.JobStatusID, error) {
	var jobStatusID domain.JobStatusID
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		jobStatusRepo, err := rf.NewJobStatusRepository(ctx)
		if err != nil {
			return err
		}

		jobHistoryRepo, err := rf.NewJobHistoryRepository(ctx)
		if err != nil {
			return err
		}

		tmpJobStatusID, err := jobStatusRepo.AddJobStatus(ctx, job)
		if err != nil {
			return err
		}

		param, err := service.NewJobHistoryAddParameter(tmpJobStatusID, job.GetName(), job.GetJobParameter(), "started")
		if err != nil {
			return err
		}

		if err := jobHistoryRepo.AddJobHistory(ctx, param); err != nil {
			return err
		}

		jobStatusID = tmpJobStatusID
		return nil
	}); err != nil {
		return "", err
	}

	return jobStatusID, nil
}

func (u *jobUsecase) registerStoppedRecord(ctx context.Context, job service.Job, jobStatusID domain.JobStatusID, status string) error {
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		jobStatusRepo, err := rf.NewJobStatusRepository(ctx)
		if err != nil {
			return err
		}

		jobHistoryRepo, err := rf.NewJobHistoryRepository(ctx)
		if err != nil {
			return err
		}

		if err := jobStatusRepo.RemoveJobStatus(ctx, jobStatusID); err != nil {
			return err
		}

		param, err := service.NewJobHistoryAddParameter(jobStatusID, job.GetName(), job.GetJobParameter(), status)
		if err != nil {
			return err
		}

		if err := jobHistoryRepo.AddJobHistory(ctx, param); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *jobUsecase) registerCompletedRecord(ctx context.Context, job service.Job, jobStatusID domain.JobStatusID) error {
	return u.registerStoppedRecord(ctx, job, jobStatusID, "completed")
}

func (u *jobUsecase) registerTimedoutRecord(ctx context.Context, job service.Job, jobStatusID domain.JobStatusID) error {
	return u.registerStoppedRecord(ctx, job, jobStatusID, "timedout")
}

func (u *jobUsecase) registerFailedRecord(ctx context.Context, job service.Job, jobStatusID domain.JobStatusID) error {
	return u.registerStoppedRecord(ctx, job, jobStatusID, "failed")
}

func (u *jobUsecase) StartJob(ctx context.Context, job service.Job) error {
	jobStatusID, err := u.registerStartedRecord(ctx, job)
	if err != nil {
		return err
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
			if err := u.registerTimedoutRecord(bg, job, jobStatusID); err != nil {
				logger.Errorf("registerTimeoutRecord. err: %v", err)
			}
		case err := <-errCh:
			if err != nil {
				logger.Infof("job failed. jobStatusId: %s", jobStatusID)
				if err := u.registerFailedRecord(bg, job, jobStatusID); err != nil {
					logger.Errorf("registerTimeoutRecord. err: %v", err)
				}
			} else {
				logger.Infof("job completed. jobStatusId: %s", jobStatusID)
				if err := u.registerCompletedRecord(bg, job, jobStatusID); err != nil {
					logger.Errorf("registerTimeoutRecord. err: %v", err)
				}
			}
		}
	}()

	return nil
}

func (u *jobUsecase) CleanOldJobs(ctx context.Context) error {

	return nil
}
