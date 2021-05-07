package bundleinstanceauth

import (
	"database/sql"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/repo"
)

type Entity struct {
	*repo.BaseEntity
	BundleID        string         `db:"bundle_id"`
	TenantID        string         `db:"tenant_id"`
	Context         sql.NullString `db:"context"`
	InputParams     sql.NullString `db:"input_params"`
	AuthValue       sql.NullString `db:"auth_value"`
	StatusCondition string         `db:"status_condition"`
	StatusTimestamp time.Time      `db:"status_timestamp"`
	StatusMessage   string         `db:"status_message"`
	StatusReason    string         `db:"status_reason"`
}

type Collection []Entity

func (c Collection) Len() int {
	return len(c)
}
