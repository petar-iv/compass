// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	context "context"

	labelfilter "github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// RuntimeCtxService is an autogenerated mock type for the RuntimeCtxService type
type RuntimeCtxService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *RuntimeCtxService) Create(ctx context.Context, in model.RuntimeContextInput) (string, error) {
	ret := _m.Called(ctx, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.RuntimeContextInput) string); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.RuntimeContextInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *RuntimeCtxService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListByFilter provides a mock function with given fields: ctx, runtimeID, filter, pageSize, cursor
func (_m *RuntimeCtxService) ListByFilter(ctx context.Context, runtimeID string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimeContextPage, error) {
	ret := _m.Called(ctx, runtimeID, filter, pageSize, cursor)

	var r0 *model.RuntimeContextPage
	if rf, ok := ret.Get(0).(func(context.Context, string, []*labelfilter.LabelFilter, int, string) *model.RuntimeContextPage); ok {
		r0 = rf(ctx, runtimeID, filter, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimeContextPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []*labelfilter.LabelFilter, int, string) error); ok {
		r1 = rf(ctx, runtimeID, filter, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRuntimeCtxService creates a new instance of RuntimeCtxService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuntimeCtxService(t testing.TB) *RuntimeCtxService {
	mock := &RuntimeCtxService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
