// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"
)

// DirectorClient is an autogenerated mock type for the DirectorClient type
type DirectorClient struct {
	mock.Mock
}

// CreateAPIDefinition provides a mock function with given fields: ctx, bundleID, apiDefinitionInput
func (_m *DirectorClient) CreateAPIDefinition(ctx context.Context, bundleID string, apiDefinitionInput graphql.APIDefinitionInput) (string, error) {
	ret := _m.Called(ctx, bundleID, apiDefinitionInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.APIDefinitionInput) string); ok {
		r0 = rf(ctx, bundleID, apiDefinitionInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.APIDefinitionInput) error); ok {
		r1 = rf(ctx, bundleID, apiDefinitionInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateBundle provides a mock function with given fields: ctx, appID, in
func (_m *DirectorClient) CreateBundle(ctx context.Context, appID string, in graphql.BundleCreateInput) (string, error) {
	ret := _m.Called(ctx, appID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.BundleCreateInput) string); ok {
		r0 = rf(ctx, appID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.BundleCreateInput) error); ok {
		r1 = rf(ctx, appID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDocument provides a mock function with given fields: ctx, bundleID, documentInput
func (_m *DirectorClient) CreateDocument(ctx context.Context, bundleID string, documentInput graphql.DocumentInput) (string, error) {
	ret := _m.Called(ctx, bundleID, documentInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.DocumentInput) string); ok {
		r0 = rf(ctx, bundleID, documentInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.DocumentInput) error); ok {
		r1 = rf(ctx, bundleID, documentInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateEventDefinition provides a mock function with given fields: ctx, bundleID, eventDefinitionInput
func (_m *DirectorClient) CreateEventDefinition(ctx context.Context, bundleID string, eventDefinitionInput graphql.EventDefinitionInput) (string, error) {
	ret := _m.Called(ctx, bundleID, eventDefinitionInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.EventDefinitionInput) string); ok {
		r0 = rf(ctx, bundleID, eventDefinitionInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graphql.EventDefinitionInput) error); ok {
		r1 = rf(ctx, bundleID, eventDefinitionInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAPIDefinition provides a mock function with given fields: ctx, apiID
func (_m *DirectorClient) DeleteAPIDefinition(ctx context.Context, apiID string) error {
	ret := _m.Called(ctx, apiID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, apiID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteBundle provides a mock function with given fields: ctx, bundleID
func (_m *DirectorClient) DeleteBundle(ctx context.Context, bundleID string) error {
	ret := _m.Called(ctx, bundleID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, bundleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDocument provides a mock function with given fields: ctx, documentID
func (_m *DirectorClient) DeleteDocument(ctx context.Context, documentID string) error {
	ret := _m.Called(ctx, documentID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, documentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEventDefinition provides a mock function with given fields: ctx, eventID
func (_m *DirectorClient) DeleteEventDefinition(ctx context.Context, eventID string) error {
	ret := _m.Called(ctx, eventID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBundle provides a mock function with given fields: ctx, appID, bundleID
func (_m *DirectorClient) GetBundle(ctx context.Context, appID string, bundleID string) (graphql.BundleExt, error) {
	ret := _m.Called(ctx, appID, bundleID)

	var r0 graphql.BundleExt
	if rf, ok := ret.Get(0).(func(context.Context, string, string) graphql.BundleExt); ok {
		r0 = rf(ctx, appID, bundleID)
	} else {
		r0 = ret.Get(0).(graphql.BundleExt)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, appID, bundleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListBundles provides a mock function with given fields: ctx, appID
func (_m *DirectorClient) ListBundles(ctx context.Context, appID string) ([]*graphql.BundleExt, error) {
	ret := _m.Called(ctx, appID)

	var r0 []*graphql.BundleExt
	if rf, ok := ret.Get(0).(func(context.Context, string) []*graphql.BundleExt); ok {
		r0 = rf(ctx, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.BundleExt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetApplicationLabel provides a mock function with given fields: ctx, appID, label
func (_m *DirectorClient) SetApplicationLabel(ctx context.Context, appID string, label graphql.LabelInput) error {
	ret := _m.Called(ctx, appID, label)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.LabelInput) error); ok {
		r0 = rf(ctx, appID, label)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateBundle provides a mock function with given fields: ctx, bundleID, in
func (_m *DirectorClient) UpdateBundle(ctx context.Context, bundleID string, in graphql.BundleUpdateInput) error {
	ret := _m.Called(ctx, bundleID, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.BundleUpdateInput) error); ok {
		r0 = rf(ctx, bundleID, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
