// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ProblemUpdateParameter is an autogenerated mock type for the ProblemUpdateParameter type
type ProblemUpdateParameter struct {
	mock.Mock
}

// GetIntProperty provides a mock function with given fields: name
func (_m *ProblemUpdateParameter) GetIntProperty(name string) (int, error) {
	ret := _m.Called(name)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNumber provides a mock function with given fields:
func (_m *ProblemUpdateParameter) GetNumber() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetProperties provides a mock function with given fields:
func (_m *ProblemUpdateParameter) GetProperties() map[string]string {
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

// GetStringProperty provides a mock function with given fields: name
func (_m *ProblemUpdateParameter) GetStringProperty(name string) (string, error) {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProblemUpdateParameter creates a new instance of ProblemUpdateParameter. It also registers a cleanup function to assert the mocks expectations.
func NewProblemUpdateParameter(t testing.TB) *ProblemUpdateParameter {
	mock := &ProblemUpdateParameter{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}