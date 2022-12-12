package job

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"

	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

const studyStatsRetentionDays = 14
const userFetchSize = 10

type JobUsecaseStat interface {
	AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin domain.SystemAdminModel, targetDate time.Time) error
	CleanStudyStats(ctx context.Context, systemAdmin domain.SystemAdminModel, expirationDate time.Time) error
}

type jobUsecaseStat struct {
	transaction service.Transaction
}

func NewJobUsecaseStat(transaction service.Transaction) JobUsecaseStat {
	return &jobUsecaseStat{
		transaction: transaction,
	}
}

func (j *jobUsecaseStat) getSystemOwnerAndRepositoryFactory(ctx context.Context, systemAdmin domain.SystemAdminModel, rf service.RepositoryFactory, userRf userS.RepositoryFactory) (userS.SystemOwner, userS.RepositoryFactory, service.RepositoryFactory, error) {

	var repositoryFactory service.RepositoryFactory
	var userRepositoryFactory userS.RepositoryFactory
	var systemOwner userS.SystemOwner
	if err := func() error {
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

func (j *jobUsecaseStat) AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin domain.SystemAdminModel, targetDate time.Time) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.AggregateStudyResultsOfAllUsers")
	defer span.End()

	pageNo := 1
	pageSize := userFetchSize
	if err := j.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		systemOwner, userRf, rf, err := j.getSystemOwnerAndRepositoryFactory(ctx, systemAdmin, rf, userRf)
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

func (j *jobUsecaseStat) CleanStudyStats(ctx context.Context, systemAdmin domain.SystemAdminModel, expirationDate time.Time) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.CleanStudyStats")
	defer span.End()

	if err := j.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		systemOwner, _, rf, err := j.getSystemOwnerAndRepositoryFactory(ctx, systemAdmin, rf, userRf)
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
