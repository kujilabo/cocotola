// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	time "time"
)

// ProblemModel is an autogenerated mock type for the ProblemModel type
type ProblemModel struct {
	mock.Mock
}

// GetCreatedAt provides a mock function with given fields:
func (_m *ProblemModel) GetCreatedAt() time.Time {
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
func (_m *ProblemModel) GetCreatedBy() uint {
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
func (_m *ProblemModel) GetID() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetNumber provides a mock function with given fields:
func (_m *ProblemModel) GetNumber() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetProblemType provides a mock function with given fields:
func (_m *ProblemModel) GetProblemType() domain.ProblemTypeName {
	ret := _m.Called()

	var r0 domain.ProblemTypeName
	if rf, ok := ret.Get(0).(func() domain.ProblemTypeName); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.ProblemTypeName)
	}

	return r0
}

// GetProperties provides a mock function with given fields: ctx
func (_m *ProblemModel) GetProperties(ctx context.Context) map[string]interface{} {
	ret := _m.Called(ctx)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]interface{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// GetUpdatedAt provides a mock function with given fields:
func (_m *ProblemModel) GetUpdatedAt() time.Time {
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
func (_m *ProblemModel) GetUpdatedBy() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// GetVersion provides a mock function with given fields:
func (_m *ProblemModel) GetVersion() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// NewProblemModel creates a new instance of ProblemModel. It also registers a cleanup function to assert the mocks expectations.
func NewProblemModel(t testing.TB) *ProblemModel {
	mock := &ProblemModel{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
