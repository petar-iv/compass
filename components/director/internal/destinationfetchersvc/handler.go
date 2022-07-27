package destinationfetchersvc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
)

type DestinationsConfig struct {
	OAuthConfig OAuth2Config
}

type HandlerConfig struct {
	DestinationsEndpoint     string `envconfig:"APP_DESTINATIONS_ON_DEMAND_HANDLER_ENDPOINT,default=/v1/fetch"`
	DestinationsInfoEndpoint string `envconfig:"APP_DESTINATIONS_GET_DESTINATION,default=/v1/info"`
	UserContextHeader        string `envconfig:"APP_USER_CONTEXT_HEADER,default=user_context"`
}

type handler struct {
	fetcher DestinationFetcher
	config  HandlerConfig
}

type DestinationFetcher interface {
	FetchDestinationsOnDemand(ctx context.Context, subaccountID string) error
	FetchDestinationsSensitiveData(ctx context.Context, subaccountID string, destinationNames []string) ([]byte, error)
}

// NewDestinationsHTTPHandler returns a new HTTP handler, responsible for handleing HTTP requests
func NewDestinationsHTTPHandler(fetcher DestinationFetcher, config HandlerConfig) *handler {
	return &handler{
		fetcher: fetcher,
		config:  config,
	}
}

func (h *handler) FetchDestinationsOnDemand(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	subaccountID := request.Header.Get(h.config.UserContextHeader)
	if subaccountID == "" {
		http.Error(writer, fmt.Sprintf("%s header is missing", h.config.UserContextHeader), http.StatusBadRequest)
		return
	}

	if err := h.fetcher.FetchDestinationsOnDemand(ctx, subaccountID); err != nil {
		http.Error(writer, fmt.Sprintf("Failed to fetch destinations for subaccount %s", subaccountID), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func getDestinationNames(namesRaw string) ([]string, error) {
	namesRawLength := len(namesRaw)
	if namesRawLength == 0 {
		return nil, fmt.Errorf("name query parameter is missing")
	}

	if namesRaw[0] != '[' || namesRaw[namesRawLength-1] != ']' {
		return nil, fmt.Errorf("%s name query parameter is invalid. Must start with '[' and end with ']'", namesRaw)
	}

	//removes brackets from query
	namesRawWithoutBrackets := namesRaw[1 : namesRawLength-1]
	names := strings.Split(namesRawWithoutBrackets, ",")

	if err := sliceContainsEmptyString(names); err != nil {
		return nil, fmt.Errorf("name parameter containes empty element")
	}

	return names, nil
}

func (h *handler) FetchDestinationsSensitiveData(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	subaccountID := request.Header.Get(h.config.UserContextHeader)
	if subaccountID == "" {
		http.Error(writer, fmt.Sprintf("%s header is missing", h.config.UserContextHeader), http.StatusBadRequest)
		return
	}

	namesRaw := request.URL.Query().Get("name")
	names, err := getDestinationNames(namesRaw)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	json, err := h.fetcher.FetchDestinationsSensitiveData(ctx, subaccountID, names)

	if err != nil {
		log.C(ctx).Errorf("Failed to fetch destination sensitive data with error: %s", err.Error())
		if apperrors.IsNotFoundError(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(writer, fmt.Sprintf("Failed to get destination info for names %s", namesRaw), http.StatusInternalServerError)
		return
	}

	writer.Write(json)
	writer.WriteHeader(http.StatusOK)
}

func sliceContainsEmptyString(s []string) error {
	for _, e := range s {
		if strings.TrimSpace(e) == "" {
			return fmt.Errorf("contains blank element")
		}
	}

	return nil
}
