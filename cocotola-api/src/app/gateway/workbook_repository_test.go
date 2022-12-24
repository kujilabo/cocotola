//go:build m

package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

func Test_workbookRepository_FindPersonalWorkbooks(t *testing.T) {
	t.Parallel()
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, sysOwner, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		workbookRepo := ts.rf.NewWorkbookRepository(ctx)

		// user1 has two workbooks
		user1 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_1", "USERNAME_1")
		student1 := testNewStudent(ctx, t, ts, user1)
		space1, _ := student1.GetPersonalSpace(ctx)
		workbook11 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB11")
		workbook12 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB12")

		// user2 has one workbook
		user2 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_2", "USERNAME_2")
		student2 := testNewStudent(ctx, t, ts, user2)
		space2, _ := student2.GetPersonalSpace(ctx)
		workbook21 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student2, userD.SpaceID(space2.GetID()), "WB21")

		type args struct {
			operator service.Student
			param    service.WorkbookSearchCondition
		}
		type want struct {
			workbookID   domain.WorkbookID
			workbookName string
		}
		tests := []struct {
			name    string
			args    args
			want    []want
			wantErr bool
		}{
			{
				name: "user1",
				args: args{
					operator: student1,
					param:    testNewWorkbookSearchCondition(t),
				},
				want: []want{
					{
						workbookID:   workbook11.GetWorkbookID(),
						workbookName: "WB11",
					},
					{
						workbookID:   workbook12.GetWorkbookID(),
						workbookName: "WB12",
					},
				},
			},
			{
				name: "user2",
				args: args{
					operator: student2,
					param:    testNewWorkbookSearchCondition(t),
				},
				want: []want{
					{
						workbookID:   workbook21.GetWorkbookID(),
						workbookName: "WB21",
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := workbookRepo.FindPersonalWorkbooks(ctx, tt.args.operator, tt.args.param)
				if (err != nil) != tt.wantErr {
					t.Errorf("workbookRepository.FindPersonalWorkbooks() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					assert.Len(t, got.GetResults(), len(tt.want))
					for i, want := range tt.want {
						assert.Equal(t, uint(want.workbookID), got.GetResults()[i].GetID())
						assert.Equal(t, want.workbookName, got.GetResults()[i].GetName())
					}
					assert.Equal(t, len(tt.want), got.GetTotalCount())
				}
			})
		}
	}

	testDB(t, fn)
}

func Test_workbookRepository_FindWorkbookByName(t *testing.T) {
	t.Parallel()
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, sysOwner, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		workbookRepo := ts.rf.NewWorkbookRepository(ctx)

		user1 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_1", "USERNAME_1")
		testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_2", "USERNAME_2")
		// user1 has two workbooks
		student1 := testNewStudent(ctx, t, ts, user1)
		space1, _ := student1.GetPersonalSpace(ctx)
		workbook11 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB11")
		testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB12")

		type args struct {
			operator service.Student
			param    string
		}
		type want struct {
			workbookID   domain.WorkbookID
			workbookName string
			audioEnabled string
		}
		tests := []struct {
			name    string
			args    args
			want    want
			wantErr bool
		}{
			{
				name: "user1",
				args: args{
					operator: student1,
					param:    "WB11",
				},
				want: want{
					workbookID:   workbook11.GetWorkbookID(),
					workbookName: "WB11",
					audioEnabled: "false",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				space, _ := tt.args.operator.GetPersonalSpace(ctx)
				got, err := workbookRepo.FindWorkbookByName(ctx, tt.args.operator, userD.SpaceID(space.GetID()), tt.args.param)
				if (err != nil) != tt.wantErr {
					t.Errorf("workbookRepository.FindWorkbookByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					assert.Equal(t, uint(tt.want.workbookID), got.GetID())
					assert.Equal(t, tt.want.workbookName, got.GetName())
					assert.Equal(t, tt.want.audioEnabled, got.GetProperties()["audioEnabled"])
				}
			})
		}
	}
	testDB(t, fn)
}

func Test_workbookRepository_FindWorkbookByID_priv(t *testing.T) {
	t.Parallel()
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, sysOwner, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		workbookRepo := ts.rf.NewWorkbookRepository(ctx)

		user1 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_1", "USERNAME_1")
		user2 := testNewAppUser(t, ctx, ts, sysOwner, owner, "LOGIN_ID_2", "USERNAME_2")

		// user1 has two workbooks(WB11, WB12)
		student1 := testNewStudent(ctx, t, ts, user1)
		space1, _ := student1.GetPersonalSpace(ctx)
		workbook11 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB11")
		workbook12 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student1, userD.SpaceID(space1.GetID()), "WB12")

		// user2 has two workbooks(WB11, WB12)
		student2 := testNewStudent(ctx, t, ts, user2)
		space2, _ := student2.GetPersonalSpace(ctx)
		workbook21 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student2, userD.SpaceID(space2.GetID()), "WB21")
		workbook22 := testNewWorkbook(t, ctx, ts.db, workbookRepo, student2, userD.SpaceID(space2.GetID()), "WB22")

		// user1 can read user1's workbooks(WB11, WB12)
		workbook11Tmp, err := workbookRepo.FindWorkbookByID(ctx, student1, workbook11.GetWorkbookID())
		assert.NoError(t, err)
		assert.Equal(t, workbook11Tmp.GetID(), workbook11.GetID())
		workbook12Tmp, err := workbookRepo.FindWorkbookByID(ctx, student1, workbook12.GetWorkbookID())
		assert.NoError(t, err)
		assert.Equal(t, workbook12Tmp.GetID(), workbook12.GetID())

		// user1 cannot read user2's workbooks(WB21, WB22)
		if _, err := workbookRepo.FindWorkbookByID(ctx, student1, workbook21.GetWorkbookID()); err != nil {
			assert.True(t, errors.Is(err, service.ErrWorkbookPermissionDenied))
		} else {
			assert.Fail(t, "err is nil")
		}
		if _, err := workbookRepo.FindWorkbookByID(ctx, student1, workbook22.GetWorkbookID()); err != nil {
			assert.True(t, errors.Is(err, service.ErrWorkbookPermissionDenied))
		} else {
			assert.Fail(t, "err is nil")
		}
	}
	testDB(t, fn)
}
