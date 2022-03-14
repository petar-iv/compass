// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// LabelUpsertService is an autogenerated mock type for the LabelUpsertService type
type LabelUpsertService struct {
	mock.Mock
}

// UpsertLabel provides a mock function with given fields: ctx, _a1, labelInput
func (_m *LabelUpsertService) UpsertLabel(ctx context.Context, _a1 string, labelInput *model.LabelInput) error {
	ret := _m.Called(ctx, _a1, labelInput)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.LabelInput) error); ok {
		r0 = rf(ctx, _a1, labelInput)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
