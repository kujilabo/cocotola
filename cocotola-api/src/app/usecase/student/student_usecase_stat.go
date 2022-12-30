//go:generate mockery --output mock --name StudentUsecaseStat
package student

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

type StudentUsecaseStat interface {
	GetStat(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Stat, error)
}

type studentUsecaseStat struct {
	transaction     service.Transaction
	pf              service.ProcessorFactory
	findStudentFunc FindStudentFunc
}

func NewStudentUsecaseStat(transaction service.Transaction, pf service.ProcessorFactory,
	findStudentFunc FindStudentFunc) StudentUsecaseStat {
	return &studentUsecaseStat{
		transaction:     transaction,
		pf:              pf,
		findStudentFunc: findStudentFunc,
	}
}

func (s *studentUsecaseStat) GetStat(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Stat, error) {
	_, span := tracer.Start(ctx, "studentUsecaseStat.GetStat")
	defer span.End()

	var result service.Stat
	if err := s.transaction.Do(ctx, func(rf service.RepositoryFactory) error {
		student, err := s.findStudent(ctx, rf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("s.findStudent. err: %w", err)
		}

		stat, err := student.GetStat(ctx)
		if err != nil {
			return liberrors.Errorf("student.GetStat. err: %w", err)
		}

		result = stat
		return nil
	}); err != nil {
		return nil, liberrors.Errorf("GetStatt. err: %w", err)
	}

	return result, nil
}

func (s *studentUsecaseStat) findStudent(ctx context.Context, rf service.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	student, err := s.findStudentFunc(ctx, rf, organizationID, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}

	return student, nil
}
