package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_studyRecordRepository_CountAnsweredProblems(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	fn := func(ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, sysOwner, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		workbookRepo := ts.rf.NewWorkbookRepository(ctx)
		studyRecordRepo := ts.rf.NewStudyRecordRepository(ctx)

		user1 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_1", "USERNAME_1")
		user2 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_2", "USERNAME_2")

		// user1 has two workbooks(WB11, WB12)
		student1 := testNewStudent(ctx, t, ts, user1)
		space1, _ := student1.GetPersonalSpace(ctx)
		workbook11 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB11")
		workbook12 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB12")

		logrus.Debug("Test_recordbookRepository_CountAnsweredProblems. A")
		// workbok11, english-word, memorization
		for _, problemID := range []int{111, 112, 113, 114} {
			problemID := domain.ProblemID(problemID)
			err := studyRecordRepo.AddRecord(ctx, sysOwner, userD.AppUserID(student1.GetID()), workbook11.GetWorkbookID(), englishWordID, memorizationID, problemID, false)
			require.NoError(t, err)
		}
		// workbok11, english-word, memorization, mastered
		for _, problemID := range []int{115, 116, 117} {
			problemID := domain.ProblemID(problemID)
			err := studyRecordRepo.AddRecord(ctx, sysOwner, userD.AppUserID(student1.GetID()), workbook11.GetWorkbookID(), englishWordID, memorizationID, problemID, true)
			require.NoError(t, err)
		}
		// workbok11, english-word, dictation
		for _, problemID := range []int{111, 112, 113} {
			problemID := domain.ProblemID(problemID)
			err := studyRecordRepo.AddRecord(ctx, sysOwner, userD.AppUserID(student1.GetID()), workbook11.GetWorkbookID(), englishWordID, dictationID, problemID, false)
			require.NoError(t, err)
		}
		// workbok12, english-word, memorization
		for _, problemID := range []int{121, 122} {
			problemID := domain.ProblemID(problemID)
			err := studyRecordRepo.AddRecord(ctx, sysOwner, userD.AppUserID(student1.GetID()), workbook12.GetWorkbookID(), englishWordID, memorizationID, problemID, false)
			require.NoError(t, err)
		}

		// user2 has one workbook(WB21)
		student2 := testNewStudent(ctx, t, ts, user2)
		space2, _ := student2.GetPersonalSpace(ctx)
		workbook21 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student2, userD.SpaceID(space2.GetID()), "WB21")
		for _, problemID := range []int{211} {
			problemID := domain.ProblemID(problemID)
			err := studyRecordRepo.AddRecord(ctx, sysOwner, userD.AppUserID(student2.GetID()), workbook21.GetWorkbookID(), englishWordID, memorizationID, problemID, false)
			require.NoError(t, err)
		}

		results, err := studyRecordRepo.CountAnsweredProblems(ctx, userD.AppUserID(user1.GetID()), today)
		require.NoError(t, err)
		assert.Equal(t, 3, len(results.Results))
		for _, result := range results.Results {
			if result.WorkbookID == workbook11.GetID() &&
				result.ProblemTypeID == englishWordID &&
				result.StudyTypeID == memorizationID {
				assert.Equal(t, 7, result.Answered)
				assert.Equal(t, 3, result.Mastered)
			} else if result.WorkbookID == workbook11.GetID() &&
				result.ProblemTypeID == englishWordID &&
				result.StudyTypeID == dictationID {
				assert.Equal(t, 3, result.Answered)
				assert.Equal(t, 0, result.Mastered)
			} else if result.WorkbookID == workbook12.GetID() &&
				result.ProblemTypeID == englishWordID &&
				result.StudyTypeID == dictationID {
				assert.Equal(t, 2, result.Answered)
				assert.Equal(t, 0, result.Mastered)
			}
		}
	}
	testDB(t, fn)
}
