package destination

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
)

const (
	destinationTable = "public.destinations"
	revisionColumn   = "revision"
)

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
	fmt.Println("#### UPSERT")
	return nil
}

func (r *repository) Delete(ctx context.Context, revision string) error {
	conditions := repo.Conditions{repo.NewNotEqualCondition(revisionColumn, revision)}
	r.deleterGlobal.DeleteManyGlobal(ctx, conditions)
	return nil
}

func (r *repository) GetSubdomain(ctx context.Context, subaccountId string) (*Subdomain, error) {
	var subdomain Subdomain

	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
	SELECT l.tenant_id, l.value #>> '{}' as value
	FROM labels l
	WHERE l.key='subdomain' AND l.tenant_id=(
	SELECT id FROM business_tenant_mappings WHERE external_tenant='%s'
	)`, subaccountId)

	err = persist.GetContext(ctx, &subdomain, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch subdomain for subaccount with id '%s' from DB", subaccountId)
	}
	return &subdomain, nil
}

