// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// FormationTemplateRepository is an autogenerated mock type for the FormationTemplateRepository type
type FormationTemplateRepository struct {
	mock.Mock
}

// GetByName provides a mock function with given fields: ctx, templateName
func (_m *FormationTemplateRepository) GetByName(ctx context.Context, templateName string) (*model.FormationTemplate, error) {
	ret := _m.Called(ctx, templateName)

	var r0 *model.FormationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.FormationTemplate); ok {
		r0 = rf(ctx, templateName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, templateName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFormationTemplateRepository creates a new instance of FormationTemplateRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationTemplateRepository(t testing.TB) *FormationTemplateRepository {
	mock := &FormationTemplateRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
