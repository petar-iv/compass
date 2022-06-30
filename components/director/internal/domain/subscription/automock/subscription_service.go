// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	resource "github.com/kyma-incubator/compass/components/director/pkg/resource"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// SubscriptionService is an autogenerated mock type for the SubscriptionService type
type SubscriptionService struct {
	mock.Mock
}

// DetermineSubscriptionFlow provides a mock function with given fields: ctx, providerID, region
func (_m *SubscriptionService) DetermineSubscriptionFlow(ctx context.Context, providerID string, region string) (resource.Type, error) {
	ret := _m.Called(ctx, providerID, region)

	var r0 resource.Type
	if rf, ok := ret.Get(0).(func(context.Context, string, string) resource.Type); ok {
		r0 = rf(ctx, providerID, region)
	} else {
		r0 = ret.Get(0).(resource.Type)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, providerID, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeTenantToApplication provides a mock function with given fields: ctx, providerID, subaccountTenantID, region, subscribedAppName
func (_m *SubscriptionService) SubscribeTenantToApplication(ctx context.Context, providerID string, subaccountTenantID string, region string, subscribedAppName string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, region, subscribedAppName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, region, subscribedAppName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, region, subscribedAppName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeTenantToRuntime provides a mock function with given fields: ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName
func (_m *SubscriptionService) SubscribeTenantToRuntime(ctx context.Context, providerID string, subaccountTenantID string, providerSubaccountID string, consumerTenantID string, region string, subscriptionAppName string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnsubscribeTenantFromApplication provides a mock function with given fields: ctx, providerID, subaccountTenantID, region
func (_m *SubscriptionService) UnsubscribeTenantFromApplication(ctx context.Context, providerID string, subaccountTenantID string, region string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, region)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, region)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnsubscribeTenantFromRuntime provides a mock function with given fields: ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region
func (_m *SubscriptionService) UnsubscribeTenantFromRuntime(ctx context.Context, providerID string, subaccountTenantID string, providerSubaccountID string, consumerTenantID string, region string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSubscriptionService creates a new instance of SubscriptionService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSubscriptionService(t testing.TB) *SubscriptionService {
	mock := &SubscriptionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
