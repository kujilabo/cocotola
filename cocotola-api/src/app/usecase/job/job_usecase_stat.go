package job

import (
	"context"
	"errors"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/lib/log"

	jobD "github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

const studyStatsRetentionDays = 14
const userFetchSize = 10

type JobUsecaseStat interface {
	AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin domain.SystemAdminModel) error
	CleanStudyStats(ctx context.Context, systemAdmin domain.SystemAdminModel, expirationDate time.Time) error
}

type jobUsecaseStat struct {
	transaction service.Transaction
	jobService  jobS.JobService
}

func NewJobUsecaseStat(transaction service.Transaction, jobService jobS.JobService) JobUsecaseStat {
	return &jobUsecaseStat{
		transaction: transaction,
		jobService:  jobService,
	}
}

func (u *jobUsecaseStat) getSystemOwnerAndRepositoryFactory(ctx context.Context, systemAdmin domain.SystemAdminModel, rf service.RepositoryFactory) (userS.SystemOwner, userS.RepositoryFactory, service.RepositoryFactory, error) {

	var repositoryFactory service.RepositoryFactory
	var userRepositoryFactory userS.RepositoryFactory
	var systemOwner userS.SystemOwner
	if err := func() error {
		userRf, err := rf.NewUserRepositoryFactory(ctx)
		if err != nil {
			return err
		}

		appUserRepo, err := userRf.NewAppUserRepository(ctx)
		if err != nil {
			return err
		}

		so, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, "cocotola")
		if err != nil {
			return err
		}

		systemOwner = so
		repositoryFactory = rf
		userRepositoryFactory = userRf
		return nil
	}(); err != nil {
		return nil, nil, nil, err
	}

	return systemOwner, userRepositoryFactory, repositoryFactory, nil
}

func (u *jobUsecaseStat) AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin domain.SystemAdminModel) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.AggregateStudyResultsOfAllUsers")
	defer span.End()
	logger := log.FromContext(ctx)

	// Check if the process has already been executed.

	makeJobName := func(targetDate time.Time) jobD.JobName {
		return jobD.JobName("AggregateStudyResults-" + targetDate.Format("2006-01-02"))
	}

	targetDateList := make([]time.Time, 0)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		jobRf, err := rf.NewJobRepositoryFactory(ctx)
		if err != nil {
			return err
		}

		jobHistoryRepo, err := jobRf.NewJobHistoryRepository(ctx)
		if err != nil {
			return err
		}

		for i := 1; i <= 7; i++ {
			date := today.AddDate(0, 0, -1*i)
			jobName := makeJobName(date)
			_, err := jobHistoryRepo.FindJobHistoryByJobName(ctx, jobName)
			if err != nil && errors.Is(err, jobS.ErrJobHistoryNotFound) {
				targetDateList = append(targetDateList, date)
			} else if err != nil {
				return err
			} else {
				logger.Debugf("%s has already been completed", jobName)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	makeFunc := func(targetDate time.Time) func(context.Context) error {
		return func(ctx context.Context) error {
			pageNo := 1
			pageSize := userFetchSize
			if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
				systemOwner, userRf, rf, err := u.getSystemOwnerAndRepositoryFactory(ctx, systemAdmin, rf)
				if err != nil {
					return err
				}

				studyStatRepo, err := rf.NewStudyStatRepository(ctx)
				if err != nil {
					return err
				}

				userRepo, err := userRf.NewAppUserRepository(ctx)
				if err != nil {
					return err
				}

				for {
					userIDs, err := userRepo.FindAppUserIDs(ctx, systemOwner, pageNo, pageSize)
					if err != nil {
						return err
					}

					if len(userIDs) == 0 {
						break
					}

					for _, userID := range userIDs {
						if err := studyStatRepo.AggregateResults(ctx, systemOwner, targetDate, userID); err != nil {
							return err
						}
					}

					pageNo++
				}

				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	}

	for _, targetDate := range targetDateList {
		jobName := makeJobName(targetDate)
		fn := makeFunc(targetDate)
		job, err := jobS.NewJob(jobName, time.Minute, false, fn)
		if err != nil {
			return err
		}
		if err := u.jobService.StartJob(ctx, job); err != nil {
			return err
		}
	}

	return nil
}

func (j *jobUsecaseStat) CleanStudyStats(ctx context.Context, systemAdmin domain.SystemAdminModel, expirationDate time.Time) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.CleanStudyStats")
	defer span.End()

	if err := j.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		systemOwner, _, rf, err := j.getSystemOwnerAndRepositoryFactory(ctx, systemAdmin, rf)
		if err != nil {
			return err
		}

		studyStatRepo, err := rf.NewStudyStatRepository(ctx)
		if err != nil {
			return err
		}

		expirationDate := time.Now().AddDate(0, 0, -1*studyStatsRetentionDays)

		if err := studyStatRepo.CleanStudyStats(ctx, systemOwner, expirationDate); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
