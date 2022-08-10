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

func (f *fetcher) SyncTenantDestinations(ctx context.Context, tenantID string) error {
	return f.svc.SyncTenantDestinations(ctx, tenantID)
}

func (f *fetcher) FetchDestinationsSensitiveData(ctx context.Context, tenantID string, names []string) ([]byte, error) {
	return f.svc.FetchDestinationsSensitiveData(ctx, tenantID, names)
}
