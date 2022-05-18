// Code generated by mockery. DO NOT EDIT.

package automock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TenantSyncService is an autogenerated mock type for the TenantSyncService type
type TenantSyncService struct {
	mock.Mock
}

// SyncTenants provides a mock function with given fields:
func (_m *TenantSyncService) SyncTenants() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTenantSyncService creates a new instance of TenantSyncService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantSyncService(t testing.TB) *TenantSyncService {
	mock := &TenantSyncService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
