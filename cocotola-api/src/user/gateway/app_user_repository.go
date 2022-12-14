package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
	"github.com/kujilabo/cocotola/lib/passwordhelper"
)

var (
	AppUserTableName = "app_user"

	SystemOwnerLoginID   = "system-owner"
	SystemStudentLoginID = "system-student"
	GuestLoginID         = "guest"

	AdministratorRole = "Administrator"
	OwnerRole         = "Owner"
	ManagerRole       = "Manager"
	UserRole          = "User"
	GuestRole         = "Guest"
	UnknownRole       = "Unknown"
)

type appUserRepository struct {
	rf service.RepositoryFactory
	db *gorm.DB
}

type appUserEntity struct {
	SinmpleModelEntity
	OrganizationID       uint
	LoginID              string
	Username             string
	HashedPassword       string
	Role                 string
	Provider             string
	ProviderID           string
	ProviderAccessToken  string
	ProviderRefreshToken string
	Removed              bool
}

func (e *appUserEntity) TableName() string {
	return AppUserTableName
}

// func toRole(role string) domain.Role {
// 	if role == "administrator" {
// 		return domain.AdministratorRole
// 	} else if role == OwnerRole {
// 		return domain.OwnerRole
// 	} else if role == "Manager" {
// 		return domain.ManagerRole
// 	} else if role == "User" {
// 		return domain.UserRole
// 	} else if role == "Guest" {
// 		return domain.GuestRole
// 	}
// 	return domain.UnknownRole
// }

// func fromRoleToString(role domain.Role) string {
// 	switch role {
// 	case domain.AdministratorRole:
// 		return AdministratorRole
// 	case domain.OwnerRole:
// 		return OwnerRole
// 	case domain.ManagerRole:
// 		return ManagerRole
// 	case domain.UserRole:
// 		return UserRole
// 	case domain.GuestRole:
// 		return GuestRole
// 	default:
// 		return UnknownRole
// 	}
// }

func (e *appUserEntity) toSystemOwner(ctx context.Context, rf service.RepositoryFactory) (service.SystemOwner, error) {
	if e.LoginID != SystemOwnerLoginID {
		return nil, liberrors.Errorf("invalid system owner. loginID: %s", e.LoginID)
	}

	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewModel. err: %w", err)
	}

	appUserModel, err := domain.NewAppUserModel(model, domain.OrganizationID(e.OrganizationID), e.LoginID, e.Username, []string{"SystemOwner"}, map[string]string{})
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	ownerModel, err := domain.NewOwnerModel(appUserModel)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOwnerModel. err: %w", err)
	}

	systemOwnerModel, err := domain.NewSystemOwnerModel(ownerModel)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewSystemOwnerModel. err: %w", err)
	}

	systemOwner, err := service.NewSystemOwner(ctx, rf, systemOwnerModel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewSystemOwner. err: %w", err)
	}

	return systemOwner, nil
}

func (e *appUserEntity) toAppUserModel(roles []string, properties map[string]string) (domain.AppUserModel, error) {
	model, err := e.toModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toModel. err: %w", err)
	}
	appUserModel, err := domain.NewAppUserModel(model, domain.OrganizationID(e.OrganizationID), e.LoginID, e.Username, roles, properties)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return appUserModel, nil
}

func (e *appUserEntity) toAppUser(ctx context.Context, rf service.RepositoryFactory, roles []string, properties map[string]string) (service.AppUser, error) {
	appUserModel, err := e.toAppUserModel(roles, properties)
	if err != nil {
		return nil, liberrors.Errorf("e.toAppUserModel. err: %w", err)
	}

	appUser, err := service.NewAppUser(ctx, rf, appUserModel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewAppUser. err: %w", err)
	}

	return appUser, nil
}

func (e *appUserEntity) toOwner(rf service.RepositoryFactory, roles []string, properties map[string]string) (service.Owner, error) {
	appUserModel, err := e.toAppUserModel(roles, properties)
	if err != nil {
		return nil, liberrors.Errorf("e.toAppUserModel. err: %w", err)
	}

	ownerModel, err := domain.NewOwnerModel(appUserModel)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewOwnerModel. err: %w", err)
	}

	return service.NewOwner(rf, ownerModel), nil
}

func NewAppUserRepository(ctx context.Context, rf service.RepositoryFactory, db *gorm.DB) service.AppUserRepository {
	if rf == nil {
		panic(errors.New("rf is nil"))
	} else if db == nil {
		panic(errors.New("db is nil"))
	}
	return &appUserRepository{
		rf: rf,
		db: db,
	}
}

