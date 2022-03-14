// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	systemauth "github.com/kyma-incubator/compass/components/director/pkg/systemauth"
	mock "github.com/stretchr/testify/mock"
)

// OneTimeTokenService is an autogenerated mock type for the OneTimeTokenService type
type OneTimeTokenService struct {
	mock.Mock
}

// IsTokenValid provides a mock function with given fields: systemAuth
func (_m *OneTimeTokenService) IsTokenValid(systemAuth *systemauth.SystemAuth) (bool, error) {
	ret := _m.Called(systemAuth)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*systemauth.SystemAuth) bool); ok {
		r0 = rf(systemAuth)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*systemauth.SystemAuth) error); ok {
		r1 = rf(systemAuth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
