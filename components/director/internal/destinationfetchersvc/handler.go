package destinationfetchersvc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
)

const (
	tenantIDKey        = "subaccountId"
	regionKey          = "region"
	destQueryParameter = "dest"
)

type HandlerConfig struct {
	DestinationsEndpoint          string `envconfig:"APP_DESTINATIONS_ON_DEMAND_HANDLER_ENDPOINT,default=/v1/fetch"`
	DestinationsSensitiveEndpoint string `envconfig:"APP_DESTINATIONS_GET_DESTINATION,default=/v1/info"`
	UserContextHeader             string `envconfig:"APP_USER_CONTEXT_HEADER,default=user_context"`
}

type handler struct {
	fetcher DestinationManager
	config  HandlerConfig
}

//go:generate mockery --name=DestinationManager --output=automock --outpkg=automock --case=underscore --disable-version-string
type DestinationManager interface {
	SyncTenantDestinations(ctx context.Context, tenantID string) error
	FetchDestinationsSensitiveData(ctx context.Context, tenantID string, destinationNames []string) ([]byte, error)
}

// NewDestinationsHTTPHandler returns a new HTTP handler, responsible for handling HTTP requests
func NewDestinationsHTTPHandler(fetcher DestinationManager, config HandlerConfig) *handler {
	return &handler{
		fetcher: fetcher,
		config:  config,
	}
}

func (h *handler) SyncTenantDestinations(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userContextHeader := request.Header.Get(h.config.UserContextHeader)
	tenantID, err := h.readTenantFromHeader(userContextHeader)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.fetcher.SyncTenantDestinations(ctx, tenantID); err != nil {
		if apperrors.IsNotFoundError(err) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(writer, fmt.Sprintf("Failed to fetch destinations for tenant %s",
			tenantID), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func getDestinationNames(namesRaw string) ([]string, error) {
	namesRawLength := len(namesRaw)
	if namesRawLength == 0 {
		return nil, fmt.Errorf("dest query parameter is missing")
	}

	if namesRaw[0] != '[' || namesRaw[namesRawLength-1] != ']' {
		return nil, fmt.Errorf("%s dest query parameter is invalid. Must start with '[' and end with ']'", namesRaw)
	}

	//removes brackets from query
	namesRawWithoutBrackets := namesRaw[1 : namesRawLength-1]
	names := strings.Split(namesRawWithoutBrackets, ",")

	if sliceContainsEmptyString(names) {
		return nil, fmt.Errorf("name parameter containes empty element")
	}

	return names, nil
}

func (h *handler) FetchDestinationsSensitiveData(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userContextHeader := request.Header.Get(h.config.UserContextHeader)
	tenantID, err := h.readTenantFromHeader(userContextHeader)
	if err != nil {
		log.C(ctx).Errorf("Failed to read userContext header with error: %s", err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	namesRaw := request.URL.Query().Get(destQueryParameter)
	names, err := getDestinationNames(namesRaw)
	if err != nil {
		log.C(ctx).Errorf("Failed to fetch sensitive data for destinations %s: %v", namesRaw, err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	json, err := h.fetcher.FetchDestinationsSensitiveData(ctx, tenantID, names)

	if err != nil {
		log.C(ctx).Errorf("Failed to fetch sensitive data for destinations %s and tenant %s: %v",
			namesRaw, tenantID, err)
		if apperrors.IsNotFoundError(err) {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(writer, fmt.Sprintf("Failed to fetch sensitive data for destinations %s", namesRaw), http.StatusInternalServerError)
		return
	}

	writer.Write(json)
	writer.WriteHeader(http.StatusOK)
}

func sliceContainsEmptyString(s []string) bool {
	for _, e := range s {
		if strings.TrimSpace(e) == "" {
			return true
		}
	}

	return false
}

func (h *handler) readTenantFromHeader(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("%s header is missing", h.config.UserContextHeader)
	}

	var headerMap map[string]string
	if err := json.Unmarshal([]byte(header), &headerMap); err != nil {
		return "", fmt.Errorf("failed to parse %s header", h.config.UserContextHeader)
	}

	tenantId, ok := headerMap[tenantIDKey]
	if !ok {
		return "", fmt.Errorf("%s not found in %s header", tenantIDKey, h.config.UserContextHeader)
	}

	return tenantId, nil
}
