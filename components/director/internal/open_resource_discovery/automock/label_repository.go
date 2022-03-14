// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// LabelRepository is an autogenerated mock type for the labelRepository type
type LabelRepository struct {
	mock.Mock
}

// ListGlobalByKeyAndObjects provides a mock function with given fields: ctx, objectType, objectIDs, key
func (_m *LabelRepository) ListGlobalByKeyAndObjects(ctx context.Context, objectType model.LabelableObject, objectIDs []string, key string) ([]*model.Label, error) {
	ret := _m.Called(ctx, objectType, objectIDs, key)

	var r0 []*model.Label
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelableObject, []string, string) []*model.Label); ok {
		r0 = rf(ctx, objectType, objectIDs, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.LabelableObject, []string, string) error); ok {
		r1 = rf(ctx, objectType, objectIDs, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
