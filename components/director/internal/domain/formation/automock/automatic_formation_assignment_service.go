// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// AutomaticFormationAssignmentService is an autogenerated mock type for the automaticFormationAssignmentService type
type AutomaticFormationAssignmentService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *AutomaticFormationAssignmentService) Create(ctx context.Context, in model.AutomaticScenarioAssignment) (model.AutomaticScenarioAssignment, error) {
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
func (_m *AutomaticFormationAssignmentService) Delete(ctx context.Context, in model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetForScenarioName provides a mock function with given fields: ctx, scenarioName
func (_m *AutomaticFormationAssignmentService) GetForScenarioName(ctx context.Context, scenarioName string) (model.AutomaticScenarioAssignment, error) {
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

// NewAutomaticFormationAssignmentService creates a new instance of AutomaticFormationAssignmentService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAutomaticFormationAssignmentService(t testing.TB) *AutomaticFormationAssignmentService {
	mock := &AutomaticFormationAssignmentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
