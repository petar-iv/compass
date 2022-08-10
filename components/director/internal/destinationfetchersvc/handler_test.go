package destinationfetchersvc_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc"
	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc/automock"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_FetchDestinationsOnDemand(t *testing.T) {
	const (
		subaccountID      = "f09ba084-0e82-49ab-ab2e-b7ecc988312d"
		userContextHeader = "user_context"
	)

	target := "/v1/fetch"

	validHandlerConfig := destinationfetchersvc.HandlerConfig{
		DestinationsEndpoint:          "/v1/fetch",
		DestinationsSensitiveEndpoint: "/v1/info",
		UserContextHeader:             userContextHeader,
	}

	reqWithUserContext := httptest.NewRequest(http.MethodPut, target, nil)
	userContext := `{"subaccountId":"` + subaccountID + `"}`
	reqWithUserContext.Header.Set(userContextHeader, userContext)
	testCases := []struct {
		Name                string
		Request             *http.Request
		DestinationManager  func() *automock.DestinationManager
		ExpectedErrorOutput string
		ExpectedStatusCode  int
	}{
		{
			Name:    "Successful fetch on-demand",
			Request: reqWithUserContext,
			DestinationManager: func() *automock.DestinationManager {
				svc := &automock.DestinationManager{}
				svc.On("SyncSubaccountDestinations", mock.Anything, subaccountID).Return(nil)
				return svc
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:    "Missing user context header",
			Request: httptest.NewRequest(http.MethodPut, target, nil),
			DestinationManager: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:    "Subaccount not found",
			Request: reqWithUserContext,
			DestinationManager: func() *automock.DestinationManager {
				svc := &automock.DestinationManager{}
				err := apperrors.NewNotFoundErrorWithMessage(resource.Label,
					subaccountID, fmt.Sprintf("subaccount %s not found", subaccountID))
				svc.On("SyncSubaccountDestinations", mock.Anything, subaccountID).Return(err)
				return svc
			},
			ExpectedErrorOutput: fmt.Sprintf("subaccount %s not found", subaccountID),
			ExpectedStatusCode:  http.StatusBadRequest,
		},
		{
			Name:    "Internal Server Error",
			Request: reqWithUserContext,
			DestinationManager: func() *automock.DestinationManager {
				svc := &automock.DestinationManager{}
				err := fmt.Errorf("Random error")
				svc.On("SyncSubaccountDestinations", mock.Anything, subaccountID).Return(err)
				return svc
			},
			ExpectedErrorOutput: fmt.Sprintf("Failed to fetch destinations for subaccount %s", subaccountID),
			ExpectedStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tf := testCase.DestinationManager()
			defer mock.AssertExpectationsForObjects(t, tf)

			handler := destinationfetchersvc.NewDestinationsHTTPHandler(tf, validHandlerConfig)
			req := testCase.Request
			w := httptest.NewRecorder()

			// WHEN
			handler.FetchDestinationsOnDemand(w, req)

			// THEN
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			if len(testCase.ExpectedErrorOutput) > 0 {
				assert.Contains(t, string(body), testCase.ExpectedErrorOutput)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.ExpectedStatusCode, resp.StatusCode)
		})
	}
}

func TestService_FetchDestinationsSensitiveData(t *testing.T) {
	const (
		subaccountID       = "f09ba084-0e82-49ab-ab2e-b7ecc988312d"
		userContextHeader  = "user_context"
		destQueryParameter = "dest"
	)

	json := []byte("{}")

	target := "/v1/info"

	validHandlerConfig := destinationfetchersvc.HandlerConfig{
		DestinationsEndpoint:          "/v1/fetch",
		DestinationsSensitiveEndpoint: "/v1/info",
		UserContextHeader:             userContextHeader,
	}

	namesRaw := "[Rand, Mat]"
	names := strings.Split("Rand, Mat", ",")
	reqWithUserContext := httptest.NewRequest(http.MethodPut, target, nil)
	userContext := `{"subaccountId":"` + subaccountID + `"}`
	reqWithUserContext.Header.Set(userContextHeader, userContext)
	testCases := []struct {
		Name                  string
		Request               *http.Request
		DestQueryParameter    string
		DestinationFetcherSvc func() *automock.DestinationManager
		ExpectedErrorOutput   string
		ExpectedStatusCode    int
	}{
		{
			Name:               "Successful fetch data fetch",
			Request:            reqWithUserContext,
			DestQueryParameter: namesRaw,
			DestinationFetcherSvc: func() *automock.DestinationManager {
				svc := &automock.DestinationManager{}
				svc.On("FetchDestinationsSensitiveData", mock.Anything, subaccountID, names).
					Return(
						func(ctx context.Context, subaccountID string, destNames []string) []byte {
							return json
						},
						func(ctx context.Context, subaccountID string, destNames []string) error {
							return nil
						},
					)
				return svc
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:    "Missing user context header",
			Request: httptest.NewRequest(http.MethodPut, target, nil),
			DestinationFetcherSvc: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:    "Missing destination query parameter.",
			Request: reqWithUserContext,
			DestinationFetcherSvc: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid destination query parameter. Missing beggining bracket",
			Request:            reqWithUserContext,
			DestQueryParameter: "Rand,Mat]",
			DestinationFetcherSvc: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid destination query parameter. Missing end bracket",
			Request:            reqWithUserContext,
			DestQueryParameter: "[Perin,Mat",
			DestinationFetcherSvc: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid destination query parameter. Empty element",
			Request:            reqWithUserContext,
			DestQueryParameter: "[Perin,]",
			DestinationFetcherSvc: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid destination query parameter. No bracket",
			Request:            reqWithUserContext,
			DestQueryParameter: "Rand,Perin",
			DestinationFetcherSvc: func() *automock.DestinationManager {
				return &automock.DestinationManager{}
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tf := testCase.DestinationFetcherSvc()
			defer mock.AssertExpectationsForObjects(t, tf)

			handler := destinationfetchersvc.NewDestinationsHTTPHandler(tf, validHandlerConfig)
			req := testCase.Request
			//req is a pointer and the changes on the previous test are kept
			req.URL.RawQuery = ""
			if len(testCase.DestQueryParameter) > 0 {
				query := req.URL.Query()
				query.Add(destQueryParameter, testCase.DestQueryParameter)
				req.URL.RawQuery = query.Encode()
			}

			w := httptest.NewRecorder()

			// WHEN
			handler.FetchDestinationsSensitiveData(w, req)

			// THEN
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			if len(testCase.ExpectedErrorOutput) > 0 {
				assert.Contains(t, string(body), testCase.ExpectedErrorOutput)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.ExpectedStatusCode, resp.StatusCode)
		})
	}
}
