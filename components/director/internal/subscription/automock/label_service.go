// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// LabelService is an autogenerated mock type for the LabelService type
type LabelService struct {
	mock.Mock
}

// CreateLabel provides a mock function with given fields: ctx, tenant, id, labelInput
func (_m *LabelService) CreateLabel(ctx context.Context, tenant string, id string, labelInput *model.LabelInput) error {
	ret := _m.Called(ctx, tenant, id, labelInput)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.LabelInput) error); ok {
		r0 = rf(ctx, tenant, id, labelInput)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLabel provides a mock function with given fields: ctx, tenant, labelInput
func (_m *LabelService) GetLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) (*model.Label, error) {
	ret := _m.Called(ctx, tenant, labelInput)

	var r0 *model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.LabelInput) *model.Label); ok {
		r0 = rf(ctx, tenant, labelInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *model.LabelInput) error); ok {
		r1 = rf(ctx, tenant, labelInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLabel provides a mock function with given fields: ctx, tenant, id, labelInput
func (_m *LabelService) UpdateLabel(ctx context.Context, tenant string, id string, labelInput *model.LabelInput) error {
	ret := _m.Called(ctx, tenant, id, labelInput)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.LabelInput) error); ok {
		r0 = rf(ctx, tenant, id, labelInput)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLabelService creates a new instance of LabelService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewLabelService(t testing.TB) *LabelService {
	mock := &LabelService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
