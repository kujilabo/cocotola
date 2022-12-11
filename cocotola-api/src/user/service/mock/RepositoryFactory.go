// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	service "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// RepositoryFactory is an autogenerated mock type for the RepositoryFactory type
type RepositoryFactory struct {
	mock.Mock
}

// NewAppUserGroupRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewAppUserGroupRepository(ctx context.Context) (service.AppUserGroupRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.AppUserGroupRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.AppUserGroupRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.AppUserGroupRepository)
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

// NewAppUserRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewAppUserRepository(ctx context.Context) (service.AppUserRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.AppUserRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.AppUserRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.AppUserRepository)
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

// NewGroupUserRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewGroupUserRepository(ctx context.Context) (service.GroupUserRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.GroupUserRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.GroupUserRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.GroupUserRepository)
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

// NewOrganizationRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewOrganizationRepository(ctx context.Context) (service.OrganizationRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.OrganizationRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.OrganizationRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.OrganizationRepository)
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

// NewRBACRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewRBACRepository(ctx context.Context) (service.RBACRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.RBACRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.RBACRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.RBACRepository)
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

// NewSpaceRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewSpaceRepository(ctx context.Context) (service.SpaceRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.SpaceRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.SpaceRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.SpaceRepository)
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

// NewUserSpaceRepository provides a mock function with given fields: ctx
func (_m *RepositoryFactory) NewUserSpaceRepository(ctx context.Context) (service.UserSpaceRepository, error) {
	ret := _m.Called(ctx)

	var r0 service.UserSpaceRepository
	if rf, ok := ret.Get(0).(func(context.Context) service.UserSpaceRepository); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.UserSpaceRepository)
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
