package destinationfetchersvc

import (
	"context"
	"fmt"
	"net/http"
)

type DestinationsConfig struct {
	OAuthConfig OAuth2Config
}

type HandlerConfig struct {
	DestinationsEndpoint string `envconfig:"default=/fetch"`
}

type handler struct {
	fetcher DestinationFetcher
}

type DestinationFetcher interface {
	FetchDestinationsOnDemand(ctx context.Context, subaccountID string) error
}

// NewDestinationsHTTPHandler returns a new HTTP handler, responsible for handleing HTTP requests
func NewDestinationsHTTPHandler(fetcher DestinationFetcher, config HandlerConfig) *handler {
	return &handler{fetcher: fetcher}
}

func (h *handler) FetchDestinationsOnDemand(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	subaccountIdHeader := "x-tenant"
	subaccountID := request.Header.Get(subaccountIdHeader)
	if subaccountID == "" {
		http.Error(writer, fmt.Sprintf("%s header is missing", subaccountIdHeader), http.StatusBadRequest)
		return
	}

	if err := h.fetcher.FetchDestinationsOnDemand(ctx, subaccountID); err != nil {
		http.Error(writer, fmt.Sprintf("Failed to fetch destinations for subaccount %s", subaccountID), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
