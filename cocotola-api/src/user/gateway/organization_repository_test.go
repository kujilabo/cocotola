package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

const invalidOrgID = 99999

func TestGetOrganization(t *testing.T) {
	fn := func(ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, _ := setupOrganization(t, ts)
		defer teardownOrganization(t, ts, orgID)

		orgRepo, err := gateway.NewOrganizationRepository(ts.db)
		require.NoError(t, err)

		// get organization registered
		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		userModel, err := domain.NewAppUserModel(model, orgID, "login_id", "username", []string{}, map[string]string{})
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
		}

		// get organization unregistered
		otherUserModel, err := domain.NewAppUserModel(model, invalidOrgID, "login_id", "username", []string{}, map[string]string{})
		assert.NoError(t, err)
		{
			_, err := orgRepo.GetOrganization(ctx, otherUserModel)
			assert.Equal(t, service.ErrOrganizationNotFound, err)
		}
	}
	testDB(t, fn)
}

func TestFindOrganizationByName(t *testing.T) {
	fn := func(ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, _ := setupOrganization(t, ts)
		defer teardownOrganization(t, ts, orgID)
		systemAdminModel := domain.NewSystemAdminModel()

		orgRepo, err := gateway.NewOrganizationRepository(ts.db)
		require.NoError(t, err)

		var orgName string

		// get organization registered
		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		userModel, err := domain.NewAppUserModel(model, orgID, "login_id", "username", []string{}, map[string]string{})
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
			orgName = org.GetName()
		}

		// find organization registered by name
		{
			org, err := orgRepo.FindOrganizationByName(ctx, systemAdminModel, orgName)
			assert.NoError(t, err)
			assert.Equal(t, orgName, org.GetName())
		}

		// find organization unregistered by name
		{
			_, err := orgRepo.FindOrganizationByName(ctx, systemAdminModel, "NOT_FOUND")
			assert.Equal(t, service.ErrOrganizationNotFound, err)
		}
	}
	testDB(t, fn)
}
