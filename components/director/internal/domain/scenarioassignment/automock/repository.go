// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *Repository) Create(ctx context.Context, _a1 model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteForScenarioName provides a mock function with given fields: ctx, tenantID, scenarioName
func (_m *Repository) DeleteForScenarioName(ctx context.Context, tenantID string, scenarioName string) error {
	ret := _m.Called(ctx, tenantID, scenarioName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenantID, scenarioName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteForTargetTenant provides a mock function with given fields: ctx, tenantID, targetTenantID
func (_m *Repository) DeleteForTargetTenant(ctx context.Context, tenantID string, targetTenantID string) error {
	ret := _m.Called(ctx, tenantID, targetTenantID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenantID, targetTenantID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetForScenarioName provides a mock function with given fields: ctx, tenantID, scenarioName
func (_m *Repository) GetForScenarioName(ctx context.Context, tenantID string, scenarioName string) (model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, tenantID, scenarioName)

	var r0 model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, string, string) model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, tenantID, scenarioName)
	} else {
		r0 = ret.Get(0).(model.AutomaticScenarioAssignment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantID, scenarioName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, tenant, pageSize, cursor
func (_m *Repository) List(ctx context.Context, tenant string, pageSize int, cursor string) (*model.AutomaticScenarioAssignmentPage, error) {
	ret := _m.Called(ctx, tenant, pageSize, cursor)

	var r0 *model.AutomaticScenarioAssignmentPage
	if rf, ok := ret.Get(0).(func(context.Context, string, int, string) *model.AutomaticScenarioAssignmentPage); ok {
		r0 = rf(ctx, tenant, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AutomaticScenarioAssignmentPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, string) error); ok {
		r1 = rf(ctx, tenant, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields: ctx, tenantID
func (_m *Repository) ListAll(ctx context.Context, tenantID string) ([]*model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, tenantID)

	var r0 []*model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, tenantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AutomaticScenarioAssignment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForTargetTenant provides a mock function with given fields: ctx, tenantID, targetTenantID
func (_m *Repository) ListForTargetTenant(ctx context.Context, tenantID string, targetTenantID string) ([]*model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, tenantID, targetTenantID)

	var r0 []*model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, tenantID, targetTenantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AutomaticScenarioAssignment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantID, targetTenantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
