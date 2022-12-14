package service

import (
	"context"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

type GroupUserRepository interface {
	AddGroupUser(ctx context.Context, operator domain.AppUserModel, appUserGroupID domain.AppUserGroupID, appUserID domain.AppUserID) error
}
