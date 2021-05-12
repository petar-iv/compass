package solution

import (
	"database/sql"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

// Solution struct represents database entity for Solution
type Solution struct {
	ID          string         `db:"id"`
	TenantID    string         `db:"tenant_id"`
	Name        string         `db:"name"`
	Version     string         `db:"version"`
	Description sql.NullString `db:"description"`
}

// EntityFromSolutionModel converts Solution model to Solution entity
func EntityFromSolutionModel(model *model.Solution) (*Solution, error) {
	var nullDescription sql.NullString
	if model.Description != nil && len(*model.Description) > 0 {
		nullDescription = sql.NullString{
			String: *model.Description,
			Valid:  true,
		}
	}

	return &Solution{
		ID:          model.ID,
		Name:        model.Name,
		TenantID:    model.Tenant,
		Version:     model.Version,
		Description: nullDescription,
	}, nil
}

// ToModel converts Solution entity to Solution model
func (e Solution) ToModel() *model.Solution {
	var description *string
	if e.Description.Valid {
		description = new(string)
		*description = e.Description.String
	}

	return &model.Solution{
		ID:          e.ID,
		Tenant:      e.TenantID,
		Name:        e.Name,
		Version:     e.Version,
		Description: description,
	}
}
