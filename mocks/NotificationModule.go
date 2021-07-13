// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// NotificationModule is an autogenerated mock type for the NotificationModule type
type NotificationModule struct {
	mock.Mock
}

// AddPushToken provides a mock function with given fields: ctx, userID, pushToken
func (_m *NotificationModule) AddPushToken(ctx context.Context, userID string, pushToken string) error {
	ret := _m.Called(ctx, userID, pushToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, userID, pushToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
