package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	libD "github.com/kujilabo/cocotola/lib/domain"
	"github.com/kujilabo/cocotola/lib/log"
)

const SystemOwnerID = 2

type SystemOwner interface {
	// AppUser
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
	rf               RepositoryFactory
	orgRepo          OrganizationRepository
	spaceRepo        SpaceRepository
	appUserGroupRepo AppUserGroupRepository
	appUserRepo      AppUserRepository
	groupUserRepo    GroupUserRepository
	rbacRepo         RBACRepository
}

func NewSystemOwner(rf RepositoryFactory, systemOwnerModel domain.SystemOwnerModel) (SystemOwner, error) {

	orgRepo, err := rf.NewOrganizationRepository()
	if err != nil {
		return nil, err
	}
	appUserRepo, err := rf.NewAppUserRepository()
	if err != nil {
		return nil, err
	}
	spaceRepo, err := rf.NewSpaceRepository()
	if err != nil {
		return nil, err
	}
	appUserGroupRepo, err := rf.NewAppUserGroupRepository()
	if err != nil {
		return nil, err
	}
	groupUserRepo, err := rf.NewGroupUserRepository()
	if err != nil {
		return nil, err
	}
	rbacRepo, err := rf.NewRBACRepository()
	if err != nil {
		return nil, err
	}

	m := &systemOwner{
		SystemOwnerModel: systemOwnerModel,
		rf:               rf,
		orgRepo:          orgRepo,
		spaceRepo:        spaceRepo,
		appUserGroupRepo: appUserGroupRepo,
		appUserRepo:      appUserRepo,
		groupUserRepo:    groupUserRepo,
		rbacRepo:         rbacRepo,
	}

	return m, libD.Validator.Struct(m)
}

func (s *systemOwner) GetOrganization(ctx context.Context) (Organization, error) {
	return s.orgRepo.GetOrganization(ctx, s)
}

func (s *systemOwner) FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error) {
	return s.appUserRepo.FindAppUserByID(ctx, s, id)
}

func (s *systemOwner) FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error) {
	return s.appUserRepo.FindAppUserByLoginID(ctx, s, loginID)
}

func (s *systemOwner) FindSystemSpace(ctx context.Context) (Space, error) {
	return s.spaceRepo.FindSystemSpace(ctx, s)
}

func (s *systemOwner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddStudent")
	appUserID, err := s.appUserRepo.AddAppUser(ctx, s, param)
	if err != nil {
		return 0, err
	}
	appUser, err := s.appUserRepo.FindAppUserByID(ctx, s, appUserID)
	if err != nil {
		return 0, err
	}

	// personalGroupID, err := s.rf.NewAppUserGroupRepository().AddPersonalGroup(s, studentID)
	// if err != nil {
	// 	return 0, err
	// }

	publicGroup, err := s.appUserGroupRepo.FindPublicGroup(ctx, s)
	if err != nil {
		return 0, err
	}
	if err := s.groupUserRepo.AddGroupUser(ctx, s, domain.AppUserGroupID(publicGroup.GetID()), domain.AppUserID(appUser.GetID())); err != nil {
		return 0, err
	}

	spaceID, err := s.spaceRepo.AddPersonalSpace(ctx, s, appUser)
	if err != nil {
		return 0, err
	}

	logger.Infof("Personal spaceID: %d", spaceID)

	spaceWriter := domain.NewSpaceWriterRole(spaceID)
	spaceObject := domain.NewSpaceObject(spaceID)
	userSubject := domain.NewUserObject(appUserID)

	if err := s.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "read"); err != nil {
		return 0, err
	}

	if err := s.rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "write"); err != nil {
		return 0, err
	}

	if err := s.rbacRepo.AddNamedGroupingPolicy(userSubject, spaceWriter); err != nil {
		return 0, err
	}

	// defaultSpace, err := s.rf.NewSpaceRepository().FindDefaultSpace(ctx, s)
	// if err != nil {
	// 	return 0, err
	// }

	// if err := s.rf.NewUserSpaceRepository().Add(ctx, appUser, SpaceID(defaultSpace.GetID())); err != nil {
	// 	return 0, err
	// }

	return appUserID, nil
}

func (s *systemOwner) AddSystemSpace(ctx context.Context) (domain.SpaceID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddSystemSpace")

	spaceID, err := s.spaceRepo.AddSystemSpace(ctx, s)
	if err != nil {
		return 0, err
	}
	return spaceID, nil
}
