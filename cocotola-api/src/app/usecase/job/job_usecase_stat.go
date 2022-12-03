package job

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"gorm.io/gorm"

	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

type JobUsecaseStat interface {
	AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin domain.SystemAdminModel, targetDate time.Time) error
}

type jobUsecaseStat struct {
	db         *gorm.DB
	rfFunc     service.RepositoryFactoryFunc
	userRfFunc userS.RepositoryFactoryFunc
}

func NewJobUsecaseStat(db *gorm.DB, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) JobUsecaseStat {
	return &jobUsecaseStat{
		db:         db,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (j *jobUsecaseStat) AggregateStudyResultsOfAllUsers(ctx context.Context, systemAdmin domain.SystemAdminModel, targetDate time.Time) error {
	_, span := tracer.Start(ctx, "jobUsecaseStat.AggregateStudyResultsOfAllUsers")
	defer span.End()

	if err := j.db.Transaction(func(tx *gorm.DB) error {
		rf, err := j.rfFunc(ctx, tx)
		if err != nil {
			return err
		}

		userRf, err := j.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		appUserRepo, err := userRf.NewAppUserRepository()
		if err != nil {
			return err
		}

		systemOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, "cocotola")
		if err != nil {
			return err
		}

		studyStatRepo, err := rf.NewStudyStatRepository(ctx)
		if err != nil {
			return err
		}

		if err := studyStatRepo.AggregateResultsOfAllUsers(ctx, systemOwner, targetDate); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
