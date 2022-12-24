package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
)

var (
	AppUserGroupTableName = "app_user_group"
)

type appUserGroupRepository struct {
	db *gorm.DB
}

type appUserGroupEntity struct {
	SinmpleModelEntity
	OrganizationID uint
	Key            string
	Name           string
	Description    string
}

func (e *appUserGroupEntity) TableName() string {
	return AppUserGroupTableName
}

func (e *appUserGroupEntity) toAppUserGroup() (service.AppUserGroup, error) {
	model, err := e.toModel()
	if err != nil {
		return nil, liberrors.Errorf("toAppUserGroup. err: %w", err)
	}

	appUserGroupMdoel, err := domain.NewAppUserGroup(model, domain.OrganizationID(e.OrganizationID), e.Key, e.Name, e.Description)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserGroup. err: %w", err)
	}

	appUserGroup, err := service.NewAppUserGroup(appUserGroupMdoel)
	if err != nil {
		return nil, liberrors.Errorf("service.NewAppUserGroup. err: %w", err)
	}

	return appUserGroup, nil
}

func NewAppUserGroupRepository(ctx context.Context, db *gorm.DB) service.AppUserGroupRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	return &appUserGroupRepository{
		db: db,
	}
}

func (r *appUserGroupRepository) FindPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (service.AppUserGroup, error) {
	_, span := tracer.Start(ctx, "appUserGroupRepository.FindPublicGroup")
	defer span.End()

	appUserGroup := appUserGroupEntity{}
	if result := r.db.Where(&appUserGroupEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "public",
	}).Find(&appUserGroup); result.Error != nil {
		return nil, result.Error
	}
	return appUserGroup.toAppUserGroup()
}

func (r *appUserGroupRepository) AddPublicGroup(ctx context.Context, operator domain.SystemOwnerModel) (domain.AppUserGroupID, error) {
	_, span := tracer.Start(ctx, "appUserGroupRepository.AddPublicGroup")
	defer span.End()

	appUserGroup := appUserGroupEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "public",
		Name:           "Public group",
	}
	if result := r.db.Create(&appUserGroup); result.Error != nil {
		return 0, liberrors.Errorf(". err: %w", libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}
	return domain.AppUserGroupID(appUserGroup.ID), nil
}

func (r *appUserGroupRepository) AddPersonalGroup(ctx context.Context, operator domain.AppUserModel) (uint, error) {
	_, span := tracer.Start(ctx, "appUserGroupRepository.AddPersonalGroup")
	defer span.End()

	appUserGroup := appUserGroupEntity{
		SinmpleModelEntity: SinmpleModelEntity{
			Version:   1,
			CreatedBy: operator.GetID(),
			UpdatedBy: operator.GetID(),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
		Key:            "#" + operator.GetLoginID(),
		Name:           "Personal group",
	}
	if result := r.db.Create(&appUserGroup); result.Error != nil {
		return 0, liberrors.Errorf(". err: %w", libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}
	return appUserGroup.ID, nil
}
