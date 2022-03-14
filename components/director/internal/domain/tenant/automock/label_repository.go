// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// LabelRepository is an autogenerated mock type for the LabelRepository type
type LabelRepository struct {
	mock.Mock
}

// ListForObject provides a mock function with given fields: ctx, _a1, objectType, objectID
func (_m *LabelRepository) ListForObject(ctx context.Context, _a1 string, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error) {
	ret := _m.Called(ctx, _a1, objectType, objectID)

	var r0 map[string]*model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string) map[string]*model.Label); ok {
		r0 = rf(ctx, _a1, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.LabelableObject, string) error); ok {
		r1 = rf(ctx, _a1, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
