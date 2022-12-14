// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	casbin "github.com/casbin/casbin/v2"
	domain "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// RBACRepository is an autogenerated mock type for the RBACRepository type
type RBACRepository struct {
	mock.Mock
}

// AddNamedGroupingPolicy provides a mock function with given fields: subject, object
func (_m *RBACRepository) AddNamedGroupingPolicy(subject domain.RBACUser, object domain.RBACRole) error {
	ret := _m.Called(subject, object)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.RBACUser, domain.RBACRole) error); ok {
		r0 = rf(subject, object)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddNamedPolicy provides a mock function with given fields: subject, object, action
func (_m *RBACRepository) AddNamedPolicy(subject domain.RBACRole, object domain.RBACObject, action domain.RBACAction) error {
	ret := _m.Called(subject, object, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.RBACRole, domain.RBACObject, domain.RBACAction) error); ok {
		r0 = rf(subject, object, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields:
func (_m *RBACRepository) Init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEnforcerWithRolesAndUsers provides a mock function with given fields: roles, users
func (_m *RBACRepository) NewEnforcerWithRolesAndUsers(roles []domain.RBACRole, users []domain.RBACUser) (*casbin.Enforcer, error) {
	ret := _m.Called(roles, users)

	var r0 *casbin.Enforcer
	if rf, ok := ret.Get(0).(func([]domain.RBACRole, []domain.RBACUser) *casbin.Enforcer); ok {
		r0 = rf(roles, users)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*casbin.Enforcer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]domain.RBACRole, []domain.RBACUser) error); ok {
		r1 = rf(roles, users)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRBACRepository creates a new instance of RBACRepository. It also registers a cleanup function to assert the mocks expectations.
func NewRBACRepository(t testing.TB) *RBACRepository {
	mock := &RBACRepository{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
