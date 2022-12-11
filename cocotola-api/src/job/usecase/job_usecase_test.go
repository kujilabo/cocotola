package usecase_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	service_mock "github.com/kujilabo/cocotola/cocotola-api/src/job/service/mock"
	"github.com/kujilabo/cocotola/cocotola-api/src/job/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

type transaction struct {
	rf service.RepositoryFactory
}

func newTransaction(rf service.RepositoryFactory) service.Transaction {
	return &transaction{
		rf: rf,
	}
}

func (t *transaction) Do(ctx context.Context, fn func(rf service.RepositoryFactory) error) error {
	return fn(t.rf)
}

func Test_StartJob_timedout(t *testing.T) {
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	var value int32
	job, err := service.NewJob(domain.JobName("job1"), time.Duration(2)*time.Second, false, func(ctx context.Context) error {
		defer wg.Done()
		for i := 0; i < 4; i++ {
			time.Sleep(time.Second)
			select {
			case <-ctx.Done():
				logrus.Info("canceled")
				return ctx.Err()
			default:
				logrus.Info("continue")
			}
		}
		logrus.Info("completed")
		atomic.AddInt32(&value, 1)
		return nil
	})
	require.NoError(t, err)

	jobStatusRepo := new(service_mock.JobStatusRepository)
	jobStatusRepo.On("AddJobStatus", anythingOfContext, job).Return(domain.JobStatusID("StatusID1"), nil)
	jobStatusRepo.On("RemoveJobStatus", anythingOfContext, domain.JobStatusID("StatusID1")).Return(nil)
	jobHistoryRepo := new(service_mock.JobHistoryRepository)
	jobHistoryRepo.On("AddJobHistory", anythingOfContext, mock.Anything).Return(nil)
	rf := new(service_mock.RepositoryFactory)
	rf.On("NewJobStatusRepository", anythingOfContext).Return(jobStatusRepo, nil)
	rf.On("NewJobHistoryRepository", anythingOfContext).Return(jobHistoryRepo, nil)
	transaction := newTransaction(rf)
	jobUsecase, err := usecase.NewJobUsecase(ctx, transaction)
	require.NoError(t, err)

	err = jobUsecase.StartJob(ctx, job)
	assert.NoError(t, err)
	wg.Wait()
	assert.Equal(t, int32(0), value)
}

func Test_StartJob_completed(t *testing.T) {
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	var value int32
	job, err := service.NewJob(domain.JobName("job1"), time.Duration(2)*time.Second, false, func(ctx context.Context) error {
		defer wg.Done()
		for i := 0; i < 1; i++ {
			time.Sleep(time.Second)
			select {
			case <-ctx.Done():
				logrus.Info("canceled")
				return ctx.Err()
			default:
				logrus.Info("continue")
			}
		}
		logrus.Info("completed")
		atomic.AddInt32(&value, 1)
		return nil
	})
	require.NoError(t, err)

	jobStatusRepo := new(service_mock.JobStatusRepository)
	jobStatusRepo.On("AddJobStatus", anythingOfContext, job).Return(domain.JobStatusID("StatusID1"), nil)
	jobStatusRepo.On("RemoveJobStatus", anythingOfContext, domain.JobStatusID("StatusID1")).Return(nil)
	jobHistoryRepo := new(service_mock.JobHistoryRepository)
	jobHistoryRepo.On("AddJobHistory", anythingOfContext, mock.Anything).Return(nil)
	rf := new(service_mock.RepositoryFactory)
	rf.On("NewJobStatusRepository", anythingOfContext).Return(jobStatusRepo, nil)
	rf.On("NewJobHistoryRepository", anythingOfContext).Return(jobHistoryRepo, nil)
	transaction := newTransaction(rf)
	jobUsecase, err := usecase.NewJobUsecase(ctx, transaction)
	require.NoError(t, err)

	err = jobUsecase.StartJob(ctx, job)
	assert.NoError(t, err)
	wg.Wait()
	assert.Equal(t, int32(1), value)
}
