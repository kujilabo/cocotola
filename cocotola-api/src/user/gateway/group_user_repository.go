package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	libG "github.com/kujilabo/cocotola/lib/gateway"
)

var (
	GroupUserTableName = "group_user"
)

type groupUserRepository struct {
	db *gorm.DB
}

type groupUserEntity struct {
	JunctionModelEntity
	OrganizationID uint
	AppUserGroupID uint
	AppUserID      uint
}

func (u *groupUserEntity) TableName() string {
	return GroupUserTableName
}

func NewGroupUserRepository(ctx context.Context, db *gorm.DB) service.GroupUserRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	return &groupUserRepository{
		db: db,
	}
}

func (r *groupUserRepository) AddGroupUser(ctx context.Context, operator domain.AppUserModel, appUserGroupID domain.AppUserGroupID, appUserID domain.AppUserID) error {
	_, span := tracer.Start(ctx, "groupUserRepository.AddGroupUser")
	defer span.End()

	groupUser := groupUserEntity{
		JunctionModelEntity: JunctionModelEntity{
			CreatedBy: operator.GetID(),
		},
		OrganizationID: uint(operator.GetOrganizationID()),
		AppUserGroupID: uint(appUserGroupID),
		AppUserID:      uint(appUserID),
	}
	if result := r.db.Create(&groupUser); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists)
	}
	return nil
}
