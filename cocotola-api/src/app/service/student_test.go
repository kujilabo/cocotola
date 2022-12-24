//go:build s

package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	service_mock "github.com/kujilabo/cocotola/cocotola-api/src/app/service/mock"
	userD "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	userD_mock "github.com/kujilabo/cocotola/cocotola-api/src/user/domain/mock"
	userS_mock "github.com/kujilabo/cocotola/cocotola-api/src/user/service/mock"
)

const (
	problemType1          = domain.ProblemTypeName("PROBLEM_TYPE_1")
	problemType2          = domain.ProblemTypeName("PROBLEM_TYPE_2")
	studentLoginIDLength  = 20
	studentUsernameLength = 20
)

func newRandStudent(t *testing.T, orgID userD.OrganizationID) domain.StudentModel {

	model, err := userD.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)
	userModel, err := userD.NewAppUserModel(model, orgID, RandString(studentLoginIDLength), RandString(studentUsernameLength), []string{}, map[string]string{})
	require.NoError(t, err)
	studentModel, err := domain.NewStudentModel(userModel)
	require.NoError(t, err)
	return studentModel
}

func student_Init(t *testing.T, ctx context.Context) (
	pf *service_mock.ProcessorFactory,
	rf *service_mock.RepositoryFactory,
	spaceRepo *userS_mock.SpaceRepository,
	workbookRepo *service_mock.WorkbookRepository,
	userQuotaRepo *service_mock.UserQuotaRepository,
	problemQuotaProcessor *service_mock.ProblemQuotaProcessor) {

	userRf := new(userS_mock.RepositoryFactory)
	spaceRepo = new(userS_mock.SpaceRepository)
	userRf.On("NewSpaceRepository", ctx).Return(spaceRepo, nil)

	workbookRepo = new(service_mock.WorkbookRepository)
	userQuotaRepo = new(service_mock.UserQuotaRepository)
	rf = new(service_mock.RepositoryFactory)
	rf.On("NewWorkbookRepository", ctx).Return(workbookRepo, nil)
	rf.On("NewUserQuotaRepository", ctx).Return(userQuotaRepo, nil)
	rf.On("NewUserRepositoryFactory", ctx).Return(userRf, nil)

	problemQuotaProcessor = new(service_mock.ProblemQuotaProcessor)
	pf = new(service_mock.ProcessorFactory)
	pf.On("NewProblemQuotaProcessor", problemType1).Return(problemQuotaProcessor, nil)
	pf.On("NewProblemQuotaProcessor", problemType2).Return(problemQuotaProcessor, nil)
	return
}

func Test_student_GetDefaultSpace(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pf, rf, spaceRepo, _, _, _ := student_Init(t, ctx)

	expected := new(userD_mock.SpaceModel)
	spaceRepo.On("FindDefaultSpace", ctx, mock.Anything).Return(expected, nil)
	studentModel := newRandStudent(t, userD.OrganizationID(1))
	student, err := service.NewStudent(ctx, pf, rf, studentModel)
	require.NoError(t, err)
	// given
	expected.On("GetKey").Return("KEY")
	// when
	actual, err := student.GetDefaultSpace(ctx)
	require.NoError(t, err)
	// then
	require.Equal(t, "KEY", actual.GetKey())
	spaceRepo.AssertCalled(t, "FindDefaultSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindDefaultSpace", 1)
}

func Test_student_GetPersonalSpace(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pf, rf, spaceRepo, _, _, _ := student_Init(t, ctx)

	expected := new(userD_mock.SpaceModel)
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(expected, nil)
	studentModel := newRandStudent(t, userD.OrganizationID(1))
	student, err := service.NewStudent(ctx, pf, rf, studentModel)
	require.NoError(t, err)
	// given
	expected.On("GetKey").Return("KEY")
	// when
	actual, err := student.GetPersonalSpace(ctx)
	require.NoError(t, err)
	// then
	require.Equal(t, "KEY", actual.GetKey())
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbooksFromPersonalSpace(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pf, rf, spaceRepo, workbookRepo, _, _ := student_Init(t, ctx)

	space := new(userD_mock.SpaceModel)
	space.On("GetID").Return(uint(100))
	spaceRepo.On("FindPersonalSpace", ctx, mock.Anything).Return(space, nil)

	studentModel := newRandStudent(t, userD.OrganizationID(1))
	student, err := service.NewStudent(ctx, pf, rf, studentModel)
	require.NoError(t, err)
	// given
	expected, err := service.NewWorkbookSearchResult(123, nil)
	require.NoError(t, err)
	workbookRepo.On("FindPersonalWorkbooks", ctx, mock.Anything, mock.Anything).Return(expected, nil)
	// when
	condition, err := service.NewWorkbookSearchCondition(1, 100, nil)
	require.NoError(t, err)
	actual, err := student.FindWorkbooksFromPersonalSpace(ctx, condition)
	require.NoError(t, err)
	// then
	require.Equal(t, 123, actual.GetTotalCount())
	spaceRepo.AssertCalled(t, "FindPersonalSpace", ctx, mock.Anything)
	spaceRepo.AssertNumberOfCalls(t, "FindPersonalSpace", 1)
}

func Test_student_FindWorkbookByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pf, rf, _, workbookRepo, _, _ := student_Init(t, ctx)

	expected := new(service_mock.Workbook)
	workbookRepo.On("FindWorkbookByID", ctx, mock.Anything, mock.Anything).Return(expected, nil)

	studentModel := newRandStudent(t, userD.OrganizationID(1))
	student, err := service.NewStudent(ctx, pf, rf, studentModel)
	require.NoError(t, err)
	// given
	expected.On("GetID").Return(uint(123))
	// when
	actual, err := student.FindWorkbookByID(ctx, domain.WorkbookID(100))
	require.NoError(t, err)
	// then
	require.Equal(t, uint(123), actual.GetID())
}

