// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	service "github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/service"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TatoebaSentenceAddParameterIterator is an autogenerated mock type for the TatoebaSentenceAddParameterIterator type
type TatoebaSentenceAddParameterIterator struct {
	mock.Mock
}

// Next provides a mock function with given fields: ctx
func (_m *TatoebaSentenceAddParameterIterator) Next(ctx context.Context) (service.TatoebaSentenceAddParameter, error) {
	ret := _m.Called(ctx)

	var r0 service.TatoebaSentenceAddParameter
	if rf, ok := ret.Get(0).(func(context.Context) service.TatoebaSentenceAddParameter); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.TatoebaSentenceAddParameter)
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

// NewTatoebaSentenceAddParameterIterator creates a new instance of TatoebaSentenceAddParameterIterator. It also registers a cleanup function to assert the mocks expectations.
func NewTatoebaSentenceAddParameterIterator(t testing.TB) *TatoebaSentenceAddParameterIterator {
	mock := &TatoebaSentenceAddParameterIterator{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
