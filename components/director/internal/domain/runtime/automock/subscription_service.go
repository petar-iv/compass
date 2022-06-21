// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// SubscriptionService is an autogenerated mock type for the SubscriptionService type
type SubscriptionService struct {
	mock.Mock
}

// SubscribeTenantToRuntime provides a mock function with given fields: ctx, providerID, subaccountTenantID, providerSubaccountID, region
func (_m *SubscriptionService) SubscribeTenantToRuntime(ctx context.Context, providerID string, subaccountTenantID string, providerSubaccountID string, region string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, providerSubaccountID, region)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, region)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnsubscribeTenantFromRuntime provides a mock function with given fields: ctx, providerID, subaccountTenantID, providerSubaccountID, region
func (_m *SubscriptionService) UnsubscribeTenantFromRuntime(ctx context.Context, providerID string, subaccountTenantID string, providerSubaccountID string, region string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, providerSubaccountID, region)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, region)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, region)
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
