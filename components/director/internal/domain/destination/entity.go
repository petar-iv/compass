package destination

import "database/sql"

type Entity struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	Type           string `db:"type"`
	URL            string `db:"url"`
	Authentication string `db:"authentication"`
	TenantID       string `db:"tenant_id"`
	BundleID       string `db:"bundle_id"`
	Revision       string `db:"revision"`
}

type Subdomain struct {
	TenantID       string         `db:"tenant_id"`
	ParentTenantID sql.NullString `db:"parent"`
	Value          string         `db:"value"`
}

type Bundle struct {
	ID string `db:"id"`
}
