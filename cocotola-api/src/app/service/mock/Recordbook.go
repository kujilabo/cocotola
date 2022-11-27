// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola/cocotola-api/src/app/service"

	testing "testing"
)

// Recordbook is an autogenerated mock type for the Recordbook type
type Recordbook struct {
	mock.Mock
}

// GetResults provides a mock function with given fields: ctx
func (_m *Recordbook) GetResults(ctx context.Context) (map[domain.ProblemID]domain.StudyRecord, error) {
	ret := _m.Called(ctx)

	var r0 map[domain.ProblemID]domain.StudyRecord
	if rf, ok := ret.Get(0).(func(context.Context) map[domain.ProblemID]domain.StudyRecord); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[domain.ProblemID]domain.StudyRecord)
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

// GetResultsSortedLevel provides a mock function with given fields: ctx
func (_m *Recordbook) GetResultsSortedLevel(ctx context.Context) ([]domain.StudyRecordWithProblemID, error) {
	ret := _m.Called(ctx)

	var r0 []domain.StudyRecordWithProblemID
	if rf, ok := ret.Get(0).(func(context.Context) []domain.StudyRecordWithProblemID); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.StudyRecordWithProblemID)
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

// GetStudent provides a mock function with given fields:
func (_m *Recordbook) GetStudent() service.Student {
	ret := _m.Called()

	var r0 service.Student
	if rf, ok := ret.Get(0).(func() service.Student); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Student)
		}
	}

	return r0
}

// GetWorkbookID provides a mock function with given fields:
func (_m *Recordbook) GetWorkbookID() domain.WorkbookID {
	ret := _m.Called()

	var r0 domain.WorkbookID
	if rf, ok := ret.Get(0).(func() domain.WorkbookID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.WorkbookID)
	}

	return r0
}

// SetResult provides a mock function with given fields: ctx, problemType, problemID, result, mastered
func (_m *Recordbook) SetResult(ctx context.Context, problemType string, problemID domain.ProblemID, result bool, mastered bool) error {
	ret := _m.Called(ctx, problemType, problemID, result, mastered)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.ProblemID, bool, bool) error); ok {
		r0 = rf(ctx, problemType, problemID, result, mastered)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRecordbook creates a new instance of Recordbook. It also registers a cleanup function to assert the mocks expectations.
func NewRecordbook(t testing.TB) *Recordbook {
	mock := &Recordbook{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
