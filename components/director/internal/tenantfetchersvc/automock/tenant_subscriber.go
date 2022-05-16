// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"
	testing "testing"

	tenantfetchersvc "github.com/kyma-incubator/compass/components/director/internal/tenantfetchersvc"
	mock "github.com/stretchr/testify/mock"
)

// TenantSubscriber is an autogenerated mock type for the TenantSubscriber type
type TenantSubscriber struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: ctx, tenantSubscriptionRequest
func (_m *TenantSubscriber) Subscribe(ctx context.Context, tenantSubscriptionRequest *tenantfetchersvc.TenantSubscriptionRequest) error {
	ret := _m.Called(ctx, tenantSubscriptionRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *tenantfetchersvc.TenantSubscriptionRequest) error); ok {
		r0 = rf(ctx, tenantSubscriptionRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsubscribe provides a mock function with given fields: ctx, tenantSubscriptionRequest
func (_m *TenantSubscriber) Unsubscribe(ctx context.Context, tenantSubscriptionRequest *tenantfetchersvc.TenantSubscriptionRequest) error {
	ret := _m.Called(ctx, tenantSubscriptionRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *tenantfetchersvc.TenantSubscriptionRequest) error); ok {
		r0 = rf(ctx, tenantSubscriptionRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTenantSubscriber creates a new instance of TenantSubscriber. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantSubscriber(t testing.TB) *TenantSubscriber {
	mock := &TenantSubscriber{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
