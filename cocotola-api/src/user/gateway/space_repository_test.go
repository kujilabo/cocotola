package gateway_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

func Test_spaceRepository_FindDefaultSpace(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(ctx context.Context, ts testService) {
		orgID, owner := setupOrganization(t, ts)
		defer teardownOrganization(t, ts, orgID)

		type args struct {
			operator domain.AppUserModel
		}

		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		spaceModel, err := domain.NewSpaceModel(model, orgID, 1, "default", "Default", "")
		assert.NoError(t, err)
		space, err := service.NewSpace(spaceModel)
		assert.NoError(t, err)
		tests := []struct {
			name string
			args args
			want service.Space
			err  error
		}{
			{
				name: "",
				args: args{
					operator: owner,
				},
				want: space,
				err:  nil,
			},
		}
		spaceRepo, err := ts.rf.NewSpaceRepository()
		require.NoError(t, err)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := spaceRepo.FindDefaultSpace(ctx, tt.args.operator)
				if err != nil && !errors.Is(err, tt.err) {
					t.Errorf("spaceRepository.FindDefaultSpace() error = %v, err %v", err, tt.err)
					return
				}
				if err == nil {
					assert.Equal(t, space.GetKey(), got.GetKey())
					assert.Equal(t, space.GetName(), got.GetName())
					assert.Equal(t, space.GetDescription(), got.GetDescription())
				}
			})
		}
	}
	testDB(t, fn)
}

func Test_spaceRepository_FindPersonalSpace(t *testing.T) {

	fn := func(ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, owner := setupOrganization(t, ts)
		defer teardownOrganization(t, ts, orgID)

		type args struct {
			operator domain.AppUserModel
		}

		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		spaceModel, err := domain.NewSpaceModel(model, orgID, 1, strconv.Itoa(int(owner.GetID())), "Default", "")
		assert.NoError(t, err)
		space, err := service.NewSpace(spaceModel)
		assert.NoError(t, err)
		tests := []struct {
			name string
			args args
			want service.Space
			err  error
		}{
			{
				name: "",
				args: args{
					operator: owner,
				},
				want: space,
				err:  nil,
			},
		}
		spaceRepo, err := ts.rf.NewSpaceRepository()
		require.NoError(t, err)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := spaceRepo.FindPersonalSpace(ctx, tt.args.operator)
				if err != nil && !errors.Is(err, tt.err) {
					t.Errorf("spaceRepository.FindPersonalSpace() error = %v, err %v", err, tt.err)
					return
				}
				if err == nil {
					assert.Equal(t, space.GetKey(), got.GetKey())
					assert.Equal(t, space.GetName(), got.GetName())
					assert.Equal(t, space.GetDescription(), got.GetDescription())
				}
			})
		}
	}
	testDB(t, fn)
}
