package destination

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
)

const destinationTable = "public.destinations"

var (
	destinationColumns = []string{"id", "name", "type", "url", "authentication", "tenant_id", "bundle_id", "revision"}
	conflictingColumns = []string{"tenant_id", "bundle_id"}
	updateColumns      = []string{"name", "type", "url", "authentication", "revision"}
)

type repository struct {
	deleterGlobal  repo.DeleterGlobal
	upserterGlobal repo.UpserterGlobal
}

func NewRepository() *repository {
	return &repository{
		deleterGlobal:  repo.NewDeleterGlobal(resource.Destination, destinationTable),
		upserterGlobal: repo.NewUpserterGlobal(resource.Destination, destinationTable, destinationColumns, conflictingColumns, updateColumns),
	}
}

func (r *repository) Upsert(ctx context.Context) error {
	return nil
}

func (r *repository) Delete(ctx context.Context) error {
	fmt.Println("in delete")
	conditions := repo.Conditions{repo.NewEqualCondition("name", "test")}
	r.deleterGlobal.DeleteManyGlobal(ctx, conditions)
	return nil
}
