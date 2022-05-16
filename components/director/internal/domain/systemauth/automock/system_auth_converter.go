// Code generated by mockery. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/pkg/model"

	testing "testing"
)

// SystemAuthConverter is an autogenerated mock type for the SystemAuthConverter type
type SystemAuthConverter struct {
	mock.Mock
}

// ToGraphQL provides a mock function with given fields: _a0
func (_m *SystemAuthConverter) ToGraphQL(_a0 *model.SystemAuth) (graphql.SystemAuth, error) {
	ret := _m.Called(_a0)

	var r0 graphql.SystemAuth
	if rf, ok := ret.Get(0).(func(*model.SystemAuth) graphql.SystemAuth); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(graphql.SystemAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.SystemAuth) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSystemAuthConverter creates a new instance of SystemAuthConverter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSystemAuthConverter(t testing.TB) *SystemAuthConverter {
	mock := &SystemAuthConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
