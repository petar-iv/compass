package destinationfetchersvc

import (
	"context"
	"net/http"
)

type HandlerConfig struct {
	DestinationsEndpoint string `envconfig:"default=/fetch"`
}

type handler struct {
	fetcher DestinationFetcher
}

type DestinationFetcher interface {
	GetDestinations(ctx context.Context, tenantID, parentTenantID string) error
}

// NewDestinationFetcherHTTPHandler returns a new HTTP handler, responsible for handleing HTTP requests
func NewDestinationFetcherHTTPHandler(config HandlerConfig) *handler {
	return &handler{fetcher: nil}
}

func (h *handler) GetDestinations(writer http.ResponseWriter, request *http.Request) {
	//ctx := request.Context()
	writer.WriteHeader(http.StatusOK)
}
