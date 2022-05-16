// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/gateway/pkg/auditlog/model"
	mock "github.com/stretchr/testify/mock"
)

// AuditlogClient is an autogenerated mock type for the AuditlogClient type
type AuditlogClient struct {
	mock.Mock
}

// LogConfigurationChange provides a mock function with given fields: ctx, change
func (_m *AuditlogClient) LogConfigurationChange(ctx context.Context, change model.ConfigurationChange) error {
	ret := _m.Called(ctx, change)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ConfigurationChange) error); ok {
		r0 = rf(ctx, change)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LogSecurityEvent provides a mock function with given fields: ctx, event
func (_m *AuditlogClient) LogSecurityEvent(ctx context.Context, event model.SecurityEvent) error {
	ret := _m.Called(ctx, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.SecurityEvent) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
