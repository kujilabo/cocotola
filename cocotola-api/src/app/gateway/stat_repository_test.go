package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	// "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/gateway"
	// userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

// func tset_statRepository_init(t *testing.T, bg context.Context, driverName string, db *gorm.DB) (userD.AppUserID, userD.AppUserID, domain.WorkbookID, domain.WorkbookID) {

// 	fn := func(ctx context.Context, ts testService) {
// 	userRepo, err := userG.NewRepositoryFactory(db)
// 	assert.NoError(t, err)
// 	_, sysOwner, owner := testInitOrganization(t, db)
// 	appUserRepo := userG.NewAppUserRepository(nil, db)

// 	rbacRepo := userG.NewRBACRepository(db)
// 	err = rbacRepo.Init()
// 	assert.NoError(t, err)

// 	userID1, err := appUserRepo.AddAppUser(bg, owner, testNewAppUserAddParameter(t, "LOGIN_ID_1", "USERNAME_1"))
// 	assert.NoError(t, err)
// 	user1, err := appUserRepo.FindAppUserByID(bg, owner, userID1)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "LOGIN_ID_1", user1.GetLoginID())
// 	userID2, err := appUserRepo.AddAppUser(bg, owner, testNewAppUserAddParameter(t, "LOGIN_ID_2", "USERNAME_2"))
// 	assert.NoError(t, err)
// 	user2, err := appUserRepo.FindAppUserByID(bg, owner, userID2)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "LOGIN_ID_2", user2.GetLoginID())

// 	englishWord := testNewProblemType(t, "english_word_problem")
// 	workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []domain.ProblemType{englishWord})
// 	spaceRepo := userG.NewSpaceRepository(db)

// 	// user1 has two workbooks
// 	student1 := testNewStudent(t, ts, user1)
// 	spaceID1, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user1)
// 	assert.NoError(t, err)
// 	workbookID11, err := workbookRepo.AddWorkbook(bg, student1, spaceID1, testNewWorkbookAddParameter(t, "WB11"))
// 	assert.NoError(t, err)
// 	workbookID12, err := workbookRepo.AddWorkbook(bg, student1, spaceID1, testNewWorkbookAddParameter(t, "WB12"))
// 	assert.NoError(t, err)

// 	// user2 has one workbook
// 	student2 := testNewStudent(t, ts, user2)
// 	spaceID2, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user2)
// 	assert.NoError(t, err)
// 	workbookID21, err := workbookRepo.AddWorkbook(bg, student2, spaceID2, testNewWorkbookAddParameter(t, "WB21"))
// 	assert.NoError(t, err)
// 	assert.GreaterOrEqual(t, uint(workbookID21), uint(1))
// 	return userID1, userID2, workbookID11, workbookID12
// }

func Test_statRepository_FindStat(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	logrus.Warnf("today: %v", today)

	fn := func(ctx context.Context, ts testService) {
		_, sysOwner, owner := testInitOrganization(t, ts)
		workbookRepo, _ := ts.rf.NewWorkbookRepository(ctx)

		user1 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_1", "USERNAME_1")
		user2 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_2", "USERNAME_2")

		// user1 has two workbooks(WB11, WB12)
		student1 := testNewStudent(t, ts, user1)
		space1, _ := student1.GetPersonalSpace(ctx)
		workbook11 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB11")
		testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB12")

		// user2 has one workbook(WB21)
		student2 := testNewStudent(t, ts, user2)
		space2, _ := student2.GetPersonalSpace(ctx)
		testNewWorkbook(t, ctx, ts.db, workbookRepo, student2, userD.SpaceID(space2.GetID()), "WB21")

		// yesterday
		ts.db.Exec("INSERT INTO study_stat (app_user_id, workbook_id, problem_type_id, study_type_id, answered, mastered, record_date) values(?, ?, ?, ?, ?, ?, ?)", user1.GetID(), workbook11.GetID(), 1, 1, 10, 20, today.AddDate(0, 0, -1))
		ts.db.Exec("INSERT INTO study_stat (app_user_id, workbook_id, problem_type_id, study_type_id, answered, mastered, record_date) values(?, ?, ?, ?, ?, ?, ?)", user1.GetID(), workbook11.GetID(), 1, 2, 11, 21, today.AddDate(0, 0, -1))
		// two days ago
		ts.db.Exec("INSERT INTO study_stat (app_user_id, workbook_id, problem_type_id, study_type_id, answered, mastered, record_date) values(?, ?, ?, ?, ?, ?, ?)", user1.GetID(), workbook11.GetID(), 1, 2, 12, 22, today.AddDate(0, 0, -2))

		statRepo := gateway.NewStatRepository(ctx, ts.db)
		stat, err := statRepo.FindStat(ctx, userD.AppUserID(user1.GetID()))
		assert.NoError(t, err)
		assert.Equal(t, stat.GetUserID(), userD.AppUserID(user1.GetID()))
		// yesterday
		assert.Equal(t, stat.GetHistory().Results[6].Date.Format(time.RFC3339), today.AddDate(0, 0, -1).Format(time.RFC3339))
		assert.Equal(t, stat.GetHistory().Results[6].Answered, 21)
		assert.Equal(t, stat.GetHistory().Results[6].Mastered, 41)
		// two days ago
		assert.Equal(t, stat.GetHistory().Results[5].Date.Format(time.RFC3339), today.AddDate(0, 0, -2).Format(time.RFC3339))
		assert.Equal(t, stat.GetHistory().Results[5].Answered, 12)
		assert.Equal(t, stat.GetHistory().Results[5].Mastered, 22)
	}
	testDB(t, fn)
}
