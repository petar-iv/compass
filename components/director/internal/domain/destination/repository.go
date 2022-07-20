package destination

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/internal/repo"
)

const destinationTable = "public.destinations"

var (
	destinationColumns = []string{"id", "name", "type", "url", "authentication", "tenant_id", "bundle_id", "revision"}
	conflictingColumns = []string{"id"}
	updateColumns      = []string{"name", "type", "url", "authentication", "revision"}
)

type repository struct {
	deleter  repo.Deleter
	upserter repo.Upserter
}

func NewRepository() *repository {
	return &repository{
		deleter:  repo.NewDeleter(destinationTable),
		upserter: repo.NewUpserter(destinationTable, destinationColumns, conflictingColumns, updateColumns),
	}
}

func (r *repository) Upsert(ctx context.Context) error {
	return nil
}

func (r *repository) Delete(ctx context.Context) error {
	return nil
}
