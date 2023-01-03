// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola/cocotola-translator-api/src/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// AdminUsecase is an autogenerated mock type for the AdminUsecase type
type AdminUsecase struct {
	mock.Mock
}

// AddTranslation provides a mock function with given fields: ctx, param
func (_m *AdminUsecase) AddTranslation(ctx context.Context, param domain.TranslationAddParameter) error {
	ret := _m.Called(ctx, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.TranslationAddParameter) error); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindTranslationByText provides a mock function with given fields: ctx, lang2, text
func (_m *AdminUsecase) FindTranslationByText(ctx context.Context, lang2 domain.Lang2, text string) ([]domain.Translation, error) {
	ret := _m.Called(ctx, lang2, text)

	var r0 []domain.Translation
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang2, string) []domain.Translation); ok {
		r0 = rf(ctx, lang2, text)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Translation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Lang2, string) error); ok {
		r1 = rf(ctx, lang2, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTranslationByTextAndPos provides a mock function with given fields: ctx, lang2, text, pos
func (_m *AdminUsecase) FindTranslationByTextAndPos(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	ret := _m.Called(ctx, lang2, text, pos)

	var r0 domain.Translation
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang2, string, domain.WordPos) domain.Translation); ok {
		r0 = rf(ctx, lang2, text, pos)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Translation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Lang2, string, domain.WordPos) error); ok {
		r1 = rf(ctx, lang2, text, pos)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTranslationsByFirstLetter provides a mock function with given fields: ctx, lang2, firstLetter
func (_m *AdminUsecase) FindTranslationsByFirstLetter(ctx context.Context, lang2 domain.Lang2, firstLetter string) ([]domain.Translation, error) {
	ret := _m.Called(ctx, lang2, firstLetter)

	var r0 []domain.Translation
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang2, string) []domain.Translation); ok {
		r0 = rf(ctx, lang2, firstLetter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Translation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Lang2, string) error); ok {
		r1 = rf(ctx, lang2, firstLetter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveTranslation provides a mock function with given fields: ctx, lang2, text, pos
func (_m *AdminUsecase) RemoveTranslation(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos) error {
	ret := _m.Called(ctx, lang2, text, pos)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang2, string, domain.WordPos) error); ok {
		r0 = rf(ctx, lang2, text, pos)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTranslation provides a mock function with given fields: ctx, lang2, text, pos, param
func (_m *AdminUsecase) UpdateTranslation(ctx context.Context, lang2 domain.Lang2, text string, pos domain.WordPos, param domain.TranslationUpdateParameter) error {
	ret := _m.Called(ctx, lang2, text, pos, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang2, string, domain.WordPos, domain.TranslationUpdateParameter) error); ok {
		r0 = rf(ctx, lang2, text, pos, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAdminUsecase creates a new instance of AdminUsecase. It also registers a cleanup function to assert the mocks expectations.
func NewAdminUsecase(t testing.TB) *AdminUsecase {
	mock := &AdminUsecase{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
