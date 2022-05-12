// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// LabelDefRepository is an autogenerated mock type for the labelDefRepository type
type LabelDefRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, def
func (_m *LabelDefRepository) Create(ctx context.Context, def model.LabelDefinition) error {
	ret := _m.Called(ctx, def)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelDefinition) error); ok {
		r0 = rf(ctx, def)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, tenant, key
func (_m *LabelDefRepository) Exists(ctx context.Context, tenant string, key string) (bool, error) {
	ret := _m.Called(ctx, tenant, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, tenant, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByKey provides a mock function with given fields: ctx, tenant, key
func (_m *LabelDefRepository) GetByKey(ctx context.Context, tenant string, key string) (*model.LabelDefinition, error) {
	ret := _m.Called(ctx, tenant, key)

	var r0 *model.LabelDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.LabelDefinition); ok {
		r0 = rf(ctx, tenant, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.LabelDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateWithVersion provides a mock function with given fields: ctx, def
func (_m *LabelDefRepository) UpdateWithVersion(ctx context.Context, def model.LabelDefinition) error {
	ret := _m.Called(ctx, def)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelDefinition) error); ok {
		r0 = rf(ctx, def)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLabelDefRepository creates a new instance of LabelDefRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewLabelDefRepository(t testing.TB) *LabelDefRepository {
	mock := &LabelDefRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
