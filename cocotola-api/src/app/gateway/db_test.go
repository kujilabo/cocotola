package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	service_mock "github.com/kujilabo/cocotola/cocotola-api/src/app/service/mock"
	jobG "github.com/kujilabo/cocotola/cocotola-api/src/job/gateway"
	jobS "github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userG "github.com/kujilabo/cocotola/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	testlibG "github.com/kujilabo/cocotola/test-lib/gateway"
)

const englishWordName = "english_word"
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
}

func testDB(t *testing.T, fn func(t *testing.T, ctx context.Context, ts testService)) {
	logrus.SetLevel(logrus.WarnLevel)
	// englishWord := testNewProblemType(t, englishWordID, englishWordName)
	// memorization := testNewStudyType(t, memorizationID, memorizationName)
	// dictation := testNewStudyType(t, dictationID, dictationName)
	// problemTypes := []domain.ProblemType{englishWord}
	// studyTypes := []domain.StudyType{memorization, dictation}

	problemAddProcessor := map[domain.ProblemTypeName]service.ProblemAddProcessor{}
	problemUpdateProcessor := map[domain.ProblemTypeName]service.ProblemUpdateProcessor{}
	problemRemoveProcessor := map[domain.ProblemTypeName]service.ProblemRemoveProcessor{}
	problemImportProcessor := map[domain.ProblemTypeName]service.ProblemImportProcessor{}
	problemQuotaProcessor := map[domain.ProblemTypeName]service.ProblemQuotaProcessor{}

	jobRff := func(ctx context.Context, db *gorm.DB) (jobS.RepositoryFactory, error) {
		return jobG.NewRepositoryFactory(ctx, db)
	}
	userRff := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(ctx, db)
	}

	pf := service.NewProcessorFactory(problemAddProcessor, problemUpdateProcessor, problemRemoveProcessor, problemImportProcessor, problemQuotaProcessor)

	ctx := context.Background()
	location := time.Local
	for driverName, db := range testlibG.ListDB() {
		driverName := driverName
		db := db
		t.Run(driverName, func(t *testing.T) {
			t.Parallel()
			logrus.Debugf("%s\n", driverName)
			sqlDB, err := db.DB()
			require.NoError(t, err)
			defer sqlDB.Close()

			problemRepository := new(service_mock.ProblemRepository)
			problemRepositories := map[domain.ProblemTypeName]func(context.Context, *gorm.DB) (service.ProblemRepository, error){
				englishWordName: func(context.Context, *gorm.DB) (service.ProblemRepository, error) {
					return problemRepository, nil
				},
			}
			rf, err := gateway.NewRepositoryFactory(ctx, db, driverName, location, jobRff, userRff, pf, problemRepositories)
			require.NoError(t, err)
			testService := testService{driverName: driverName, db: db, pf: pf, rf: rf}

			fn(t, ctx, testService)
		})
	}
}

