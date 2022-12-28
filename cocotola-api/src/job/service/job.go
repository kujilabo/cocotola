//go:generate mockery --output mock --name Job
package service

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola/cocotola-api/src/job/domain"
)

type Job interface {
	GetName() domain.JobName
	GetJobParameter() string
	GetTimeout() time.Duration
	IsAllowedConcurrentExecution() bool
	Run(ctx context.Context) error
}

type job struct {
	Name                       domain.JobName
	JobParameter               string
	Timeout                    time.Duration
	AllowedConcurrentExecution bool
	Func                       func(context.Context) error
}

type JobFunc func(context.Context) error

func NewJob(name domain.JobName, timeout time.Duration, allowedConcurrentExecution bool, fn JobFunc) (Job, error) {
	return &job{
		Name:                       name,
		Timeout:                    timeout,
		AllowedConcurrentExecution: allowedConcurrentExecution,
		Func:                       fn,
	}, nil
}

func (m *job) GetName() domain.JobName {
	return m.Name
}

func (m *job) GetJobParameter() string {
	return m.JobParameter
}

func (m *job) GetTimeout() time.Duration {
	return m.Timeout
}

func (m *job) IsAllowedConcurrentExecution() bool {
	return m.AllowedConcurrentExecution
}

func (m *job) Run(ctx context.Context) error {
	return m.Func(ctx)
}
