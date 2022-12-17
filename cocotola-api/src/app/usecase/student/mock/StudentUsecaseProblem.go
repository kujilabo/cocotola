// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	appdomain "github.com/kujilabo/cocotola/cocotola-api/src/app/domain"

	domain "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"

	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola/cocotola-api/src/app/service"

	testing "testing"
)

// StudentUsecaseProblem is an autogenerated mock type for the StudentUsecaseProblem type
type StudentUsecaseProblem struct {
	mock.Mock
}

// AddProblem provides a mock function with given fields: ctx, organizationID, operatorID, param
func (_m *StudentUsecaseProblem) AddProblem(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, param service.ProblemAddParameter) ([]appdomain.ProblemID, error) {
	ret := _m.Called(ctx, organizationID, operatorID, param)

	var r0 []appdomain.ProblemID
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemAddParameter) []appdomain.ProblemID); ok {
		r0 = rf(ctx, organizationID, operatorID, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]appdomain.ProblemID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemAddParameter) error); ok {
		r1 = rf(ctx, organizationID, operatorID, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllProblemsByWorkbookID provides a mock function with given fields: ctx, organizationID, operatorID, workbookID
func (_m *StudentUsecaseProblem) FindAllProblemsByWorkbookID(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, workbookID appdomain.WorkbookID) (service.ProblemSearchResult, error) {
	ret := _m.Called(ctx, organizationID, operatorID, workbookID)

	var r0 service.ProblemSearchResult
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID) service.ProblemSearchResult); ok {
		r0 = rf(ctx, organizationID, operatorID, workbookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.ProblemSearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID) error); ok {
		r1 = rf(ctx, organizationID, operatorID, workbookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindProblemByID provides a mock function with given fields: ctx, organizationID, operatorID, id
func (_m *StudentUsecaseProblem) FindProblemByID(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, id service.ProblemSelectParameter1) (appdomain.ProblemModel, error) {
	ret := _m.Called(ctx, organizationID, operatorID, id)

	var r0 appdomain.ProblemModel
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemSelectParameter1) appdomain.ProblemModel); ok {
		r0 = rf(ctx, organizationID, operatorID, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(appdomain.ProblemModel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemSelectParameter1) error); ok {
		r1 = rf(ctx, organizationID, operatorID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindProblemIDs provides a mock function with given fields: ctx, organizationID, operatorID, workbookID
func (_m *StudentUsecaseProblem) FindProblemIDs(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, workbookID appdomain.WorkbookID) ([]appdomain.ProblemID, error) {
	ret := _m.Called(ctx, organizationID, operatorID, workbookID)

	var r0 []appdomain.ProblemID
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID) []appdomain.ProblemID); ok {
		r0 = rf(ctx, organizationID, operatorID, workbookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]appdomain.ProblemID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID) error); ok {
		r1 = rf(ctx, organizationID, operatorID, workbookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindProblemsByProblemIDs provides a mock function with given fields: ctx, organizationID, operatorID, workbookID, param
func (_m *StudentUsecaseProblem) FindProblemsByProblemIDs(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, workbookID appdomain.WorkbookID, param service.ProblemIDsCondition) (service.ProblemSearchResult, error) {
	ret := _m.Called(ctx, organizationID, operatorID, workbookID, param)

	var r0 service.ProblemSearchResult
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID, service.ProblemIDsCondition) service.ProblemSearchResult); ok {
		r0 = rf(ctx, organizationID, operatorID, workbookID, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.ProblemSearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID, service.ProblemIDsCondition) error); ok {
		r1 = rf(ctx, organizationID, operatorID, workbookID, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindProblemsByWorkbookID provides a mock function with given fields: ctx, organizationID, operatorID, workbookID, param
func (_m *StudentUsecaseProblem) FindProblemsByWorkbookID(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, workbookID appdomain.WorkbookID, param service.ProblemSearchCondition) (service.ProblemSearchResult, error) {
	ret := _m.Called(ctx, organizationID, operatorID, workbookID, param)

	var r0 service.ProblemSearchResult
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID, service.ProblemSearchCondition) service.ProblemSearchResult); ok {
		r0 = rf(ctx, organizationID, operatorID, workbookID, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.ProblemSearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID, service.ProblemSearchCondition) error); ok {
		r1 = rf(ctx, organizationID, operatorID, workbookID, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ImportProblems provides a mock function with given fields: ctx, organizationID, operatorID, workbookID, newIterator
func (_m *StudentUsecaseProblem) ImportProblems(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, workbookID appdomain.WorkbookID, newIterator func(appdomain.WorkbookID, string) (service.ProblemAddParameterIterator, error)) error {
	ret := _m.Called(ctx, organizationID, operatorID, workbookID, newIterator)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, appdomain.WorkbookID, func(appdomain.WorkbookID, string) (service.ProblemAddParameterIterator, error)) error); ok {
		r0 = rf(ctx, organizationID, operatorID, workbookID, newIterator)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveProblem provides a mock function with given fields: ctx, organizationID, operatorID, id
func (_m *StudentUsecaseProblem) RemoveProblem(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, id service.ProblemSelectParameter2) error {
	ret := _m.Called(ctx, organizationID, operatorID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemSelectParameter2) error); ok {
		r0 = rf(ctx, organizationID, operatorID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateProblem provides a mock function with given fields: ctx, organizationID, operatorID, id, param
func (_m *StudentUsecaseProblem) UpdateProblem(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error {
	ret := _m.Called(ctx, organizationID, operatorID, id, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemSelectParameter2, service.ProblemUpdateParameter) error); ok {
		r0 = rf(ctx, organizationID, operatorID, id, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateProblemProperty provides a mock function with given fields: ctx, organizationID, operatorID, id, param
func (_m *StudentUsecaseProblem) UpdateProblemProperty(ctx context.Context, organizationID domain.OrganizationID, operatorID domain.AppUserID, id service.ProblemSelectParameter2, param service.ProblemUpdateParameter) error {
	ret := _m.Called(ctx, organizationID, operatorID, id, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.OrganizationID, domain.AppUserID, service.ProblemSelectParameter2, service.ProblemUpdateParameter) error); ok {
		r0 = rf(ctx, organizationID, operatorID, id, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStudentUsecaseProblem creates a new instance of StudentUsecaseProblem. It also registers a cleanup function to assert the mocks expectations.
func NewStudentUsecaseProblem(t testing.TB) *StudentUsecaseProblem {
	mock := &StudentUsecaseProblem{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}