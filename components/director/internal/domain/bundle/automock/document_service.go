// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// DocumentService is an autogenerated mock type for the DocumentService type
type DocumentService struct {
	mock.Mock
}

// CreateInBundle provides a mock function with given fields: ctx, appID, bundleID, in
func (_m *DocumentService) CreateInBundle(ctx context.Context, appID string, bundleID string, in model.DocumentInput) (string, error) {
	ret := _m.Called(ctx, appID, bundleID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, model.DocumentInput) string); ok {
		r0 = rf(ctx, appID, bundleID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, model.DocumentInput) error); ok {
		r1 = rf(ctx, appID, bundleID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForBundle provides a mock function with given fields: ctx, id, bundleID
func (_m *DocumentService) GetForBundle(ctx context.Context, id string, bundleID string) (*model.Document, error) {
	ret := _m.Called(ctx, id, bundleID)

	var r0 *model.Document
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Document); ok {
		r0 = rf(ctx, id, bundleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Document)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, bundleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByBundleIDs provides a mock function with given fields: ctx, bundleIDs, pageSize, cursor
func (_m *DocumentService) ListByBundleIDs(ctx context.Context, bundleIDs []string, pageSize int, cursor string) ([]*model.DocumentPage, error) {
	ret := _m.Called(ctx, bundleIDs, pageSize, cursor)

	var r0 []*model.DocumentPage
	if rf, ok := ret.Get(0).(func(context.Context, []string, int, string) []*model.DocumentPage); ok {
		r0 = rf(ctx, bundleIDs, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DocumentPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, int, string) error); ok {
		r1 = rf(ctx, bundleIDs, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
