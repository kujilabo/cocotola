package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	"github.com/kujilabo/cocotola/lib/log"
)

type SystemOwner interface {
	domain.SystemOwnerModel

	GetOrganization(ctxc context.Context) (Organization, error)

	FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error)

	FindSystemSpace(ctx context.Context) (Space, error)

	AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error)

	AddSystemSpace(ctx context.Context) (domain.SpaceID, error)
}

type systemOwner struct {
	domain.SystemOwnerModel
	orgRepo          OrganizationRepository
	spaceRepo        SpaceRepository
	appUserGroupRepo AppUserGroupRepository
	appUserRepo      AppUserRepository
	groupUserRepo    GroupUserRepository
	rbacRepo         RBACRepository
}

func NewSystemOwner(ctx context.Context, rf RepositoryFactory, systemOwnerModel domain.SystemOwnerModel) (SystemOwner, error) {

	orgRepo := rf.NewOrganizationRepository(ctx)
	appUserRepo := rf.NewAppUserRepository(ctx)
	spaceRepo := rf.NewSpaceRepository(ctx)
	appUserGroupRepo := rf.NewAppUserGroupRepository(ctx)
	groupUserRepo := rf.NewGroupUserRepository(ctx)
	rbacRepo := rf.NewRBACRepository(ctx)

	m := &systemOwner{
		SystemOwnerModel: systemOwnerModel,
		orgRepo:          orgRepo,
		spaceRepo:        spaceRepo,
		appUserGroupRepo: appUserGroupRepo,
		appUserRepo:      appUserRepo,
		groupUserRepo:    groupUserRepo,
		rbacRepo:         rbacRepo,
	}

	if err := libD.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libD.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *systemOwner) GetOrganization(ctx context.Context) (Organization, error) {
	org, err := m.orgRepo.GetOrganization(ctx, m)
	if err != nil {
		return nil, liberrors.Errorf("m.orgRepo.GetOrganization. err: %w", err)
	}

	return org, nil

}

func (m *systemOwner) FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error) {
	appUser, err := m.appUserRepo.FindAppUserByID(ctx, m, id)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindAppUserByID. err: %w", err)
	}

	return appUser, nil
}

func (m *systemOwner) FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error) {
	appUser, err := m.appUserRepo.FindAppUserByLoginID(ctx, m, loginID)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindAppUserByLoginID. err: %w", err)
	}

	return appUser, nil
}

func (m *systemOwner) FindSystemSpace(ctx context.Context) (Space, error) {
	space, err := m.spaceRepo.FindSystemSpace(ctx, m)
	if err != nil {
		return nil, liberrors.Errorf("m.appUserRepo.FindSystemSpace. err: %w", err)
	}

	return space, nil
}

func (m *systemOwner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddStudent")
	appUserID, err := m.appUserRepo.AddAppUser(ctx, m, param)
	if err != nil {
		return 0, liberrors.Errorf("m.appUserRepo.AddAppUser. err: %w", err)
	}
	appUser, err := m.appUserRepo.FindAppUserByID(ctx, m, appUserID)
	if err != nil {
		return 0, liberrors.Errorf("m.appUserRepo.FindAppUserByID. err: %w", err)
	}

	// personalGroupID, err := m.rf.NewAppUserGroupRepository().AddPersonalGroup(m, studentID)
	// if err != nil {
	// 	return 0, err
	// }

	publicGroup, err := m.appUserGroupRepo.FindPublicGroup(ctx, m)
	if err != nil {
		return 0, liberrors.Errorf("m.appUserGroupRepo.FindPublicGroup. err: %w", err)
	}
	if err := m.groupUserRepo.AddGroupUser(ctx, m, publicGroup.GetAppUerGroupID(), appUser.GetAppUserID()); err != nil {
		return 0, liberrors.Errorf("m.groupUserRepo.AddGroupUser. err: %w", err)
	}

	spaceID, err := m.spaceRepo.AddPersonalSpace(ctx, m, appUser)
	if err != nil {
		return 0, liberrors.Errorf("m.spaceRepo.AddPersonalSpace. err: %w", err)
	}

	logger.Infof("Personal spaceID: %d", spaceID)

	spaceWriter := domain.NewSpaceWriterRole(spaceID)
	spaceObject := domain.NewSpaceObject(spaceID)
	userSubject := domain.NewUserObject(appUserID)

	if err := m.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "read"); err != nil {
		return 0, liberrors.Errorf("problemRepo.AddNamedPolicy(read). err: %w", err)
	}

	if err := m.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "write"); err != nil {
		return 0, liberrors.Errorf("problemRepo.AddNamedPolicy(write). err: %w", err)
	}

	if err := m.rbacRepo.AddNamedGroupingPolicy(userSubject, spaceWriter); err != nil {
		return 0, liberrors.Errorf("problemRepo.AddNamedGroupingPolicy. err: %w", err)
	}

	// defaultSpace, err := m.rf.NewSpaceRepository().FindDefaultSpace(ctx, s)
	// if err != nil {
	// 	return 0, err
	// }

	// if err := m.rf.NewUserSpaceRepository().Add(ctx, appUser, SpaceID(defaultSpace.GetID())); err != nil {
	// 	return 0, err
	// }

	return appUserID, nil
}

func (m *systemOwner) AddSystemSpace(ctx context.Context) (domain.SpaceID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddSystemSpace")

	spaceID, err := m.spaceRepo.AddSystemSpace(ctx, m)
	if err != nil {
		return 0, liberrors.Errorf("m.spaceRepo.AddSystemSpace. err: %w", err)
	}
	return spaceID, nil
}
