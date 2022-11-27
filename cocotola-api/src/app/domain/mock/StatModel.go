// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	userdomain "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
)

// StatModel is an autogenerated mock type for the StatModel type
type StatModel struct {
	mock.Mock
}

// GetHistory provides a mock function with given fields:
func (_m *StatModel) GetHistory() domain.StatHistory {
	ret := _m.Called()

	var r0 domain.StatHistory
	if rf, ok := ret.Get(0).(func() domain.StatHistory); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.StatHistory)
	}

	return r0
}

// GetUserID provides a mock function with given fields:
func (_m *StatModel) GetUserID() userdomain.AppUserID {
	ret := _m.Called()

	var r0 userdomain.AppUserID
	if rf, ok := ret.Get(0).(func() userdomain.AppUserID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(userdomain.AppUserID)
	}

	return r0
}

// NewStatModel creates a new instance of StatModel. It also registers a cleanup function to assert the mocks expectations.
func NewStatModel(t testing.TB) *StatModel {
	mock := &StatModel{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
