package repo

import (
	"context"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
)

// UnionLister is an interface for listing tenant scoped entities with either externally managed tenant accesses (m2m table or view) or embedded tenant in them.
// It lists entities based on multiple parent queries. For each parent a separate list with a separate tenant isolation subquery is created and the end result is a union of all the results.
type UnionLister interface {
	// List stores the result into dest and returns the total count of tuples for each id from ids
	List(ctx context.Context, resourceType resource.Type, tenant string, ids []string, idsColumn string, pageSize int, cursor string, orderBy OrderByParams, dest Collection, additionalConditions ...Condition) (map[string]int, error)
}

// UnionListerGlobal is an interface for listing global entities.
// It lists entities based on multiple parent queries. For each parent a separate list with a separate tenant isolation subquery is created and the end result is a union of all the results.
type UnionListerGlobal interface {
	ListGlobal(ctx context.Context, ids []string, idsColumn string, pageSize int, cursor string, orderBy OrderByParams, dest Collection, additionalConditions ...Condition) (map[string]int, error)
	SetSelectedColumns(selectedColumns []string)
	Clone() *unionLister
}

type unionLister struct {
	tableName       string
	selectedColumns string
	tenantColumn    *string
	resourceType    resource.Type
}

// NewUnionListerWithEmbeddedTenant is a constructor for UnionLister about entities with tenant embedded in them.
func NewUnionListerWithEmbeddedTenant(tableName string, tenantColumn string, selectedColumns []string) UnionLister {
	return &unionLister{
		tableName:       tableName,
		selectedColumns: strings.Join(selectedColumns, ", "),
		tenantColumn:    &tenantColumn,
	}
}

// NewUnionLister is a constructor for UnionLister about entities with externally managed tenant accesses (m2m table or view)
func NewUnionLister(tableName string, selectedColumns []string) UnionLister {
	return &unionLister{
		tableName:       tableName,
		selectedColumns: strings.Join(selectedColumns, ", "),
	}
}

// NewUnionListerGlobal is a constructor for UnionListerGlobal about global entities.
func NewUnionListerGlobal(resourceType resource.Type, tableName string, selectedColumns []string) UnionListerGlobal {
	return &unionLister{
		tableName:       tableName,
		selectedColumns: strings.Join(selectedColumns, ", "),
		resourceType:    resourceType,
	}
}

// SetSelectedColumns sets the selected columns for the lister
func (l *unionLister) SetSelectedColumns(selectedColumns []string) {
	l.selectedColumns = strings.Join(selectedColumns, ", ")
}

// Clone returns a copy of the lister
func (l *unionLister) Clone() *unionLister {
	var clonedLister unionLister

	clonedLister.resourceType = l.resourceType
	clonedLister.tableName = l.tableName
	clonedLister.selectedColumns = l.selectedColumns
	clonedLister.tenantColumn = l.tenantColumn

	return &clonedLister
}

// List lists tenant scoped entities based on multiple parent queries. For each parent a separate list with a separate tenant isolation subquery is created and the end result is a union of all the results.
// If the tenantColumn is configured the isolation is based on equal condition on tenantColumn.
// If the tenantColumn is not configured an entity with externally managed tenant accesses in m2m table / view is assumed.
func (l *unionLister) List(ctx context.Context, resourceType resource.Type, tenant string, ids []string, idscolumn string, pageSize int, cursor string, orderBy OrderByParams, dest Collection, additionalConditions ...Condition) (map[string]int, error) {
	if tenant == "" {
		return nil, apperrors.NewTenantRequiredError()
	}

	if l.tenantColumn != nil {
		additionalConditions = append(Conditions{NewEqualCondition(*l.tenantColumn, tenant)}, additionalConditions...)
		return l.list(ctx, resourceType, pageSize, cursor, orderBy, ids, idscolumn, dest, additionalConditions...)
	}

	tenantIsolation, err := NewTenantIsolationCondition(resourceType, tenant, false)
	if err != nil {
		return nil, err
	}

	additionalConditions = append(additionalConditions, tenantIsolation)

	return l.list(ctx, resourceType, pageSize, cursor, orderBy, ids, idscolumn, dest, additionalConditions...)
}

// ListGlobal lists global entities without tenant isolation.
func (l *unionLister) ListGlobal(ctx context.Context, ids []string, idscolumn string, pageSize int, cursor string, orderBy OrderByParams, dest Collection, additionalConditions ...Condition) (map[string]int, error) {
	return l.list(ctx, l.resourceType, pageSize, cursor, orderBy, ids, idscolumn, dest, additionalConditions...)
}

type queryStruct struct {
	args      []interface{}
	statement string
}

func (l *unionLister) list(ctx context.Context, resourceType resource.Type, pageSize int, cursor string, orderBy OrderByParams, ids []string, idsColumn string, dest Collection, conditions ...Condition) (map[string]int, error) {
	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return nil, err
	}

	offset, err := pagination.DecodeOffsetCursor(cursor)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding page cursor")
	}

	queries, err := l.buildQueries(ids, idsColumn, conditions, orderBy, pageSize, offset)
	if err != nil {
		return nil, err
	}

	stmts := make([]string, 0, len(queries))
	for _, q := range queries {
		stmts = append(stmts, q.statement)
	}

	args := make([]interface{}, 0, len(queries))
	for _, q := range queries {
		args = append(args, q.args...)
	}

	query := buildUnionQuery(stmts)

	err = persist.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return nil, persistence.MapSQLError(ctx, err, resourceType, resource.List, "while fetching list page of objects from '%s' table", l.tableName)
	}

	totalCount, err := l.getTotalCount(ctx, resourceType, persist, idsColumn, []string{idsColumn}, OrderByParams{NewAscOrderBy(idsColumn)}, conditions)
	if err != nil {
		return nil, err
	}

	return totalCount, nil
}

func (l *unionLister) buildQueries(ids []string, idsColumn string, conditions []Condition, orderBy OrderByParams, limit int, offset int) ([]queryStruct, error) {
	queries := make([]queryStruct, 0, len(ids))
	for _, id := range ids {
		query, args, err := buildSelectQueryWithLimitAndOffset(l.tableName, l.selectedColumns, append(conditions, NewEqualCondition(idsColumn, id)), orderBy, limit, offset, false)
		if err != nil {
			return nil, errors.Wrap(err, "while building list query")
		}

		queries = append(queries, queryStruct{
			args:      args,
			statement: query,
		})
	}
	return queries, nil
}

type idToCount struct {
	ID    string `db:"id"`
	Count int    `db:"total_count"`
}

func (l *unionLister) getTotalCount(ctx context.Context, resourceType resource.Type, persist persistence.PersistenceOp, idsColumn string, groupBy GroupByParams, orderBy OrderByParams, conditions Conditions) (map[string]int, error) {
	query, args, err := buildCountQuery(l.tableName, idsColumn, conditions, groupBy, orderBy, true)
	if err != nil {
		return nil, err
	}

	var counts []idToCount
	err = persist.SelectContext(ctx, &counts, query, args...)
	if err != nil {
		return nil, persistence.MapSQLError(ctx, err, resourceType, resource.List, "while counting objects from '%s' table", l.tableName)
	}

	totalCount := make(map[string]int)
	for _, c := range counts {
		totalCount[c.ID] = c.Count
	}

	return totalCount, nil
}
