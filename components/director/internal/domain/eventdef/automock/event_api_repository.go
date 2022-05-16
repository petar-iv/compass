// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// EventAPIRepository is an autogenerated mock type for the EventAPIRepository type
type EventAPIRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, tenant, item
func (_m *EventAPIRepository) Create(ctx context.Context, tenant string, item *model.EventDefinition) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.EventDefinition) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tenantID, id
func (_m *EventAPIRepository) Delete(ctx context.Context, tenantID string, id string) error {
	ret := _m.Called(ctx, tenantID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenantID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllByBundleID provides a mock function with given fields: ctx, tenantID, bundleID
func (_m *EventAPIRepository) DeleteAllByBundleID(ctx context.Context, tenantID string, bundleID string) error {
	ret := _m.Called(ctx, tenantID, bundleID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenantID, bundleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, tenantID, id
func (_m *EventAPIRepository) GetByID(ctx context.Context, tenantID string, id string) (*model.EventDefinition, error) {
	ret := _m.Called(ctx, tenantID, id)

	var r0 *model.EventDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.EventDefinition); ok {
		r0 = rf(ctx, tenantID, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.EventDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForBundle provides a mock function with given fields: ctx, tenant, id, bundleID
func (_m *EventAPIRepository) GetForBundle(ctx context.Context, tenant string, id string, bundleID string) (*model.EventDefinition, error) {
	ret := _m.Called(ctx, tenant, id, bundleID)

	var r0 *model.EventDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *model.EventDefinition); ok {
		r0 = rf(ctx, tenant, id, bundleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.EventDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, tenant, id, bundleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByApplicationID provides a mock function with given fields: ctx, tenantID, appID
func (_m *EventAPIRepository) ListByApplicationID(ctx context.Context, tenantID string, appID string) ([]*model.EventDefinition, error) {
	ret := _m.Called(ctx, tenantID, appID)

	var r0 []*model.EventDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*model.EventDefinition); ok {
		r0 = rf(ctx, tenantID, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.EventDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantID, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByBundleIDs provides a mock function with given fields: ctx, tenantID, bundleIDs, bundleRefs, totalCounts, pageSize, cursor
func (_m *EventAPIRepository) ListByBundleIDs(ctx context.Context, tenantID string, bundleIDs []string, bundleRefs []*model.BundleReference, totalCounts map[string]int, pageSize int, cursor string) ([]*model.EventDefinitionPage, error) {
	ret := _m.Called(ctx, tenantID, bundleIDs, bundleRefs, totalCounts, pageSize, cursor)

	var r0 []*model.EventDefinitionPage
	if rf, ok := ret.Get(0).(func(context.Context, string, []string, []*model.BundleReference, map[string]int, int, string) []*model.EventDefinitionPage); ok {
		r0 = rf(ctx, tenantID, bundleIDs, bundleRefs, totalCounts, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.EventDefinitionPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []string, []*model.BundleReference, map[string]int, int, string) error); ok {
		r1 = rf(ctx, tenantID, bundleIDs, bundleRefs, totalCounts, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, tenant, item
func (_m *EventAPIRepository) Update(ctx context.Context, tenant string, item *model.EventDefinition) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.EventDefinition) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEventAPIRepository creates a new instance of EventAPIRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventAPIRepository(t testing.TB) *EventAPIRepository {
	mock := &EventAPIRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
