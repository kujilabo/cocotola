package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

var rfFunc func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)

func InitSystemAdmin(rfFuncArg func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)) {
	if rfFuncArg == nil {
		panic(errors.New("invalid argment"))
	}
	rfFunc = rfFuncArg
}

type SystemAdmin interface {
	domain.SystemAdminModel

	FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error)

	FindOrganizationByName(ctx context.Context, name string) (Organization, error)

	AddOrganization(ctx context.Context, parma OrganizationAddParameter) (domain.OrganizationID, error)
}

type systemAdmin struct {
	domain.SystemAdminModel
	rf RepositoryFactory
}

func NewSystemAdmin(rf RepositoryFactory) SystemAdmin {
	return &systemAdmin{
		SystemAdminModel: domain.NewSystemAdminModel(),
		rf:               rf,
	}
}
func NewSystemAdminFromDB(ctx context.Context, db *gorm.DB) (SystemAdmin, error) {
	rf, err := rfFunc(ctx, db)
	if err != nil {
		return nil, err
	}
	return NewSystemAdmin(rf), nil
}

func (s *systemAdmin) FindSystemOwnerByOrganizationID(ctx context.Context, organizationID domain.OrganizationID) (SystemOwner, error) {
	appUserRepo, err := s.rf.NewAppUserRepository()
	if err != nil {
		return nil, err
	}
	return appUserRepo.FindSystemOwnerByOrganizationID(ctx, s, organizationID)
}

func (s *systemAdmin) FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error) {
	appUserRepo, err := s.rf.NewAppUserRepository()
	if err != nil {
		return nil, err
	}
	return appUserRepo.FindSystemOwnerByOrganizationName(ctx, s, organizationName)
}

func (s *systemAdmin) FindOrganizationByName(ctx context.Context, name string) (Organization, error) {
	orgRepo, err := s.rf.NewOrganizationRepository()
	if err != nil {
		return nil, err
	}
	return orgRepo.FindOrganizationByName(ctx, s, name)
}

func (s *systemAdmin) AddOrganization(ctx context.Context, param OrganizationAddParameter) (domain.OrganizationID, error) {
	logger := log.FromContext(ctx)
	orgRepo, err := s.rf.NewOrganizationRepository()
	if err != nil {
		return 0, err
	}
	// add organization
	organizationID, err := orgRepo.AddOrganization(ctx, s, param)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddOrganization. error: %w", err)
	}
	appUserRepo, err := s.rf.NewAppUserRepository()
	if err != nil {
		return 0, err
	}

	// add system owner
	systemOwnerID, err := appUserRepo.AddSystemOwner(ctx, s, organizationID)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, s, param.GetName())
	if err != nil {
		return 0, liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	// // add system student
	// systemStudentID, err := s.rf.NewAppUserRepository().AddSystemStudent(ctx, systemOwner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddSystemStudent. error: %w", err)
	// }

	// add owner
	ownerID, err := appUserRepo.AddFirstOwner(ctx, systemOwner, param.GetFirstOwner())
	if err != nil {
		return 0, liberrors.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	owner, err := appUserRepo.FindOwnerByLoginID(ctx, systemOwner, param.GetFirstOwner().GetLoginID())
	if err != nil {
		return 0, liberrors.Errorf("failed to FindOwnerByLoginID. error: %w", err)
	}

	appUserGroupRepo, err := s.rf.NewAppUserGroupRepository()
	if err != nil {
		return 0, err
	}

	// add public group
	publicGroupID, err := appUserGroupRepo.AddPublicGroup(ctx, systemOwner)
	if err != nil {
		return 0, liberrors.Errorf("failed to AddPublicGroup. error: %w", err)
	}

	groupUserRepo, err := s.rf.NewGroupUserRepository()
	if err != nil {
		return 0, err
	}
	// public-group <-> owner
	if err := groupUserRepo.AddGroupUser(ctx, systemOwner, publicGroupID, ownerID); err != nil {
		return 0, liberrors.Errorf("failed to AddGroupUser. error: %w", err)
	}

	spaceRepo, err := s.rf.NewSpaceRepository()
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
