package destinationfetchersvc

import (
	"context"
)

type fetcher struct {
	svc DestinationService
}

// NewFetcher creates new fetcher
func NewFetcher(svc DestinationService) *fetcher {
	return &fetcher{svc: svc}
}

func (f *fetcher) FetchDestinationsOnDemand(ctx context.Context, tenantID, parentTenantID string) error {
	f.svc.SyncDestinations()
	return nil
}