func (r *appUserRepository) FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.SystemOwner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindSystemOwnerByOrganizationID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where("organization_id = ? and removed = 0", organizationID).
		Where("login_id = ?", SystemOwnerLoginID).
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, liberrors.Errorf("system owner not found. organization ID: %d, err: %w", organizationID, service.ErrSystemOwnerNotFound)
		}
		return nil, result.Error
	}
	return appUser.toSystemOwner(ctx, r.rf)
}

func (r *appUserRepository) FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminModel, organizationName string) (service.SystemOwner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindSystemOwnerByOrganizationName")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Table("organization").Select("app_user.*").
		Where("organization.name = ? and app_user.removed = 0", organizationName).
		Where("login_id = ?", SystemOwnerLoginID).
		Joins("inner join app_user on organization.id = app_user.organization_id").
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, liberrors.Errorf("system owner not found. organization name: %s, err: %w", organizationName, service.ErrSystemOwnerNotFound)
		}

		return nil, result.Error
	}
	return appUser.toSystemOwner(ctx, r.rf)
}

func (r *appUserRepository) FindAppUserByID(ctx context.Context, operator domain.AppUserModel, id domain.AppUserID) (service.AppUser, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			ID: uint(id),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	roles := []string{appUser.Role}
	properties := map[string]string{}

	return appUser.toAppUser(ctx, r.rf, roles, properties)
}

func (r *appUserRepository) FindAppUserByLoginID(ctx context.Context, operator domain.AppUserModel, loginID string) (service.AppUser, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByLoginID")
	defer span.End()

	if loginID == "" {
		return nil, errors.New("invalid parameter")
	}

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        loginID,
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	roles := []string{appUser.Role}
	properties := map[string]string{}

	return appUser.toAppUser(ctx, r.rf, roles, properties)
}

func (r *appUserRepository) FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwnerModel, loginID string) (service.Owner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindOwnerByLoginID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        loginID,
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	roles := []string{appUser.Role}
	properties := map[string]string{}

	return appUser.toOwner(r.rf, roles, properties)
}

func (r *appUserRepository) addAppUser(ctx context.Context, appUserEntity *appUserEntity) (domain.AppUserID, error) {
	if result := r.db.Create(appUserEntity); result.Error != nil {
		return 0, liberrors.Errorf(". err: %w", libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}
	return domain.AppUserID(appUserEntity.ID), nil
}

func (r *appUserRepository) AddAppUser(ctx context.Context, operator domain.OwnerModel, param service.AppUserAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddAppUser")
	defer span.End()

	hashedPassword := ""
	password, ok := param.GetProperties()["password"]
	if ok {
		hashed, err := passwordhelper.HashPassword(password)
		if err != nil {
			return 0, liberrors.Errorf("passwordhelper.HashPassword. err: %w", err)
		}

		hashedPassword = hashed
	}

	appUserEntity := appUserEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        param.GetLoginID(),
		Username:       param.GetUsername(),
		HashedPassword: hashedPassword,
		Role:           UserRole,
	}
	return r.addAppUser(ctx, &appUserEntity)
}

func (r *appUserRepository) AddSystemOwner(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddSystemOwner")
	defer span.End()

	appUserEntity := appUserEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		OrganizationID: uint(organizationID),
		LoginID:        SystemOwnerLoginID,
		Username:       "SystemOwner",
		Role:           OwnerRole,
	}
	return r.addAppUser(ctx, &appUserEntity)
}

func (r *appUserRepository) AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param service.FirstOwnerAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddFirstOwner")
	defer span.End()

	hashedPassword, err := passwordhelper.HashPassword(param.GetPassword())
	if err != nil {
		return 0, liberrors.Errorf("passwordhelper.HashPassword. err: %w", err)
	}

	appUserEntity := appUserEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        param.GetLoginID(),
		Username:       param.GetUsername(),
		HashedPassword: hashedPassword,
		Role:           OwnerRole,
	}
	return r.addAppUser(ctx, &appUserEntity)
}

func (r *appUserRepository) FindAppUserIDs(ctx context.Context, operator domain.SystemOwnerModel, pageNo, pageSize int) ([]domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserIDs")
	defer span.End()

	limit := pageSize
	offset := (pageNo - 1) * pageSize

	var entities []appUserEntity
	if result := r.db.Where(&appUserEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
	}).Limit(limit).Offset(offset).Order("id").Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	ids := make([]domain.AppUserID, len(entities))
	for i, entity := range entities {
		ids[i] = domain.AppUserID(entity.ID)
	}

	return ids, nil
}
