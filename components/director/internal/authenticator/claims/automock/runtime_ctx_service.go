// Code generated by mockery. DO NOT EDIT.

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