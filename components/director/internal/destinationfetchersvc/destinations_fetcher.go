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

func (f *fetcher) FetchDestinationsOnDemand(ctx context.Context, subaccountID string) error {
	return f.svc.SyncSubaccountDestinations(ctx, subaccountID)
}

func (f *fetcher) GetDestinationsInfo(ctx context.Context, subaccountID string, names []string) ([]byte, error) {
	return f.svc.GetDestinationsInfo(ctx, subaccountID, names)
}
