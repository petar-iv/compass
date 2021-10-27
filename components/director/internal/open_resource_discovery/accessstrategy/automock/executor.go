// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, client, url
func (_m *Executor) Execute(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	ret := _m.Called(ctx, client, url)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, *http.Client, string) *http.Response); ok {
		r0 = rf(ctx, client, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *http.Client, string) error); ok {
		r1 = rf(ctx, client, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}