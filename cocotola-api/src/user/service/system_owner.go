package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
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

	orgRepo, err := rf.NewOrganizationRepository(ctx)
	if err != nil {
		return nil, err
	}
	appUserRepo, err := rf.NewAppUserRepository(ctx)
	if err != nil {
		return nil, err
	}
	spaceRepo, err := rf.NewSpaceRepository(ctx)
	if err != nil {
		return nil, err
	}
	appUserGroupRepo, err := rf.NewAppUserGroupRepository(ctx)
	if err != nil {
		return nil, err
	}
	groupUserRepo, err := rf.NewGroupUserRepository(ctx)
	if err != nil {
		return nil, err
	}
	rbacRepo, err := rf.NewRBACRepository(ctx)
	if err != nil {
		return nil, err
	}

	m := &systemOwner{
		SystemOwnerModel: systemOwnerModel,
		orgRepo:          orgRepo,
		spaceRepo:        spaceRepo,
		appUserGroupRepo: appUserGroupRepo,
		appUserRepo:      appUserRepo,
		groupUserRepo:    groupUserRepo,
		rbacRepo:         rbacRepo,
	}

	return m, libD.Validator.Struct(m)
}

func (m *systemOwner) GetOrganization(ctx context.Context) (Organization, error) {
	return m.orgRepo.GetOrganization(ctx, m)
}

func (m *systemOwner) FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error) {
	return m.appUserRepo.FindAppUserByID(ctx, m, id)
}

func (m *systemOwner) FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error) {
	return m.appUserRepo.FindAppUserByLoginID(ctx, m, loginID)
}

func (m *systemOwner) FindSystemSpace(ctx context.Context) (Space, error) {
	return m.spaceRepo.FindSystemSpace(ctx, m)
}

func (m *systemOwner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddStudent")
	appUserID, err := m.appUserRepo.AddAppUser(ctx, m, param)
	if err != nil {
		return 0, err
	}
	appUser, err := m.appUserRepo.FindAppUserByID(ctx, m, appUserID)
	if err != nil {
		return 0, err
	}

	// personalGroupID, err := m.rf.NewAppUserGroupRepository().AddPersonalGroup(m, studentID)
	// if err != nil {
	// 	return 0, err
	// }

	publicGroup, err := m.appUserGroupRepo.FindPublicGroup(ctx, m)
	if err != nil {
		return 0, err
	}
	if err := m.groupUserRepo.AddGroupUser(ctx, m, publicGroup.GetAppUerGroupID(), appUser.GetAppUserID()); err != nil {
		return 0, err
	}

	spaceID, err := m.spaceRepo.AddPersonalSpace(ctx, m, appUser)
	if err != nil {
		return 0, err
	}

	logger.Infof("Personal spaceID: %d", spaceID)

	spaceWriter := domain.NewSpaceWriterRole(spaceID)
	spaceObject := domain.NewSpaceObject(spaceID)
	userSubject := domain.NewUserObject(appUserID)

	if err := m.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "read"); err != nil {
		return 0, err
	}

	if err := m.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "write"); err != nil {
		return 0, err
	}

	if err := m.rbacRepo.AddNamedGroupingPolicy(userSubject, spaceWriter); err != nil {
		return 0, err
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
		return 0, err
	}
	return spaceID, nil
}
