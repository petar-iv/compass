// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FetchRequestRepository is an autogenerated mock type for the FetchRequestRepository type
type FetchRequestRepository struct {
	mock.Mock
}

// Update provides a mock function with given fields: ctx, tenant, item
func (_m *FetchRequestRepository) Update(ctx context.Context, tenant string, item *model.FetchRequest) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.FetchRequest) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
