// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// FetchRequestConverter is an autogenerated mock type for the FetchRequestConverter type
type FetchRequestConverter struct {
	mock.Mock
}

// ToGraphQL provides a mock function with given fields: in
func (_m *FetchRequestConverter) ToGraphQL(in *model.FetchRequest) (*graphql.FetchRequest, error) {
	ret := _m.Called(in)

	var r0 *graphql.FetchRequest
	if rf, ok := ret.Get(0).(func(*model.FetchRequest) *graphql.FetchRequest); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.FetchRequest) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFetchRequestConverter creates a new instance of FetchRequestConverter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFetchRequestConverter(t testing.TB) *FetchRequestConverter {
	mock := &FetchRequestConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
