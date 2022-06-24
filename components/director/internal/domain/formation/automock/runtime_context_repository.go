// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// RuntimeContextRepository is an autogenerated mock type for the runtimeContextRepository type
type RuntimeContextRepository struct {
	mock.Mock
}

// Exists provides a mock function with given fields: ctx, tenant, id
func (_m *RuntimeContextRepository) Exists(ctx context.Context, tenant string, id string) (bool, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields: ctx, tenant
func (_m *RuntimeContextRepository) ListAll(ctx context.Context, tenant string) ([]*model.RuntimeContext, error) {
	ret := _m.Called(ctx, tenant)

	var r0 []*model.RuntimeContext
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.RuntimeContext); ok {
		r0 = rf(ctx, tenant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.RuntimeContext)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRuntimeContextRepository creates a new instance of RuntimeContextRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuntimeContextRepository(t testing.TB) *RuntimeContextRepository {
	mock := &RuntimeContextRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
