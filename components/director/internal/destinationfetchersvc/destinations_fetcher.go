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

func (f *fetcher) FetchDestinationsOnDemand(ctx context.Context, subaccountID string, region string) error {
	return f.svc.SyncSubaccountDestinations(ctx, subaccountID, region)
}

func (f *fetcher) FetchDestinationsSensitiveData(ctx context.Context, subaccountID string, region string, names []string) ([]byte, error) {
	return f.svc.FetchDestinationsSensitiveData(ctx, subaccountID, region, names)
}
