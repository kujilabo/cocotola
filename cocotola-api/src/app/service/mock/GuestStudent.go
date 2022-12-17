// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola/cocotola-api/src/app/service"

	testing "testing"

	userservice "github.com/kujilabo/cocotola/cocotola-api/src/user/service"
)

// GuestStudent is an autogenerated mock type for the GuestStudent type
type GuestStudent struct {
	mock.Mock
}

// FindWorkbooksFromPublicSpace provides a mock function with given fields: ctx, condition
func (_m *GuestStudent) FindWorkbooksFromPublicSpace(ctx context.Context, condition service.WorkbookSearchCondition) (service.WorkbookSearchResult, error) {
	ret := _m.Called(ctx, condition)

	var r0 service.WorkbookSearchResult
	if rf, ok := ret.Get(0).(func(context.Context, service.WorkbookSearchCondition) service.WorkbookSearchResult); ok {
		r0 = rf(ctx, condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.WorkbookSearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.WorkbookSearchCondition) error); ok {
		r1 = rf(ctx, condition)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAppUserID provides a mock function with given fields:
func (_m *GuestStudent) GetAppUserID() domain.AppUserID {
	ret := _m.Called()

	var r0 domain.AppUserID
	if rf, ok := ret.Get(0).(func() domain.AppUserID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.AppUserID)
	}

	return r0
}

// GetDefaultSpace provides a mock function with given fields: ctx
func (_m *GuestStudent) GetDefaultSpace(ctx context.Context) (userservice.Space, error) {
	ret := _m.Called(ctx)

	var r0 userservice.Space
	if rf, ok := ret.Get(0).(func(context.Context) userservice.Space); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.Space)
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

// GetID provides a mock function with given fields:
func (_m *GuestStudent) GetID() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetLoginID provides a mock function with given fields:
func (_m *GuestStudent) GetLoginID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetOrganizationID provides a mock function with given fields:
func (_m *GuestStudent) GetOrganizationID() domain.OrganizationID {
	ret := _m.Called()

	var r0 domain.OrganizationID
	if rf, ok := ret.Get(0).(func() domain.OrganizationID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.OrganizationID)
	}

	return r0
}

// GetProperties provides a mock function with given fields:
func (_m *GuestStudent) GetProperties() map[string]string {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// GetRoles provides a mock function with given fields:
func (_m *GuestStudent) GetRoles() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetUsername provides a mock function with given fields:
func (_m *GuestStudent) GetUsername() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsStudentModel provides a mock function with given fields:
func (_m *GuestStudent) IsStudentModel() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewGuestStudent creates a new instance of GuestStudent. It also registers a cleanup function to assert the mocks expectations.
func NewGuestStudent(t testing.TB) *GuestStudent {
	mock := &GuestStudent{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
