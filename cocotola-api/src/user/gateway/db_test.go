package gateway_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

const orgNameLength = 8

type testService struct {
	driverName string
	db         *gorm.DB
	rf         service.RepositoryFactory
}

func testDB(t *testing.T, fn func(ctx context.Context, ts testService)) {
	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		logrus.Debugf("%s\n", driverName)
		sqlDB, err := db.DB()
		require.NoError(t, err)
		defer sqlDB.Close()

		rf, err := gateway.NewRepositoryFactory(ctx, db)
		require.NoError(t, err)

		testService := testService{driverName: driverName, db: db, rf: rf}

		fn(ctx, testService)
	}
}

func setupOrganization(ctx context.Context, t *testing.T, ts testService) (domain.OrganizationID, service.Owner) {
	bg := context.Background()
	orgName := RandString(orgNameLength)
	sysAd, err := service.NewSystemAdmin(ctx, ts.rf)
	assert.NoError(t, err)

	firstOwnerAddParam, err := service.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := service.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
	assert.NoError(t, err)

	orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)

	appUserRepo := gateway.NewAppUserRepository(ctx, ts.rf, ts.db)
	sysOwnerID, err := appUserRepo.AddSystemOwner(bg, sysAd, orgID)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(sysOwnerID)), 0)

	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(bg, sysAd, orgName)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(sysOwnerID)), 0)

	firstOwnerID, err := appUserRepo.AddFirstOwner(bg, sysOwner, firstOwnerAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(firstOwnerID)), 0)

	firstOwner, err := appUserRepo.FindOwnerByLoginID(bg, sysOwner, "OWNER_ID")
	assert.NoError(t, err)

	spaceRepo := gateway.NewSpaceRepository(ctx, ts.db)
	_, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	assert.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	assert.NoError(t, err)

	return orgID, firstOwner
}

func teardownOrganization(t *testing.T, ts testService, orgID domain.OrganizationID) {
	// delete all organizations
	ts.db.Exec("delete from space where organization_id = ?", uint(orgID))
	ts.db.Exec("delete from app_user where organization_id = ?", uint(orgID))
	ts.db.Exec("delete from organization where id = ?", uint(orgID))
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&appUserEntity{})
	// db.Where("true").Delete(&organizationEntity{})
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) service.AppUserAddParameter {
	p, err := service.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}
