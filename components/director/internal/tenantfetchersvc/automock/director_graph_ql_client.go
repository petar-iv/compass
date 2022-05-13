// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// DirectorGraphQLClient is an autogenerated mock type for the DirectorGraphQLClient type
type DirectorGraphQLClient struct {
	mock.Mock
}

// SubscribeTenantToRuntime provides a mock function with given fields: ctx, providerID, subaccountID, region
func (_m *DirectorGraphQLClient) SubscribeTenantToRuntime(ctx context.Context, providerID string, subaccountID string, region string) error {
	ret := _m.Called(ctx, providerID, subaccountID, region)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, providerID, subaccountID, region)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnsubscribeTenantFromRuntime provides a mock function with given fields: ctx, providerID, subaccountID, region
func (_m *DirectorGraphQLClient) UnsubscribeTenantFromRuntime(ctx context.Context, providerID string, subaccountID string, region string) error {
	ret := _m.Called(ctx, providerID, subaccountID, region)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, providerID, subaccountID, region)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteTenants provides a mock function with given fields: _a0, _a1
func (_m *DirectorGraphQLClient) WriteTenants(_a0 context.Context, _a1 []graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDirectorGraphQLClient creates a new instance of DirectorGraphQLClient. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewDirectorGraphQLClient(t testing.TB) *DirectorGraphQLClient {
	mock := &DirectorGraphQLClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
