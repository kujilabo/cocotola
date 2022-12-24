package job

import (
	"context"
	"errors"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	jobD "github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

const studyStatsRetentionDays = 14
const userFetchSize = 10

type JobUsecaseStat interface {
	AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin userD.SystemAdminModel) error
	CleanStudyStats(ctx context.Context, systemAdmin userD.SystemAdminModel, expirationDate time.Time) error
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

func (u *jobUsecaseStat) getSystemOwnerAndRepositoryFactory(ctx context.Context, systemAdmin userD.SystemAdminModel, rf service.RepositoryFactory) (userS.SystemOwner, userS.RepositoryFactory, service.RepositoryFactory, error) {
	var repositoryFactory service.RepositoryFactory
	var userRepositoryFactory userS.RepositoryFactory
	var systemOwner userS.SystemOwner
	if err := func() error {
		userRf, err := rf.NewUserRepositoryFactory(ctx)
		if err != nil {
			return liberrors.Errorf("rf.NewUserRepositoryFactory. err: %w", err)
		}

		appUserRepo := userRf.NewAppUserRepository(ctx)

		so, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, service.OrganizationName)
		if err != nil {
			return liberrors.Errorf("appUserRepo.FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemOwner = so
		repositoryFactory = rf
		userRepositoryFactory = userRf
		return nil
	}(); err != nil {
		return nil, nil, nil, liberrors.Errorf("getSystemOwnerAndRepositoryFactory. err: %w", err)
	}

	return systemOwner, userRepositoryFactory, repositoryFactory, nil
}

func (u *jobUsecaseStat) getTargetDateList(ctx context.Context, makeJobName func(targetDate time.Time) jobD.JobName) ([]time.Time, error) {
	logger := log.FromContext(ctx)
	targetDateList := make([]time.Time, 0)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		jobRf, err := rf.NewJobRepositoryFactory(ctx)
		if err != nil {
			return liberrors.Errorf("rf.NewJobRepositoryFactory. err: %w", err)
		}

		jobHistoryRepo := jobRf.NewJobHistoryRepository(ctx)

		for i := 1; i <= 7; i++ {
			date := today.AddDate(0, 0, -1*i)
			jobName := makeJobName(date)
			_, err := jobHistoryRepo.FindJobHistoryByJobName(ctx, jobName)
			if err != nil && errors.Is(err, jobS.ErrJobHistoryNotFound) {
				targetDateList = append(targetDateList, date)
			} else if err != nil {
				return liberrors.Errorf("jobHistoryRepo.FindJobHistoryByJobName. err: %w", err)
			} else {
				logger.Debugf("%s has already been completed", jobName)
			}
		}
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("getTargetDateList. err: %w", err)
	}

	return targetDateList, nil
}

func (u *jobUsecaseStat) aggregateStudyResultsOfAllUsersByDate(ctx context.Context, systemAdmin userD.SystemAdminModel, targetDate time.Time) error {
	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		systemOwner, userRf, rf, err := u.getSystemOwnerAndRepositoryFactory(ctx, systemAdmin, rf)
		if err != nil {
			return liberrors.Errorf("u.getSystemOwnerAndRepositoryFactory. err: %w", err)
		}

		studyStatRepo := rf.NewStudyStatRepository(ctx)
		userRepo := userRf.NewAppUserRepository(ctx)

		aggregateResultseFunc := func(ctx context.Context, userID userD.AppUserID) error {
			if err := studyStatRepo.AggregateResults(ctx, systemOwner, targetDate, userID); err != nil {
				return liberrors.Errorf("studyStatRepo.AggregateResults. err: %w", err)
			}
			return nil
		}

		if err := u.doSomethingForAllUsersWithAppUserID(ctx, userRepo, systemOwner, aggregateResultseFunc); err != nil {
			return liberrors.Errorf("u.doSomethingForAllUsersWithAppUserID. err: %w", err)
		}

		return nil
	}); err != nil {
		return liberrors.Errorf("AggregateStudyResultsOfAllUsers. err: %w", err)
	}

	return nil
}

func (u *jobUsecaseStat) AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin userD.SystemAdminModel) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.AggregateStudyResultsOfAllUsers")
	defer span.End()

	// Check if the process has already been executed.

	makeJobName := func(targetDate time.Time) jobD.JobName {
		return jobD.JobName("AggregateStudyResults-" + targetDate.Format("2006-01-02"))
	}

	makeJobFunc := func(targetDate time.Time) jobS.JobFunc {
		return func(ctx context.Context) error {
			return u.aggregateStudyResultsOfAllUsersByDate(ctx, systemAdmin, targetDate)
		}
	}

	targetDateList, err := u.getTargetDateList(ctx, makeJobName)
	if err != nil {
		return err
	}

	for _, targetDate := range targetDateList {
		jobName := makeJobName(targetDate)
		fn := makeJobFunc(targetDate)
		job, err := jobS.NewJob(jobName, time.Minute, false, fn)
		if err != nil {
			return liberrors.Errorf("jobS.NewJob. err: %w", err)
		}
		if err := u.jobService.StartJob(ctx, job); err != nil {
			return liberrors.Errorf("u.jobService.StartJob. err: %w", err)
		}
	}

	return nil
}

func (u *jobUsecaseStat) doSomethingForAllUsersWithAppUserID(ctx context.Context, userRepo userS.AppUserRepository, systemOwner userD.SystemOwnerModel, fn func(ctx context.Context, userID userD.AppUserID) error) error {
	pageNo := 1
	pageSize := userFetchSize
	for {
		userIDs, err := userRepo.FindAppUserIDs(ctx, systemOwner, pageNo, pageSize)
		if err != nil {
			return liberrors.Errorf("userRepo.FindAppUserIDs. err: %w", err)
		}

		if len(userIDs) == 0 {
			break
		}

		for _, userID := range userIDs {
			if err := fn(ctx, userID); err != nil {
				return liberrors.Errorf("fn. err: %w", err)
			}
		}

		pageNo++
	}

	return nil
}

func (u *jobUsecaseStat) CleanStudyStats(ctx context.Context, systemAdmin userD.SystemAdminModel, expirationDate time.Time) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.CleanStudyStats")
	defer span.End()

	if err := u.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		systemOwner, _, rf, err := u.getSystemOwnerAndRepositoryFactory(ctx, systemAdmin, rf)
		if err != nil {
			return liberrors.Errorf("u.getSystemOwnerAndRepositoryFactory. err: %w", err)
		}

		studyStatRepo := rf.NewStudyStatRepository(ctx)

		expirationDate := time.Now().AddDate(0, 0, -1*studyStatsRetentionDays)

		if err := studyStatRepo.CleanStudyStats(ctx, systemOwner, expirationDate); err != nil {
			return liberrors.Errorf("studyStatRepo.CleanStudyStats. err: %w", err)
		}

		return nil
	}); err != nil {
		return liberrors.Errorf("CleanStudyStats. err: %w", err)
	}

	return nil
}
