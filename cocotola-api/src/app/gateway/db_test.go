package gateway_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

func testInitOrganization(t *testing.T, db *gorm.DB) (userD.OrganizationID, userS.SystemOwner, userS.Owner) {
	log.Println("testInitOrganization")
	bg := context.Background()
	sysAd, err := userS.NewSystemAdminFromDB(bg, db)
	assert.NoError(t, err)

	// delete all organizations
	result := db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from workbook")
	assert.NoError(t, result.Error)
	result = db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from space")
	assert.NoError(t, result.Error)
	result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from app_user")
	assert.NoError(t, result.Error)
	result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from organization")
	assert.NoError(t, result.Error)

	firstOwnerAddParam, err := userS.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := userS.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	orgRepo := userG.NewOrganizationRepository(db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)
	log.Printf("OrgID: %d \n", uint(orgID))
	org, err := orgRepo.FindOrganizationByID(bg, sysAd, orgID)
	assert.NoError(t, err)
	log.Printf("OrgID: %d \n", org.GetID())

	appUserRepo := userG.NewAppUserRepository(nil, db)
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

	spaceRepo := userG.NewSpaceRepository(db)
	_, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	assert.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	assert.NoError(t, err)

	return orgID, sysOwner, firstOwner
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) userS.AppUserAddParameter {
	p, err := userS.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}
