// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// AsaService is an autogenerated mock type for the asaService type
type AsaService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *AsaService) Create(ctx context.Context, in model.AutomaticScenarioAssignment) (model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, in)

	var r0 model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(model.AutomaticScenarioAssignment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in
func (_m *AsaService) Delete(ctx context.Context, in model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteManyForSameTargetTenant provides a mock function with given fields: ctx, in
func (_m *AsaService) DeleteManyForSameTargetTenant(ctx context.Context, in []*model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetForScenarioName provides a mock function with given fields: ctx, scenarioName
func (_m *AsaService) GetForScenarioName(ctx context.Context, scenarioName string) (model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, scenarioName)

	var r0 model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, string) model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, scenarioName)
	} else {
		r0 = ret.Get(0).(model.AutomaticScenarioAssignment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, scenarioName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, pageSize, cursor
func (_m *AsaService) List(ctx context.Context, pageSize int, cursor string) (*model.AutomaticScenarioAssignmentPage, error) {
	ret := _m.Called(ctx, pageSize, cursor)

	var r0 *model.AutomaticScenarioAssignmentPage
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *model.AutomaticScenarioAssignmentPage); ok {
		r0 = rf(ctx, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AutomaticScenarioAssignmentPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForTargetTenant provides a mock function with given fields: ctx, targetTenantInternalID
func (_m *AsaService) ListForTargetTenant(ctx context.Context, targetTenantInternalID string) ([]*model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, targetTenantInternalID)

	var r0 []*model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, targetTenantInternalID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AutomaticScenarioAssignment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, targetTenantInternalID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
