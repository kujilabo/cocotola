// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	appservice "github.com/kujilabo/cocotola/cocotola-api/src/app/service"

	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola/cocotola-api/src/job/service"

	testing "testing"

	userservice "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

// RepositoryFactory is an autogenerated mock type for the RepositoryFactory type
type RepositoryFactory struct {
	mock.Mock
}

// NewJobRepositoryFactory provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewJobRepositoryFactory(ctx context.Context) (service.RepositoryFactory, error) {
	ret := _m.Called(ctx)

	var r0 service.RepositoryFactory
	if rf, ok := ret.Get(0).(func(context.Context) service.RepositoryFactory); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.RepositoryFactory)
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

// NewProblemRepository provides a mock function with given fields: ctx, problemType
func (_m *RepositoryFactory) NewProblemRepository(ctx context.Context, problemType string) (appservice.ProblemRepository, error) {
	ret := _m.Called(ctx, problemType)

	var r0 appservice.ProblemRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) appservice.ProblemRepository); ok {
		r0 = rf(ctx, problemType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.ProblemRepository)
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
func (_m *RepositoryFactory) NewProblemTypeRepository(ctx context.Context) appservice.ProblemTypeRepository {
	ret := _m.Called(ctx)

	var r0 appservice.ProblemTypeRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.ProblemTypeRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.ProblemTypeRepository)
		}
	}

	return r0
}

// NewRecordbookRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewRecordbookRepository(ctx context.Context) appservice.RecordbookRepository {
	ret := _m.Called(ctx)

	var r0 appservice.RecordbookRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.RecordbookRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.RecordbookRepository)
		}
	}

	return r0
}

// NewStatRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStatRepository(ctx context.Context) appservice.StatRepository {
	ret := _m.Called(ctx)

	var r0 appservice.StatRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.StatRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.StatRepository)
		}
	}

	return r0
}

// NewStudyRecordRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStudyRecordRepository(ctx context.Context) appservice.StudyRecordRepository {
	ret := _m.Called(ctx)

	var r0 appservice.StudyRecordRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.StudyRecordRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.StudyRecordRepository)
		}
	}

	return r0
}

// NewStudyStatRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStudyStatRepository(ctx context.Context) appservice.StudyStatRepository {
	ret := _m.Called(ctx)

	var r0 appservice.StudyStatRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.StudyStatRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.StudyStatRepository)
		}
	}

	return r0
}

// NewStudyTypeRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewStudyTypeRepository(ctx context.Context) appservice.StudyTypeRepository {
	ret := _m.Called(ctx)

	var r0 appservice.StudyTypeRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.StudyTypeRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.StudyTypeRepository)
		}
	}

	return r0
}

// NewUserQuotaRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewUserQuotaRepository(ctx context.Context) appservice.UserQuotaRepository {
	ret := _m.Called(ctx)

	var r0 appservice.UserQuotaRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.UserQuotaRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.UserQuotaRepository)
		}
	}

	return r0
}

// NewUserRepositoryFactory provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewUserRepositoryFactory(ctx context.Context) (userservice.RepositoryFactory, error) {
	ret := _m.Called(ctx)

	var r0 userservice.RepositoryFactory
	if rf, ok := ret.Get(0).(func(context.Context) userservice.RepositoryFactory); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.RepositoryFactory)
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
func (_m *RepositoryFactory) NewWorkbookRepository(ctx context.Context) appservice.WorkbookRepository {
	ret := _m.Called(ctx)

	var r0 appservice.WorkbookRepository
	if rf, ok := ret.Get(0).(func(context.Context) appservice.WorkbookRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appservice.WorkbookRepository)
		}
	}

	return r0
}

// NewRepositoryFactory creates a new instance of RepositoryFactory. It also registers a cleanup function to assert the mocks expectations.
func NewRepositoryFactory(t testing.TB) *RepositoryFactory {
	mock := &RepositoryFactory{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
