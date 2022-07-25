package destinationfetchersvc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	GetDestinationsInfo(ctx context.Context, subaccountID string, destinationNames []string) ([]byte, error)
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

func (h *handler) GetDestinationInfo(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	subaccountID := request.Header.Get(h.config.UserContextHeader)
	if subaccountID == "" {
		http.Error(writer, fmt.Sprintf("%s header is missing", h.config.UserContextHeader), http.StatusBadRequest)
		return
	}

	namesRaw := request.URL.Query().Get("name")
	//removes brackets from query
	namesRawWithoutBrackets := namesRaw[1 : len(namesRaw)-1]
	names := strings.Split(namesRawWithoutBrackets, ",")

	json, err := h.fetcher.GetDestinationsInfo(ctx, subaccountID, names)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get destination info for names %s", namesRawWithoutBrackets), http.StatusBadRequest)
		return
	}

	writer.Write(json)
	writer.WriteHeader(http.StatusOK)
}