func Test_student_CheckQuota(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		problemType domain.ProblemTypeName
		name        service.QuotaName
	}
	tests := []struct {
		name              string
		isExceeded        bool
		problemTypeSuffix string
		quotaUnit         service.QuotaUnit
		quotaLimit        int
		args              args
		err               error
	}{
		{
			name:              "QuotaNameSize,isNotExceeded",
			isExceeded:        false,
			problemTypeSuffix: "_size",
			quotaUnit:         service.QuotaUnitPersitance,
			quotaLimit:        234,
			args: args{
				problemType: problemType1,
				name:        service.QuotaNameSize,
			},
			err: nil,
		},
		{
			name:              "QuotaNameSize,isExceeded",
			isExceeded:        true,
			problemTypeSuffix: "_size",
			quotaUnit:         service.QuotaUnitPersitance,
			quotaLimit:        234,
			args: args{
				problemType: problemType2,
				name:        service.QuotaNameSize,
			},
			err: service.ErrQuotaExceeded,
		},
		{
			name:              "QuotaNameUpdate,isNotExceeded",
			isExceeded:        false,
			problemTypeSuffix: "_update",
			quotaUnit:         service.QuotaUnitDay,
			quotaLimit:        345,
			args: args{
				problemType: problemType1,
				name:        service.QuotaNameUpdate,
			},
			err: nil,
		},
		{
			name:              "QuotaNameUpdate,isExceeded",
			isExceeded:        true,
			problemTypeSuffix: "_update",
			quotaUnit:         service.QuotaUnitDay,
			quotaLimit:        345,
			args: args{
				problemType: problemType2,
				name:        service.QuotaNameUpdate,
			},
			err: service.ErrQuotaExceeded,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pf, rf, _, _, userQuotaRepo, problemQuotaProcessor := student_Init(t, ctx)
			userQuotaRepo.On("IsExceeded", mock.Anything, mock.Anything, mock.Anything, string(tt.args.problemType)+tt.problemTypeSuffix, tt.quotaUnit, tt.quotaLimit).Return(tt.isExceeded, nil)
			problemQuotaProcessor.On("GetUnitForSizeQuota").Return(service.QuotaUnitPersitance)
			problemQuotaProcessor.On("GetLimitForSizeQuota").Return(tt.quotaLimit)
			problemQuotaProcessor.On("GetUnitForUpdateQuota").Return(service.QuotaUnitDay)
			problemQuotaProcessor.On("GetLimitForUpdateQuota").Return(tt.quotaLimit)

			studentModel := newRandStudent(t, userD.OrganizationID(1))

			student, err := service.NewStudent(ctx, pf, rf, studentModel)
			require.NoError(t, err)
			require.NotNil(t, student)
			err = student.CheckQuota(ctx, tt.args.problemType, tt.args.name)
			if err == nil && tt.err != nil {
				t.Errorf("student.CheckQuota() error = %v, err %v", err, tt.err)
			} else if err != nil && tt.err == nil {
				t.Errorf("student.CheckQuota() error = %v, err %v", err, tt.err)
			} else if err != nil && tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("student.CheckQuota() error = %v, err %v", err, tt.err)
			}
			userQuotaRepo.AssertCalled(t, "IsExceeded", mock.Anything, mock.Anything, mock.Anything, string(tt.args.problemType)+tt.problemTypeSuffix, tt.quotaUnit, tt.quotaLimit)
			userQuotaRepo.AssertNumberOfCalls(t, "IsExceeded", 1)
		})
	}
}
