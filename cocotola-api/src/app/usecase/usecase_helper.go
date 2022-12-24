package usecase

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func FindStudent(ctx context.Context, pf service.ProcessorFactory, rf service.RepositoryFactory, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	userRf, err := rf.NewUserRepositoryFactory(ctx)
	if err != nil {
		return nil, liberrors.Errorf("rf.NewUserRepositoryFactory. err: %w", err)
	}

	systemAdmin, err := userS.NewSystemAdmin(ctx, userRf)
	if err != nil {
		return nil, liberrors.Errorf("userS.NewSystemAdmin. err: %w", err)
	}

	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, liberrors.Errorf("systemAdmin.FindSystemOwnerByOrganizationID. err: %w", err)
	}

	appUser, err := systemOwner.FindAppUserByID(ctx, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("systemOwner.FindAppUserByID. err: %w", err)
	}

	studentModel, err := domain.NewStudentModel(appUser)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewStudentModel. err: %w", err)
	}

	student, err := service.NewStudent(ctx, pf, rf, studentModel)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return student, nil
}
