package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

// func TestAddUser(t *testing.T) {
// 	for _, db := range dbList() {
// 		bg := context.Background()
// 		sqlDB, err := db.DB()
// 		assert.NoError(t, err)
// 		defer sqlDB.Close()

// 		model := user.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
// 		appUser, err := user.NewAppUser(nil, model, user.OrganizationID(1), "loginid", "username", []string{}, map[string]string{})
// 		assert.NoError(t, err)

// 		db.Debug().Where("id <> ?", 1).Delete(&appUserEntity{})
// 		repo := NewAppUserRepository(nil, db)
// 		_, err = repo.FindAppUserByID(bg, appUser, user.AppUserID(1))
// 		assert.NoError(t, err)
// 		// db.Delete(&organizationEntity{})

// 		// organizationID, err := initialize(db)
// 		// assert.NoError(t, err)
// 		// assert.Greater(t, organizationID, uint(0))
// 	}
// }

// func Test_appUserRepository_addAppUser(t *testing.T) {
// 	type fields struct {
// 		rf domain.RepositoryFactory
// 		db *gorm.DB
// 	}
// 	type args struct {
// 		ctx           context.Context
// 		appUserEntity *appUserEntity
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    domain.AppUserID
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &appUserRepository{
// 				rf: tt.fields.rf,
// 				db: tt.fields.db,
// 			}
// 			got, err := r.addAppUser(tt.args.ctx, tt.args.appUserEntity)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("appUserRepository.addAppUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("appUserRepository.addAppUser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_appUserRepository_AddAppUser(t *testing.T) {
	t.Parallel()
	fn := func(ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		type args struct {
			operator domain.OwnerModel
			param    service.AppUserAddParameter
		}
		tests := []struct {
			name string
			args args
			err  error
		}{
			{
				name: "success",
				args: args{
					operator: owner,
					param:    testNewAppUserAddParameter(t, "LOGIN_ID", "USERNAME"),
				},
				err: nil,
			},
			{
				name: "duplicated",
				args: args{
					operator: owner,
					param:    testNewAppUserAddParameter(t, "LOGIN_ID", "USERNAME"),
				},
				err: service.ErrAppUserAlreadyExists,
			},
		}
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.rf, ts.db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := appUserRepo.AddAppUser(ctx, tt.args.operator, tt.args.param)
				if err != nil && !errors.Is(err, tt.err) {
					t.Errorf("AddAppUser() error = %v, err %v", err, tt.err)
					return
				}
				if err == nil {
					assert.Greater(t, uint(got), uint(0))
				}
			})
		}
	}

	testDB(t, fn)
}
