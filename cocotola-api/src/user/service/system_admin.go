package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type SystemAdmin interface {
	domain.SystemAdminModel

	FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error)

	FindOrganizationByName(ctx context.Context, name string) (Organization, error)

	AddOrganization(ctx context.Context, parma OrganizationAddParameter) (domain.OrganizationID, error)
}

type systemAdmin struct {
	domain.SystemAdminModel
	rf          RepositoryFactory
	orgRepo     OrganizationRepository
	appUserRepo AppUserRepository
}

func NewSystemAdmin(ctx context.Context, rf RepositoryFactory) (SystemAdmin, error) {
	orgRepo := rf.NewOrganizationRepository(ctx)
	appUserRepo := rf.NewAppUserRepository(ctx)

	return &systemAdmin{
		SystemAdminModel: domain.NewSystemAdminModel(),
		rf:               rf,
		orgRepo:          orgRepo,
		appUserRepo:      appUserRepo,
	}, nil
}

func (m *systemAdmin) FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error) {
	sysOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationID(ctx, m, organizationID)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindSystemOwnerByOrganizationID. error: %w", err)
	}

	return sysOwner, nil
}

func (m *systemAdmin) FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error) {
	sysOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, organizationName)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindSystemOwnerByOrganizationName. error: %w", err)
	}

	return sysOwner, nil
}

func (m *systemAdmin) FindOrganizationByName(ctx context.Context, name string) (Organization, error) {
	org, err := m.orgRepo.FindOrganizationByName(ctx, m, name)
	if err != nil {
		return nil, liberrors.Errorf("m.orgRepo.FindOrganizationByName. error: %w", err)
	}

	return org, nil
}

func (m *systemAdmin) AddOrganization(ctx context.Context, param OrganizationAddParameter) (domain.OrganizationID, error) {
	logger := log.FromContext(ctx)

	// add organization
	organizationID, err := m.orgRepo.AddOrganization(ctx, m, param)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddOrganization. error: %w", err)
	}

	// add system owner
	systemOwnerID, err := m.appUserRepo.AddSystemOwner(ctx, m, organizationID)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, param.GetName())
	if err != nil {
		return 0, liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	// // add system student
	// systemStudentID, err := s.rf.NewAppUserRepository().AddSystemStudent(ctx, systemOwner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddSystemStudent. error: %w", err)
	// }

	// add owner
	ownerID, err := m.appUserRepo.AddFirstOwner(ctx, systemOwner, param.GetFirstOwner())
	if err != nil {
		return 0, liberrors.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	owner, err := m.appUserRepo.FindOwnerByLoginID(ctx, systemOwner, param.GetFirstOwner().GetLoginID())
	if err != nil {
		return 0, liberrors.Errorf("failed to FindOwnerByLoginID. error: %w", err)
	}

	appUserGroupRepo := m.rf.NewAppUserGroupRepository(ctx)

	// add public group
	publicGroupID, err := appUserGroupRepo.AddPublicGroup(ctx, systemOwner)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddPublicGroup. error: %w", err)
	}

	groupUserRepo := m.rf.NewGroupUserRepository(ctx)

	// public-group <-> owner
	if err := groupUserRepo.AddGroupUser(ctx, systemOwner, publicGroupID, ownerID); err != nil {
		return 0, liberrors.Errorf("failed to AddGroupUser. error: %w", err)
	}

	spaceRepo := m.rf.NewSpaceRepository(ctx)

	// add default space
	spaceID, err := spaceRepo.AddDefaultSpace(ctx, systemOwner)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddDefaultSpace. error: %w", err)
	}

	logger.Infof("SystemOwnerID:%d, owner: %+v, spaceID: %d", systemOwnerID, owner, spaceID)
	// logger.Infof("SystemOwnerID:%d, SystemStudentID:%d, owner: %+v, spaceID: %d", systemOwnerID, systemStudentID, owner, spaceID)

	// // add personal group
	// personalGroupID, err := s.appUserGroupRepositor.AddPublicGroup(owner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddPersonalGroup. error: %w", err)
	// }

	// // personal-group <-> owner
	// if err := s.groupUserRepository.AddGroupUser(systemOwner, personalGroupID, ownerID); err != nil {
	// 	return 0, fmt.Errorf("failed to AddGroupUser. error: %w", err)
	// }

	return organizationID, nil
}
