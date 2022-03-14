// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	resource "github.com/kyma-incubator/compass/components/director/pkg/resource"
)

// TenantService is an autogenerated mock type for the TenantService type
type TenantService struct {
	mock.Mock
}

// GetLowestOwnerForResource provides a mock function with given fields: ctx, resourceType, objectID
func (_m *TenantService) GetLowestOwnerForResource(ctx context.Context, resourceType resource.Type, objectID string) (string, error) {
	ret := _m.Called(ctx, resourceType, objectID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, resource.Type, string) string); ok {
		r0 = rf(ctx, resourceType, objectID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, resource.Type, string) error); ok {
		r1 = rf(ctx, resourceType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
