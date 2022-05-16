// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/pkg/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/kyma-incubator/compass/components/director/internal/repo"

	testing "testing"
)

// SystemAuthRepo is an autogenerated mock type for the SystemAuthRepo type
type SystemAuthRepo struct {
	mock.Mock
}

// ListGlobalWithConditions provides a mock function with given fields: ctx, conditions
func (_m *SystemAuthRepo) ListGlobalWithConditions(ctx context.Context, conditions repo.Conditions) ([]model.SystemAuth, error) {
	ret := _m.Called(ctx, conditions)

	var r0 []model.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, repo.Conditions) []model.SystemAuth); ok {
		r0 = rf(ctx, conditions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.SystemAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, repo.Conditions) error); ok {
		r1 = rf(ctx, conditions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSystemAuthRepo creates a new instance of SystemAuthRepo. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSystemAuthRepo(t testing.TB) *SystemAuthRepo {
	mock := &SystemAuthRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
