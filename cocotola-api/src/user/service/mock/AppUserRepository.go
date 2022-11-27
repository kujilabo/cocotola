// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola/cocotola-api/src/user/service"

	testing "testing"
)

// AppUserRepository is an autogenerated mock type for the AppUserRepository type
type AppUserRepository struct {
	mock.Mock
}

// AddAppUser provides a mock function with given fields: ctx, operator, param
func (_m *AppUserRepository) AddAppUser(ctx context.Context, operator domain.OwnerModel, param service.AppUserAddParameter) (domain.AppUserID, error) {
	ret := _m.Called(ctx, operator, param)

	var r0 domain.AppUserID
	if rf, ok := ret.Get(0).(func(context.Context, domain.OwnerModel, service.AppUserAddParameter) domain.AppUserID); ok {
		r0 = rf(ctx, operator, param)
	} else {
		r0 = ret.Get(0).(domain.AppUserID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.OwnerModel, service.AppUserAddParameter) error); ok {
		r1 = rf(ctx, operator, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddFirstOwner provides a mock function with given fields: ctx, operator, param
func (_m *AppUserRepository) AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param service.FirstOwnerAddParameter) (domain.AppUserID, error) {
	ret := _m.Called(ctx, operator, param)

	var r0 domain.AppUserID
	if rf, ok := ret.Get(0).(func(context.Context, domain.SystemOwnerModel, service.FirstOwnerAddParameter) domain.AppUserID); ok {
		r0 = rf(ctx, operator, param)
	} else {
		r0 = ret.Get(0).(domain.AppUserID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.SystemOwnerModel, service.FirstOwnerAddParameter) error); ok {
		r1 = rf(ctx, operator, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddSystemOwner provides a mock function with given fields: ctx, operator, organizationID
func (_m *AppUserRepository) AddSystemOwner(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.AppUserID, error) {
	ret := _m.Called(ctx, operator, organizationID)

	var r0 domain.AppUserID
	if rf, ok := ret.Get(0).(func(context.Context, domain.SystemAdminModel, domain.OrganizationID) domain.AppUserID); ok {
		r0 = rf(ctx, operator, organizationID)
	} else {
		r0 = ret.Get(0).(domain.AppUserID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.SystemAdminModel, domain.OrganizationID) error); ok {
		r1 = rf(ctx, operator, organizationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAppUserByID provides a mock function with given fields: ctx, operator, id
func (_m *AppUserRepository) FindAppUserByID(ctx context.Context, operator domain.AppUserModel, id domain.AppUserID) (service.AppUser, error) {
	ret := _m.Called(ctx, operator, id)

	var r0 service.AppUser
	if rf, ok := ret.Get(0).(func(context.Context, domain.AppUserModel, domain.AppUserID) service.AppUser); ok {
		r0 = rf(ctx, operator, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.AppUser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AppUserModel, domain.AppUserID) error); ok {
		r1 = rf(ctx, operator, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAppUserByLoginID provides a mock function with given fields: ctx, operator, loginID
func (_m *AppUserRepository) FindAppUserByLoginID(ctx context.Context, operator domain.AppUserModel, loginID string) (service.AppUser, error) {
	ret := _m.Called(ctx, operator, loginID)

	var r0 service.AppUser
	if rf, ok := ret.Get(0).(func(context.Context, domain.AppUserModel, string) service.AppUser); ok {
		r0 = rf(ctx, operator, loginID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.AppUser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AppUserModel, string) error); ok {
		r1 = rf(ctx, operator, loginID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAppUserIDs provides a mock function with given fields: ctx, operator, pageNo, pageSize
func (_m *AppUserRepository) FindAppUserIDs(ctx context.Context, operator domain.SystemOwnerModel, pageNo int, pageSize int) ([]domain.AppUserID, error) {
	ret := _m.Called(ctx, operator, pageNo, pageSize)

	var r0 []domain.AppUserID
	if rf, ok := ret.Get(0).(func(context.Context, domain.SystemOwnerModel, int, int) []domain.AppUserID); ok {
		r0 = rf(ctx, operator, pageNo, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.AppUserID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.SystemOwnerModel, int, int) error); ok {
		r1 = rf(ctx, operator, pageNo, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOwnerByLoginID provides a mock function with given fields: ctx, operator, loginID
func (_m *AppUserRepository) FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwnerModel, loginID string) (service.Owner, error) {
	ret := _m.Called(ctx, operator, loginID)

	var r0 service.Owner
	if rf, ok := ret.Get(0).(func(context.Context, domain.SystemOwnerModel, string) service.Owner); ok {
		r0 = rf(ctx, operator, loginID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Owner)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.SystemOwnerModel, string) error); ok {
		r1 = rf(ctx, operator, loginID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSystemOwnerByOrganizationID provides a mock function with given fields: ctx, operator, organizationID
func (_m *AppUserRepository) FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.SystemOwner, error) {
	ret := _m.Called(ctx, operator, organizationID)

	var r0 service.SystemOwner
	if rf, ok := ret.Get(0).(func(context.Context, domain.SystemAdminModel, domain.OrganizationID) service.SystemOwner); ok {
		r0 = rf(ctx, operator, organizationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.SystemOwner)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.SystemAdminModel, domain.OrganizationID) error); ok {
		r1 = rf(ctx, operator, organizationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSystemOwnerByOrganizationName provides a mock function with given fields: ctx, operator, organizationName
func (_m *AppUserRepository) FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminModel, organizationName string) (service.SystemOwner, error) {
	ret := _m.Called(ctx, operator, organizationName)

	var r0 service.SystemOwner
	if rf, ok := ret.Get(0).(func(context.Context, domain.SystemAdminModel, string) service.SystemOwner); ok {
		r0 = rf(ctx, operator, organizationName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.SystemOwner)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.SystemAdminModel, string) error); ok {
		r1 = rf(ctx, operator, organizationName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAppUserRepository creates a new instance of AppUserRepository. It also registers a cleanup function to assert the mocks expectations.
func NewAppUserRepository(t testing.TB) *AppUserRepository {
	mock := &AppUserRepository{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
