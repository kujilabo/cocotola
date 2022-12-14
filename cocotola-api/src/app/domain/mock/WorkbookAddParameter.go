// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// WorkbookAddParameter is an autogenerated mock type for the WorkbookAddParameter type
type WorkbookAddParameter struct {
	mock.Mock
}

// GetLang2 provides a mock function with given fields:
func (_m *WorkbookAddParameter) GetLang2() domain.Lang2 {
	ret := _m.Called()

	var r0 domain.Lang2
	if rf, ok := ret.Get(0).(func() domain.Lang2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Lang2)
		}
	}

	return r0
}

// GetName provides a mock function with given fields:
func (_m *WorkbookAddParameter) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetProblemType provides a mock function with given fields:
func (_m *WorkbookAddParameter) GetProblemType() domain.ProblemTypeName {
	ret := _m.Called()

	var r0 domain.ProblemTypeName
	if rf, ok := ret.Get(0).(func() domain.ProblemTypeName); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.ProblemTypeName)
	}

	return r0
}

// GetProperties provides a mock function with given fields:
func (_m *WorkbookAddParameter) GetProperties() map[string]string {
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

// GetQuestionText provides a mock function with given fields:
func (_m *WorkbookAddParameter) GetQuestionText() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewWorkbookAddParameter creates a new instance of WorkbookAddParameter. It also registers a cleanup function to assert the mocks expectations.
func NewWorkbookAddParameter(t testing.TB) *WorkbookAddParameter {
	mock := &WorkbookAddParameter{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
