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
	orgRepo, err := rf.NewOrganizationRepository(ctx)
	if err != nil {
		return nil, err
	}
	appUserRepo, err := rf.NewAppUserRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &systemAdmin{
		SystemAdminModel: domain.NewSystemAdminModel(),
		rf:               rf,
		orgRepo:          orgRepo,
		appUserRepo:      appUserRepo,
	}, nil
}

func (m *systemAdmin) FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error) {
	return m.appUserRepo.FindSystemOwnerByOrganizationID(ctx, m, organizationID)
}

func (m *systemAdmin) FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error) {
	return m.appUserRepo.FindSystemOwnerByOrganizationName(ctx, m, organizationName)
}

func (m *systemAdmin) FindOrganizationByName(ctx context.Context, name string) (Organization, error) {
	return m.orgRepo.FindOrganizationByName(ctx, m, name)
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

	appUserGroupRepo, err := m.rf.NewAppUserGroupRepository(ctx)
	if err != nil {
		return 0, err
	}

	// add public group
	publicGroupID, err := appUserGroupRepo.AddPublicGroup(ctx, systemOwner)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddPublicGroup. error: %w", err)
	}

	groupUserRepo, err := m.rf.NewGroupUserRepository(ctx)
	if err != nil {
		return 0, err
	}
	// public-group <-> owner
	if err := groupUserRepo.AddGroupUser(ctx, systemOwner, publicGroupID, ownerID); err != nil {
		return 0, liberrors.Errorf("failed to AddGroupUser. error: %w", err)
	}

	spaceRepo, err := m.rf.NewSpaceRepository(ctx)
	if err != nil {
		return 0, err
	}
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
