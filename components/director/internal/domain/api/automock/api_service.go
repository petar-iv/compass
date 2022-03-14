// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// APIService is an autogenerated mock type for the APIService type
type APIService struct {
	mock.Mock
}

// CreateInBundle provides a mock function with given fields: ctx, appID, bundleID, in, spec
func (_m *APIService) CreateInBundle(ctx context.Context, appID string, bundleID string, in model.APIDefinitionInput, spec *model.SpecInput) (string, error) {
	ret := _m.Called(ctx, appID, bundleID, in, spec)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, model.APIDefinitionInput, *model.SpecInput) string); ok {
		r0 = rf(ctx, appID, bundleID, in, spec)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, model.APIDefinitionInput, *model.SpecInput) error); ok {
		r1 = rf(ctx, appID, bundleID, in, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *APIService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *APIService) Get(ctx context.Context, id string) (*model.APIDefinition, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.APIDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.APIDefinition); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFetchRequests provides a mock function with given fields: ctx, specIDs
func (_m *APIService) ListFetchRequests(ctx context.Context, specIDs []string) ([]*model.FetchRequest, error) {
	ret := _m.Called(ctx, specIDs)

	var r0 []*model.FetchRequest
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*model.FetchRequest); ok {
		r0 = rf(ctx, specIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, specIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, in, spec
func (_m *APIService) Update(ctx context.Context, id string, in model.APIDefinitionInput, spec *model.SpecInput) error {
	ret := _m.Called(ctx, id, in, spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.APIDefinitionInput, *model.SpecInput) error); ok {
		r0 = rf(ctx, id, in, spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
