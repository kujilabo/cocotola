// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	service "github.com/kujilabo/cocotola/cocotola-translator-api/src/service"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

// Do provides a mock function with given fields: ctx, fn
func (_m *Transaction) Do(ctx context.Context, fn func(service.RepositoryFactory) error) error {
	ret := _m.Called(ctx, fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(service.RepositoryFactory) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransaction creates a new instance of Transaction. It also registers a cleanup function to assert the mocks expectations.
func NewTransaction(t testing.TB) *Transaction {
	mock := &Transaction{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
