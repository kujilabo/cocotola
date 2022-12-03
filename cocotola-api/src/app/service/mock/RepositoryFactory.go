// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	service "github.com/kujilabo/cocotola/cocotola-api/src/app/service"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// RepositoryFactory is an autogenerated mock type for the RepositoryFactory type
type RepositoryFactory struct {
	mock.Mock
}

// NewProblemRepository provides a mock function with given fields: ctx, problemType
func (_m *RepositoryFactory) NewProblemRepository(ctx context.Context, problemType string) (service.ProblemRepository, error) {
	ret := _m.Called(ctx, problemType)

	var r0 service.ProblemRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) service.ProblemRepository); ok {
		r0 = rf(ctx, problemType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.ProblemRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, problemType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProblemTypeRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewProblemTypeRepository(ctx context.Context) (service.ProblemTypeRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.ProblemTypeRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.ProblemTypeRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.ProblemTypeRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRecordbookRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewRecordbookRepository(ctx context.Context) (service.RecordbookRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.RecordbookRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.RecordbookRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.RecordbookRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStatRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStatRepository(ctx context.Context) (service.StatRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.StatRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.StatRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.StatRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStudyRecordRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStudyRecordRepository(ctx context.Context) (service.StudyRecordRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.StudyRecordRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.StudyRecordRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.StudyRecordRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStudyStatRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStudyStatRepository(ctx context.Context) (service.StudyStatRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.StudyStatRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.StudyStatRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.StudyStatRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStudyTypeRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStudyTypeRepository(ctx context.Context) (service.StudyTypeRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.StudyTypeRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.StudyTypeRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.StudyTypeRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserQuotaRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewUserQuotaRepository(ctx context.Context) (service.UserQuotaRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.UserQuotaRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.UserQuotaRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.UserQuotaRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWorkbookRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewWorkbookRepository(ctx context.Context) (service.WorkbookRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.WorkbookRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.WorkbookRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.WorkbookRepository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepositoryFactory creates a new instance of RepositoryFactory. It also registers a cleanup function to assert the mocks expectations.
func NewRepositoryFactory(t testing.TB) *RepositoryFactory {
	mock := &RepositoryFactory{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
