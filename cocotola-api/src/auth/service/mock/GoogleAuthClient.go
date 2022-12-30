// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	service "github.com/kujilabo/cocotola/cocotola-api/src/auth/service"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// GoogleAuthClient is an autogenerated mock type for the GoogleAuthClient type
type GoogleAuthClient struct {
	mock.Mock
}

// RetrieveAccessToken provides a mock function with given fields: ctx, code
func (_m *GoogleAuthClient) RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error) {
	ret := _m.Called(ctx, code)

	var r0 *service.GoogleAuthResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) *service.GoogleAuthResponse); ok {
		r0 = rf(ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.GoogleAuthResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveUserInfo provides a mock function with given fields: ctx, googleAuthResponse
func (_m *GoogleAuthClient) RetrieveUserInfo(ctx context.Context, googleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error) {
	ret := _m.Called(ctx, googleAuthResponse)

	var r0 *service.GoogleUserInfo
	if rf, ok := ret.Get(0).(func(context.Context, *service.GoogleAuthResponse) *service.GoogleUserInfo); ok {
		r0 = rf(ctx, googleAuthResponse)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.GoogleUserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *service.GoogleAuthResponse) error); ok {
		r1 = rf(ctx, googleAuthResponse)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGoogleAuthClient creates a new instance of GoogleAuthClient. It also registers a cleanup function to assert the mocks expectations.
func NewGoogleAuthClient(t testing.TB) *GoogleAuthClient {
	mock := &GoogleAuthClient{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
