//go:generate mockery --output mock --name StudentUsecaseStat
package student

import (
	"context"

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
	transaction service.Transaction
	pf          service.ProcessorFactory
}

func NewStudentUsecaseStat(transaction service.Transaction, pf service.ProcessorFactory) StudentUsecaseStat {
	return &studentUsecaseStat{
		transaction: transaction,
		pf:          pf,
	}
}

func (s *studentUsecaseStat) GetStat(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Stat, error) {
	_, span := tracer.Start(ctx, "studentUsecaseStat.GetStat")
	defer span.End()

	var result service.Stat
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory, userRf userS.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, userRf, organizationID, operatorID)
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

func (s *studentUsecaseStat) findStudent(ctx context.Context, rf service.RepositoryFactory, userRf userS.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}

	return student, nil
}
