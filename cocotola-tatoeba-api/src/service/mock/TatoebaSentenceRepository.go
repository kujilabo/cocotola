// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	service "github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/service"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TatoebaSentenceRepository is an autogenerated mock type for the TatoebaSentenceRepository type
type TatoebaSentenceRepository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, param
func (_m *TatoebaSentenceRepository) Add(ctx context.Context, param service.TatoebaSentenceAddParameter) error {
	ret := _m.Called(ctx, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, service.TatoebaSentenceAddParameter) error); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainsSentenceBySentenceNumber provides a mock function with given fields: ctx, sentenceNumber
func (_m *TatoebaSentenceRepository) ContainsSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (bool, error) {
	ret := _m.Called(ctx, sentenceNumber)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, sentenceNumber)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, sentenceNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTatoebaSentenceBySentenceNumber provides a mock function with given fields: ctx, sentenceNumber
func (_m *TatoebaSentenceRepository) FindTatoebaSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (service.TatoebaSentence, error) {
	ret := _m.Called(ctx, sentenceNumber)

	var r0 service.TatoebaSentence
	if rf, ok := ret.Get(0).(func(context.Context, int) service.TatoebaSentence); ok {
		r0 = rf(ctx, sentenceNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.TatoebaSentence)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, sentenceNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTatoebaSentencePairs provides a mock function with given fields: ctx, param
func (_m *TatoebaSentenceRepository) FindTatoebaSentencePairs(ctx context.Context, param service.TatoebaSentenceSearchCondition) (service.TatoebaSentencePairSearchResult, error) {
	ret := _m.Called(ctx, param)

	var r0 service.TatoebaSentencePairSearchResult
	if rf, ok := ret.Get(0).(func(context.Context, service.TatoebaSentenceSearchCondition) service.TatoebaSentencePairSearchResult); ok {
		r0 = rf(ctx, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.TatoebaSentencePairSearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.TatoebaSentenceSearchCondition) error); ok {
		r1 = rf(ctx, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTatoebaSentenceRepository creates a new instance of TatoebaSentenceRepository. It also registers a cleanup function to assert the mocks expectations.
func NewTatoebaSentenceRepository(t testing.TB) *TatoebaSentenceRepository {
	mock := &TatoebaSentenceRepository{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
