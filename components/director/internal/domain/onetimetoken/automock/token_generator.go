// Code generated by mockery. DO NOT EDIT.

package automock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TokenGenerator is an autogenerated mock type for the TokenGenerator type
type TokenGenerator struct {
	mock.Mock
}

// NewToken provides a mock function with given fields:
func (_m *TokenGenerator) NewToken() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTokenGenerator creates a new instance of TokenGenerator. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenGenerator(t testing.TB) *TokenGenerator {
	mock := &TokenGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
