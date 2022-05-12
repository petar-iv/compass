// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// GlobalRegistryService is an autogenerated mock type for the GlobalRegistryService type
type GlobalRegistryService struct {
	mock.Mock
}

// ListGlobalResources provides a mock function with given fields: ctx
func (_m *GlobalRegistryService) ListGlobalResources(ctx context.Context) (map[string]bool, error) {
	ret := _m.Called(ctx)

	var r0 map[string]bool
	if rf, ok := ret.Get(0).(func(context.Context) map[string]bool); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]bool)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SyncGlobalResources provides a mock function with given fields: ctx
func (_m *GlobalRegistryService) SyncGlobalResources(ctx context.Context) (map[string]bool, error) {
	ret := _m.Called(ctx)

	var r0 map[string]bool
	if rf, ok := ret.Get(0).(func(context.Context) map[string]bool); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]bool)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGlobalRegistryService creates a new instance of GlobalRegistryService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewGlobalRegistryService(t testing.TB) *GlobalRegistryService {
	mock := &GlobalRegistryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
