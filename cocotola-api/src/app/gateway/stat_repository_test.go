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
	libD "github.com/kujilabo/cocotola/lib/domain"
)

func Test_statRepository_FindStat(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// logrus.Warnf("today: %v", today)

	logrus.SetLevel(logrus.WarnLevel)
	fn := func(ctx context.Context, ts testService) {
		logrus.SetLevel(logrus.DebugLevel)
		orgID, sysOwner, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		workbookRepo, _ := ts.rf.NewWorkbookRepository(ctx)

		user1 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_1", "USERNAME_1")
		user2 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_2", "USERNAME_2")

		// user1 has two workbooks(WB11, WB12)
		student1 := testNewStudent(ctx, t, ts, user1)
		space1, _ := student1.GetPersonalSpace(ctx)
		workbook11 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB11")
		testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB12")

		// user2 has one workbook(WB21)
		student2 := testNewStudent(ctx, t, ts, user2)
		space2, _ := student2.GetPersonalSpace(ctx)
		testNewWorkbook(t, ctx, ts.db, workbookRepo, student2, userD.SpaceID(space2.GetID()), "WB21")

		// yesterday
		ts.db.Debug().Exec("INSERT INTO study_stat (id, organization_id, app_user_id, workbook_id, problem_type_id, study_type_id, answered, mastered, record_date) values(?, ?, ?, ?, ?, ?, ?, ?, ?)", libD.NewULID(), uint(orgID), user1.GetID(), workbook11.GetID(), 1, 1, 10, 20, today.AddDate(0, 0, -1))
		ts.db.Debug().Exec("INSERT INTO study_stat (id, organization_id, app_user_id, workbook_id, problem_type_id, study_type_id, answered, mastered, record_date) values(?, ?, ?, ?, ?, ?, ?, ?, ?)", libD.NewULID(), uint(orgID), user1.GetID(), workbook11.GetID(), 1, 2, 11, 21, today.AddDate(0, 0, -1))
		// two days ago
		ts.db.Debug().Exec("INSERT INTO study_stat (id, organization_id, app_user_id, workbook_id, problem_type_id, study_type_id, answered, mastered, record_date) values(?, ?, ?, ?, ?, ?, ?, ?, ?)", libD.NewULID(), uint(orgID), user1.GetID(), workbook11.GetID(), 1, 2, 12, 22, today.AddDate(0, 0, -2))

		statRepo, _ := gateway.NewStatRepository(ctx, ts.db)
		stat, err := statRepo.FindStat(ctx, userD.AppUserID(user1.GetID()))
		assert.NoError(t, err)

		// for _, s := range stat.GetHistory().Results {
		// 	logrus.Debugf("stat: %+v", s)
		// }
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
