// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	systemauth "github.com/kyma-incubator/compass/components/director/internal/domain/systemauth"
	pkgsystemauth "github.com/kyma-incubator/compass/components/director/pkg/systemauth"
	mock "github.com/stretchr/testify/mock"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: in
func (_m *Converter) FromEntity(in systemauth.Entity) (pkgsystemauth.SystemAuth, error) {
	ret := _m.Called(in)

	var r0 pkgsystemauth.SystemAuth
	if rf, ok := ret.Get(0).(func(systemauth.Entity) pkgsystemauth.SystemAuth); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(pkgsystemauth.SystemAuth)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(systemauth.Entity) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *Converter) ToEntity(in pkgsystemauth.SystemAuth) (systemauth.Entity, error) {
	ret := _m.Called(in)

	var r0 systemauth.Entity
	if rf, ok := ret.Get(0).(func(pkgsystemauth.SystemAuth) systemauth.Entity); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(systemauth.Entity)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(pkgsystemauth.SystemAuth) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
