// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// GqlConverter is an autogenerated mock type for the gqlConverter type
type GqlConverter struct {
	mock.Mock
}

// FromInputGraphQL provides a mock function with given fields: in, targetTenantInternalID
func (_m *GqlConverter) FromInputGraphQL(in graphql.AutomaticScenarioAssignmentSetInput, targetTenantInternalID string) model.AutomaticScenarioAssignment {
	ret := _m.Called(in, targetTenantInternalID)

	var r0 model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(graphql.AutomaticScenarioAssignmentSetInput, string) model.AutomaticScenarioAssignment); ok {
		r0 = rf(in, targetTenantInternalID)
	} else {
		r0 = ret.Get(0).(model.AutomaticScenarioAssignment)
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: in, targetTenantExternalID
func (_m *GqlConverter) ToGraphQL(in model.AutomaticScenarioAssignment, targetTenantExternalID string) graphql.AutomaticScenarioAssignment {
	ret := _m.Called(in, targetTenantExternalID)

	var r0 graphql.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(model.AutomaticScenarioAssignment, string) graphql.AutomaticScenarioAssignment); ok {
		r0 = rf(in, targetTenantExternalID)
	} else {
		r0 = ret.Get(0).(graphql.AutomaticScenarioAssignment)
	}

	return r0
}