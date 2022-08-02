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

func (f *fetcher) FetchDestinationsOnDemand(ctx context.Context, userContext *UserContext) error {
	return f.svc.SyncSubaccountDestinations(ctx, userContext)
}

func (f *fetcher) FetchDestinationsSensitiveData(ctx context.Context, userContext *UserContext, names []string) ([]byte, error) {
	return f.svc.FetchDestinationsSensitiveData(ctx, userContext, names)
}
