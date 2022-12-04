package gateway_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

const englishWordName = "english_word_problem"
const englishWordID = 1
const memorizationName = "memorization"
const memorizationID = 1
const dictationName = "dictation"
const dictationID = 2
const orgNameLength = 8

type testService struct {
	driverName string
	db         *gorm.DB
	pf         service.ProcessorFactory
	rf         service.RepositoryFactory
	userRf     userS.RepositoryFactory
}

func testDB(t *testing.T, fn func(ctx context.Context, ts testService)) {
	logrus.SetLevel(logrus.WarnLevel)
	englishWord := testNewProblemType(t, englishWordID, englishWordName)
	memorization := testNewStudyType(t, memorizationID, memorizationName)
	dictation := testNewStudyType(t, dictationID, dictationName)
	problemTypes := []domain.ProblemType{englishWord}
	studyTypes := []domain.StudyType{memorization, dictation}

	problemAddProcessor := map[string]service.ProblemAddProcessor{}
	problemUpdateProcessor := map[string]service.ProblemUpdateProcessor{}
	problemRemoveProcessor := map[string]service.ProblemRemoveProcessor{}
	problemImportProcessor := map[string]service.ProblemImportProcessor{}
	problemQuotaProcessor := map[string]service.ProblemQuotaProcessor{}

	pf := service.NewProcessorFactory(problemAddProcessor, problemUpdateProcessor, problemRemoveProcessor, problemImportProcessor, problemQuotaProcessor)

	ctx := context.Background()
	for driverName, db := range testlibG.ListDB() {
		logrus.Debugf("%s\n", driverName)
		sqlDB, err := db.DB()
		require.NoError(t, err)
		defer sqlDB.Close()

		rbacRepo, err := userG.NewRBACRepository(db)
		require.NoError(t, err)
		err = rbacRepo.Init()
		require.NoError(t, err)

		userRf, err := userG.NewRepositoryFactory(db)
		require.NoError(t, err)
		problemRepositories := map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error){}
		rf, err := gateway.NewRepositoryFactory(ctx, db, driverName, userRfFunc, pf, problemTypes, studyTypes, problemRepositories)
		require.NoError(t, err)
		testService := testService{driverName: driverName, db: db, pf: pf, rf: rf, userRf: userRf}

		fn(ctx, testService)
	}
}

func setupOrganization(t *testing.T, ts testService) (userD.OrganizationID, userS.SystemOwner, userS.Owner) {
	bg := context.Background()
	orgName := RandString(orgNameLength)
	sysAd, err := userS.NewSystemAdminFromDB(bg, ts.db)
	assert.NoError(t, err)

	firstOwnerAddParam, err := userS.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := userS.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
	assert.NoError(t, err)
	orgRepo, err := userG.NewOrganizationRepository(ts.db)
	require.NoError(t, err)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)
	logrus.Debugf("OrgID: %d \n", uint(orgID))
	org, err := orgRepo.FindOrganizationByID(bg, sysAd, orgID)
	assert.NoError(t, err)
	logrus.Debugf("OrgID: %d \n", org.GetID())

	appUserRepo, err := userG.NewAppUserRepository(ts.userRf, ts.db)
	assert.NoError(t, err)
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

	spaceRepo, err := userG.NewSpaceRepository(ts.db)
	require.NoError(t, err)
	_, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	assert.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	assert.NoError(t, err)

	return orgID, sysOwner, firstOwner
}

func teardownOrganization(t *testing.T, ts testService, orgID userD.OrganizationID) {
	// delete all organizations
	result := ts.db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from study_record where organization_id = ?", uint(orgID))
	assert.NoError(t, result.Error)
	result = ts.db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from workbook where organization_id = ?", uint(orgID))
	assert.NoError(t, result.Error)
	result = ts.db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from space where organization_id = ?", uint(orgID))
	assert.NoError(t, result.Error)
	result = ts.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from app_user where organization_id = ?", uint(orgID))
	assert.NoError(t, result.Error)
	result = ts.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from organization where id = ?", uint(orgID))
	assert.NoError(t, result.Error)
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) userS.AppUserAddParameter {
	p, err := userS.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}

func testNewProblemType(t *testing.T, id uint, name string) domain.ProblemType {
	p, err := domain.NewProblemType(id, name)
	assert.NoError(t, err)
	return p
}

func testNewStudyType(t *testing.T, id uint, name string) domain.StudyType {
	p, err := domain.NewStudyType(id, name)
	assert.NoError(t, err)
	return p
}

func testNewWorkbookSearchCondition(t *testing.T) service.WorkbookSearchCondition {
	p, err := service.NewWorkbookSearchCondition(1, 10, []userD.SpaceID{})
	assert.NoError(t, err)
	return p
}

func testNewWorkbookAddParameter(t *testing.T, name string) service.WorkbookAddParameter {
	p, err := service.NewWorkbookAddParameter("english_word_problem", name, domain.Lang2JA, "", map[string]string{"audioEnabled": "false"})
	assert.NoError(t, err)
	return p
}

func testNewAppUser(t *testing.T, ctx context.Context, ts testService, sysOwner userS.SystemOwner, owner userS.Owner, loginID, username string) userS.AppUser {
	appUserRepo, err := userG.NewAppUserRepository(ts.userRf, ts.db)
	assert.NoError(t, err)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username))
	assert.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	assert.NoError(t, err)
	assert.Equal(t, loginID, user1.GetLoginID())

	spaceRepo, err := userG.NewSpaceRepository(ts.db)
	require.NoError(t, err)

	spaceID1, err := spaceRepo.AddPersonalSpace(ctx, sysOwner, user1)
	assert.NoError(t, err)

	space, err := user1.GetPersonalSpace(ctx)
	require.NoError(t, err)
	assert.Equal(t, spaceID1, userD.SpaceID(space.GetID()))

	return user1
}

func testNewStudentModel(t *testing.T, appUserModel userD.AppUserModel) domain.StudentModel {
	s, err := domain.NewStudentModel(appUserModel)
	assert.NoError(t, err)
	return s
}

func testNewStudent(t *testing.T, testService testService, appUserModel userD.AppUserModel) service.Student {
	studentModel := testNewStudentModel(t, appUserModel)
	s, err := service.NewStudent(testService.pf, testService.rf, testService.userRf, studentModel)
	assert.NoError(t, err)
	return s
}

func testNewWorkbook(t *testing.T, ctx context.Context, db *gorm.DB, workbookRepo service.WorkbookRepository, student service.Student, spaceID userD.SpaceID, workbookName string) service.Workbook {
	workbookID11, err := workbookRepo.AddWorkbook(ctx, student, spaceID, testNewWorkbookAddParameter(t, workbookName))
	assert.NoError(t, err)
	assert.Greater(t, int(workbookID11), 0)
	workbook, err := workbookRepo.FindWorkbookByID(ctx, student, workbookID11)
	assert.NoError(t, err)
	return workbook
}
