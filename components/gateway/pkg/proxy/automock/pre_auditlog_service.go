// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	proxy "github.com/kyma-incubator/compass/components/gateway/pkg/proxy"
	mock "github.com/stretchr/testify/mock"
)

// PreAuditlogService is an autogenerated mock type for the PreAuditlogService type
type PreAuditlogService struct {
	mock.Mock
}

// Log provides a mock function with given fields: ctx, msg
func (_m *PreAuditlogService) Log(ctx context.Context, msg proxy.AuditlogMessage) error {
	ret := _m.Called(ctx, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, proxy.AuditlogMessage) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PreLog provides a mock function with given fields: ctx, msg
func (_m *PreAuditlogService) PreLog(ctx context.Context, msg proxy.AuditlogMessage) error {
	ret := _m.Called(ctx, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, proxy.AuditlogMessage) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
