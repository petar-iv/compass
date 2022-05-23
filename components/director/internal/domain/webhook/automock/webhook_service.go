// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// WebhookService is an autogenerated mock type for the WebhookService type
type WebhookService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, resourceID, in, objectType
func (_m *WebhookService) Create(ctx context.Context, resourceID string, in model.WebhookInput, objectType model.WebhookReferenceObjectType) (string, error) {
	ret := _m.Called(ctx, resourceID, in, objectType)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, model.WebhookInput, model.WebhookReferenceObjectType) string); ok {
		r0 = rf(ctx, resourceID, in, objectType)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.WebhookInput, model.WebhookReferenceObjectType) error); ok {
		r1 = rf(ctx, resourceID, in, objectType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id, objectType
func (_m *WebhookService) Delete(ctx context.Context, id string, objectType model.WebhookReferenceObjectType) error {
	ret := _m.Called(ctx, id, objectType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.WebhookReferenceObjectType) error); ok {
		r0 = rf(ctx, id, objectType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id, objectType
func (_m *WebhookService) Get(ctx context.Context, id string, objectType model.WebhookReferenceObjectType) (*model.Webhook, error) {
	ret := _m.Called(ctx, id, objectType)

	var r0 *model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string, model.WebhookReferenceObjectType) *model.Webhook); ok {
		r0 = rf(ctx, id, objectType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.WebhookReferenceObjectType) error); ok {
		r1 = rf(ctx, id, objectType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllApplicationWebhooks provides a mock function with given fields: ctx, applicationID
func (_m *WebhookService) ListAllApplicationWebhooks(ctx context.Context, applicationID string) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, applicationID)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Webhook); ok {
		r0 = rf(ctx, applicationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, applicationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForRuntime provides a mock function with given fields: ctx, runtimeID
func (_m *WebhookService) ListForRuntime(ctx context.Context, runtimeID string) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, runtimeID)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Webhook); ok {
		r0 = rf(ctx, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, in, objectType
func (_m *WebhookService) Update(ctx context.Context, id string, in model.WebhookInput, objectType model.WebhookReferenceObjectType) error {
	ret := _m.Called(ctx, id, in, objectType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.WebhookInput, model.WebhookReferenceObjectType) error); ok {
		r0 = rf(ctx, id, in, objectType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWebhookService creates a new instance of WebhookService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebhookService(t testing.TB) *WebhookService {
	mock := &WebhookService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
