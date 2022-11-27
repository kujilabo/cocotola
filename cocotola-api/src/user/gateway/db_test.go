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

		rf, err := gateway.NewRepositoryFactory(db)
		require.NoError(t, err)

		testService := testService{driverName: driverName, db: db, rf: rf}

		fn(ctx, testService)
	}
}

func testInitOrganization(t *testing.T, ts testService) (domain.OrganizationID, service.Owner) {
	bg := context.Background()
	sysAd, err := service.NewSystemAdminFromDB(bg, ts.db)
	assert.NoError(t, err)

	firstOwnerAddParam, err := service.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := service.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)

	// delete all organizations
	ts.db.Exec("delete from space")
	ts.db.Exec("delete from app_user")
	ts.db.Exec("delete from organization")
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&appUserEntity{})
	// db.Where("true").Delete(&organizationEntity{})

	orgRepo, err := gateway.NewOrganizationRepository(ts.db)
	require.NoError(t, err)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)

	appUserRepo, err := gateway.NewAppUserRepository(ts.rf, ts.db)
	assert.NoError(t, err)
	sysOwnerID, err := appUserRepo.AddSystemOwner(bg, sysAd, orgID)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(sysOwnerID)), 0)

	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(bg, sysAd, "ORG_NAME")
	assert.NoError(t, err)
	assert.Greater(t, int(uint(sysOwnerID)), 0)

	firstOwnerID, err := appUserRepo.AddFirstOwner(bg, sysOwner, firstOwnerAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(firstOwnerID)), 0)

	firstOwner, err := appUserRepo.FindOwnerByLoginID(bg, sysOwner, "OWNER_ID")
	assert.NoError(t, err)

	spaceRepo, err := gateway.NewSpaceRepository(ts.db)
	require.NoError(t, err)
	_, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	assert.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	assert.NoError(t, err)

	return orgID, firstOwner
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) service.AppUserAddParameter {
	p, err := service.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}
