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
	FetchDestinationsOnDemand(ctx context.Context, tenantID, parentTenantID string) error
}

// NewDestinationsHTTPHandler returns a new HTTP handler, responsible for handleing HTTP requests
func NewDestinationsHTTPHandler(fetcher DestinationFetcher, config HandlerConfig) *handler {
	return &handler{fetcher: fetcher}
}

func (h *handler) FetchDestinationsOnDemand(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.fetcher.FetchDestinationsOnDemand(ctx, "tenant id", "parent tenant id")
	writer.WriteHeader(http.StatusOK)
}
