// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	time "time"
)

// StudentModel is an autogenerated mock type for the StudentModel type
type StudentModel struct {
	mock.Mock
}

// GetAppUserID provides a mock function with given fields:
func (_m *StudentModel) GetAppUserID() domain.AppUserID {
	ret := _m.Called()

	var r0 domain.AppUserID
	if rf, ok := ret.Get(0).(func() domain.AppUserID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.AppUserID)
	}

	return r0
}

// GetCreatedAt provides a mock function with given fields:
func (_m *StudentModel) GetCreatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetCreatedBy provides a mock function with given fields:
func (_m *StudentModel) GetCreatedBy() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetID provides a mock function with given fields:
func (_m *StudentModel) GetID() uint {
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
func (_m *StudentModel) GetLoginID() string {
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
func (_m *StudentModel) GetOrganizationID() domain.OrganizationID {
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
func (_m *StudentModel) GetProperties() map[string]string {
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
func (_m *StudentModel) GetRoles() []string {
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

// GetUpdatedAt provides a mock function with given fields:
func (_m *StudentModel) GetUpdatedAt() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetUpdatedBy provides a mock function with given fields:
func (_m *StudentModel) GetUpdatedBy() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetUsername provides a mock function with given fields:
func (_m *StudentModel) GetUsername() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetVersion provides a mock function with given fields:
func (_m *StudentModel) GetVersion() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// IsStudentModel provides a mock function with given fields:
func (_m *StudentModel) IsStudentModel() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewStudentModel creates a new instance of StudentModel. It also registers a cleanup function to assert the mocks expectations.
func NewStudentModel(t testing.TB) *StudentModel {
	mock := &StudentModel{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
