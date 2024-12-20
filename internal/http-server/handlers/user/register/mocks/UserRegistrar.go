// Code generated by mockery v2.30.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserRegistrar is an autogenerated mock type for the UserRegistrar type
type UserRegistrar struct {
	mock.Mock
}

// IsUserExist provides a mock function with given fields: ctx, login
func (_m *UserRegistrar) IsUserExist(ctx context.Context, login string) (bool, error) {
	ret := _m.Called(ctx, login)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, login)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, login)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUser provides a mock function with given fields: ctx, login, password, date_created
func (_m *UserRegistrar) SaveUser(ctx context.Context, login string, password string, date_created string) (int64, error) {
	ret := _m.Called(ctx, login, password, date_created)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (int64, error)); ok {
		return rf(ctx, login, password, date_created)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) int64); ok {
		r0 = rf(ctx, login, password, date_created)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, login, password, date_created)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserRegistrar creates a new instance of UserRegistrar. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRegistrar(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRegistrar {
	mock := &UserRegistrar{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
