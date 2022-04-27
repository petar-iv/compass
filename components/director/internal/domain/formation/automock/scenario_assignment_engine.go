// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ScenarioAssignmentEngine is an autogenerated mock type for the scenarioAssignmentEngine type
type ScenarioAssignmentEngine struct {
	mock.Mock
}

// GetScenariosFromMatchingASAs provides a mock function with given fields: ctx, runtimeID
func (_m *ScenarioAssignmentEngine) GetScenariosFromMatchingASAs(ctx context.Context, runtimeID string) ([]string, error) {
	ret := _m.Called(ctx, runtimeID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}