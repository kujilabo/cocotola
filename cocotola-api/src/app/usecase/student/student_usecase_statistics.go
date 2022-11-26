package student

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/usecase"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type StudentUsecaseStat interface {
	GetStat(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Stat, error)
}

type studentUsecaseStat struct {
	db         *gorm.DB
	pf         service.ProcessorFactory
	rfFunc     service.RepositoryFactoryFunc
	userRfFunc userS.RepositoryFactoryFunc
}

func NewStudentUsecaseStat(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) StudentUsecaseStat {
	return &studentUsecaseStat{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *studentUsecaseStat) GetStat(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Stat, error) {
	var result service.Stat
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, err := s.findStudent(ctx, tx, organizationID, operatorID)
		if err != nil {
			return err
		}

		stat, err := student.GetStat(ctx)
		if err != nil {
			return err
		}

		result = stat
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *studentUsecaseStat) findStudent(ctx context.Context, db *gorm.DB, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	rf, err := s.rfFunc(ctx, db)
	if err != nil {
		return nil, liberrors.Errorf("failed to rfFunc. err: %w", err)
	}
	userRepo, err := s.userRfFunc(ctx, db)
	if err != nil {
		return nil, liberrors.Errorf("failed to userRepo. err: %w", err)
	}
	student, err := usecase.FindStudent(ctx, s.pf, rf, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}

	return student, nil
}
