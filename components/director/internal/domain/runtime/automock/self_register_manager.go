// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// SelfRegisterManager is an autogenerated mock type for the SelfRegisterManager type
type SelfRegisterManager struct {
	mock.Mock
}

// CleanupSelfRegisteredRuntime provides a mock function with given fields: ctx, selfRegisterLabelValue, region
func (_m *SelfRegisterManager) CleanupSelfRegisteredRuntime(ctx context.Context, selfRegisterLabelValue string, region string) error {
	ret := _m.Called(ctx, selfRegisterLabelValue, region)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, selfRegisterLabelValue, region)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSelfRegDistinguishingLabelKey provides a mock function with given fields:
func (_m *SelfRegisterManager) GetSelfRegDistinguishingLabelKey() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// PrepareRuntimeForSelfRegistration provides a mock function with given fields: ctx, in, id
func (_m *SelfRegisterManager) PrepareRuntimeForSelfRegistration(ctx context.Context, in model.RuntimeRegisterInput, id string) (map[string]interface{}, error) {
	ret := _m.Called(ctx, in, id)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context, model.RuntimeRegisterInput, string) map[string]interface{}); ok {
		r0 = rf(ctx, in, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.RuntimeRegisterInput, string) error); ok {
		r1 = rf(ctx, in, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSelfRegisterManager creates a new instance of SelfRegisterManager. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSelfRegisterManager(t testing.TB) *SelfRegisterManager {
	mock := &SelfRegisterManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
