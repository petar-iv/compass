// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	director "github.com/kyma-incubator/compass/components/hydrator/internal/director"

	mock "github.com/stretchr/testify/mock"

	systemauth "github.com/kyma-incubator/compass/components/director/pkg/systemauth"
)

// DirectorClient is an autogenerated mock type for the DirectorClient type
type DirectorClient struct {
	mock.Mock
}

// GetSystemAuthByID provides a mock function with given fields: ctx, authID
func (_m *DirectorClient) GetSystemAuthByID(ctx context.Context, authID string) (*systemauth.SystemAuth, error) {
	ret := _m.Called(ctx, authID)

	var r0 *systemauth.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, string) *systemauth.SystemAuth); ok {
		r0 = rf(ctx, authID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*systemauth.SystemAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, authID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantByExternalID provides a mock function with given fields: ctx, tenantID
func (_m *DirectorClient) GetTenantByExternalID(ctx context.Context, tenantID string) (*graphql.Tenant, error) {
	ret := _m.Called(ctx, tenantID)

	var r0 *graphql.Tenant
	if rf, ok := ret.Get(0).(func(context.Context, string) *graphql.Tenant); ok {
		r0 = rf(ctx, tenantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.Tenant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSystemAuth provides a mock function with given fields: ctx, sysAuth
func (_m *DirectorClient) UpdateSystemAuth(ctx context.Context, sysAuth *systemauth.SystemAuth) (director.UpdateAuthResult, error) {
	ret := _m.Called(ctx, sysAuth)

	var r0 director.UpdateAuthResult
	if rf, ok := ret.Get(0).(func(context.Context, *systemauth.SystemAuth) director.UpdateAuthResult); ok {
		r0 = rf(ctx, sysAuth)
	} else {
		r0 = ret.Get(0).(director.UpdateAuthResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *systemauth.SystemAuth) error); ok {
		r1 = rf(ctx, sysAuth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