func setupOrganization(ctx context.Context, t *testing.T, ts testService) (userD.OrganizationID, userS.SystemOwner, userS.Owner) {
	orgName := RandString(orgNameLength)
	userRf, err := ts.rf.NewUserRepositoryFactory(ctx)
	require.NoError(t, err)
	sysAd, err := userS.NewSystemAdmin(ctx, userRf)
	require.NoError(t, err)

	firstOwnerAddParam, err := userS.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	require.NoError(t, err)
	orgAddParam, err := userS.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
	require.NoError(t, err)
	orgRepo := userG.NewOrganizationRepository(ctx, ts.db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(ctx, sysAd, orgAddParam)
	require.NoError(t, err)
	require.Greater(t, int(uint(orgID)), 0)
	logrus.Debugf("OrgID: %d \n", uint(orgID))
	org, err := orgRepo.FindOrganizationByID(ctx, sysAd, orgID)
	require.NoError(t, err)
	logrus.Debugf("OrgID: %d \n", org.GetID())

	appUserRepo := userG.NewAppUserRepository(ctx, userRf, ts.db)
	sysOwnerID, err := appUserRepo.AddSystemOwner(ctx, sysAd, orgID)
	require.NoError(t, err)
	require.Greater(t, int(uint(sysOwnerID)), 0)

	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, orgName)
	require.NoError(t, err)
	require.Greater(t, int(uint(sysOwnerID)), 0)

	firstOwnerID, err := appUserRepo.AddFirstOwner(ctx, sysOwner, firstOwnerAddParam)
	require.NoError(t, err)
	require.Greater(t, int(uint(firstOwnerID)), 0)

	firstOwner, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, "OWNER_ID")
	require.NoError(t, err)

	spaceRepo := userG.NewSpaceRepository(ctx, ts.db)
	_, err = spaceRepo.AddDefaultSpace(ctx, sysOwner)
	require.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(ctx, sysOwner, firstOwner)
	require.NoError(t, err)

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
	require.NoError(t, err)
	return p
}

func testNewProblemType(t *testing.T, id uint, name string) domain.ProblemType {
	p, err := domain.NewProblemType(id, name)
	require.NoError(t, err)
	return p
}

func testNewStudyType(t *testing.T, id uint, name string) domain.StudyType {
	p, err := domain.NewStudyType(id, name)
	require.NoError(t, err)
	return p
}

func testNewWorkbookSearchCondition(t *testing.T) service.WorkbookSearchCondition {
	p, err := service.NewWorkbookSearchCondition(1, 10, []userD.SpaceID{})
	require.NoError(t, err)
	return p
}

func testNewWorkbookAddParameter(t *testing.T, name string) service.WorkbookAddParameter {
	p, err := service.NewWorkbookAddParameter("english_word", name, domain.Lang2JA, "", map[string]string{"audioEnabled": "false"})
	require.NoError(t, err)
	return p
}

func testNewAppUser(t *testing.T, ctx context.Context, ts testService, sysOwner userS.SystemOwner, owner userS.Owner, loginID, username string) userS.AppUser {
	userRf, err := ts.rf.NewUserRepositoryFactory(ctx)
	require.NoError(t, err)
	appUserRepo := userG.NewAppUserRepository(ctx, userRf, ts.db)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username))
	require.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	require.NoError(t, err)
	require.Equal(t, loginID, user1.GetLoginID())

	spaceRepo := userG.NewSpaceRepository(ctx, ts.db)

	spaceID1, err := spaceRepo.AddPersonalSpace(ctx, sysOwner, user1)
	require.NoError(t, err)

	space, err := user1.GetPersonalSpace(ctx)
	require.NoError(t, err)
	require.Equal(t, spaceID1, userD.SpaceID(space.GetID()))

	return user1
}

func testNewStudentModel(t *testing.T, appUserModel userD.AppUserModel) domain.StudentModel {
	s, err := domain.NewStudentModel(appUserModel)
	require.NoError(t, err)
	return s
}

func testNewStudent(ctx context.Context, t *testing.T, ts testService, appUserModel userD.AppUserModel) service.Student {
	studentModel := testNewStudentModel(t, appUserModel)
	s, err := service.NewStudent(ctx, ts.pf, ts.rf, studentModel)
	require.NoError(t, err)
	return s
}

func testNewWorkbook(t *testing.T, ctx context.Context, db *gorm.DB, workbookRepo service.WorkbookRepository, student service.Student, spaceID userD.SpaceID, workbookName string) service.Workbook {
	workbookID11, err := workbookRepo.AddWorkbook(ctx, student, spaceID, testNewWorkbookAddParameter(t, workbookName))
	require.NoError(t, err)
	require.Greater(t, int(workbookID11), 0)
	workbook, err := workbookRepo.FindWorkbookByID(ctx, student, workbookID11)
	require.NoError(t, err)
	return workbook
}
