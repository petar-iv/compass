// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	runtimectx "github.com/kyma-incubator/compass/components/director/internal/domain/runtime_context"
	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *runtimectx.RuntimeContext) *model.RuntimeContext {
	ret := _m.Called(entity)

	var r0 *model.RuntimeContext
	if rf, ok := ret.Get(0).(func(*runtimectx.RuntimeContext) *model.RuntimeContext); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimeContext)
		}
	}

	return r0
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.RuntimeContext) *runtimectx.RuntimeContext {
	ret := _m.Called(in)

	var r0 *runtimectx.RuntimeContext
	if rf, ok := ret.Get(0).(func(*model.RuntimeContext) *runtimectx.RuntimeContext); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtimectx.RuntimeContext)
		}
	}

	return r0
}