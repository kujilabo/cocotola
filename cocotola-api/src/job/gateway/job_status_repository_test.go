package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
	"github.com/kujilabo/cocotola/cocotola-api/src/job/gateway"
	"github.com/kujilabo/cocotola/cocotola-api/src/job/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func emptyFunc(ctx context.Context) error {
	return nil
}

func Test_jobStatusRepository_AddJobStatus_allowedConcurrencyExecution_is_true(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(ctx context.Context, ts testService) {
		setupJob(t, ts)
		defer teardownJob(t, ts)
		jobStatusRepo, err := gateway.NewJobStatusRepository(ctx, ts.db)
		require.NoError(t, err)
		job, err := service.NewJob(domain.JobName("job"), time.Second, true, emptyFunc)
		require.NoError(t, err)
		jobStatusID1, err1 := jobStatusRepo.AddJobStatus(ctx, job)
		require.NoError(t, err1)
		require.Equal(t, 36, len(jobStatusID1))
		require.NoError(t, err1)
		jobStatusID2, err2 := jobStatusRepo.AddJobStatus(ctx, job)
		require.NoError(t, err2)
		require.Equal(t, 36, len(jobStatusID2))
	}
	testDB(t, fn)
}

func Test_jobStatusRepository_AddJobStatus_allowedConcurrencyExecution_is_false(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)

	fn := func(ctx context.Context, ts testService) {
		setupJob(t, ts)
		defer teardownJob(t, ts)
		jobStatusRepo, err := gateway.NewJobStatusRepository(ctx, ts.db)
		require.NoError(t, err)
		job, err := service.NewJob(domain.JobName("job"), time.Second, false, emptyFunc)
		require.NoError(t, err)
		jobStatusID1, err1 := jobStatusRepo.AddJobStatus(ctx, job)
		require.NoError(t, err1)
		require.Equal(t, 36, len(jobStatusID1))
		jobStatusID2, err2 := jobStatusRepo.AddJobStatus(ctx, job)
		require.Error(t, err2)
		require.Equal(t, 0, len(jobStatusID2))
	}
	testDB(t, fn)
}

func Test_jobStatusRepository_RemoveExpiredJobStatus(t *testing.T) {
	fn := func(ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		setupJob(t, ts)
		defer teardownJob(t, ts)
		jobStatusRepo, err := gateway.NewJobStatusRepository(ctx, ts.db)
		require.NoError(t, err)

		job1, err := service.NewJob(domain.JobName("job1"), time.Second, false, emptyFunc)
		require.NoError(t, err)
		jobStatusID1, err := jobStatusRepo.AddJobStatus(ctx, job1)
		require.NoError(t, err)
		require.Equal(t, 36, len(jobStatusID1))

		job2, err := service.NewJob(domain.JobName("job2"), time.Duration(3)*time.Second, false, emptyFunc)
		require.NoError(t, err)
		jobStatusID2, err := jobStatusRepo.AddJobStatus(ctx, job2)
		require.NoError(t, err)
		require.Equal(t, 36, len(jobStatusID2))

		time.Sleep(time.Duration(2) * time.Second)

		deleted, err := jobStatusRepo.RemoveExpiredJobStatus(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, deleted)

		// job1 was deleted
		foundJob1, err := jobStatusRepo.FindJobStatusByJobName(ctx, domain.JobName("job1"))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(foundJob1))

		// job2 was not deleted
		foundJob2, err := jobStatusRepo.FindJobStatusByJobName(ctx, domain.JobName("job2"))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(foundJob2))
	}
	testDB(t, fn)
}
