package solution

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
)

const solutionTable string = `public.solutions`

var (
	solutionColumns = []string{"id", "tenant_id", "name", "description", "version"}
	tenantColumn    = "tenant_id"
)

//go:generate mockery --name=EntityConverter --output=automock --outpkg=automock --case=underscore
type EntityConverter interface {
	MultipleFromEntities(entities SolutionCollection) []*model.Solution
}

type pgRepository struct {
	existQuerier       repo.ExistQuerier
	singleGetter       repo.SingleGetter
	singleGetterGlobal repo.SingleGetterGlobal
	deleter            repo.Deleter
	pageableQuerier    repo.PageableQuerier
	lister             repo.Lister
	creator            repo.Creator
	updater            repo.Updater
	conv               EntityConverter
}

func NewRepository(conv EntityConverter) *pgRepository {
	return &pgRepository{
		existQuerier:       repo.NewExistQuerier(resource.Solution, solutionTable, tenantColumn),
		singleGetter:       repo.NewSingleGetter(resource.Solution, solutionTable, tenantColumn, solutionColumns),
		singleGetterGlobal: repo.NewSingleGetterGlobal(resource.Solution, solutionTable, solutionColumns),
		deleter:            repo.NewDeleter(resource.Solution, solutionTable, tenantColumn),
		pageableQuerier:    repo.NewPageableQuerier(resource.Solution, solutionTable, tenantColumn, solutionColumns),
		lister:             repo.NewLister(resource.Solution, solutionTable, tenantColumn, solutionColumns),
		creator:            repo.NewCreator(resource.Solution, solutionTable, solutionColumns),
		updater:            repo.NewUpdater(resource.Solution, solutionTable, []string{"name", "description", "status_condition", "status_timestamp"}, tenantColumn, []string{"id"}),
		conv:               conv,
	}
}

func (r *pgRepository) Exists(ctx context.Context, tenant, id string) (bool, error) {
	return r.existQuerier.Exists(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)})
}

func (r *pgRepository) Delete(ctx context.Context, tenant string, id string) error {
	return r.deleter.DeleteOne(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)})
}

func (r *pgRepository) GetByID(ctx context.Context, tenant, id string) (*model.Solution, error) {
	var solutionEnt Solution
	if err := r.singleGetter.Get(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)}, repo.NoOrderBy, &solutionEnt); err != nil {
		return nil, err
	}

	return solutionEnt.ToModel(), nil
}

func (r *pgRepository) GetByFiltersAndID(ctx context.Context, tenant, id string, filter []*labelfilter.LabelFilter) (*model.Solution, error) {
	tenantID, err := uuid.Parse(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing tenant as UUID")
	}

	additionalConditions := repo.Conditions{repo.NewEqualCondition("id", id)}

	filterSubquery, args, err := label.FilterQuery(model.SolutionLabelableObject, label.IntersectSet, tenantID, filter)
	if err != nil {
		return nil, errors.Wrap(err, "while building filter query")
	}
	if filterSubquery != "" {
		additionalConditions = append(additionalConditions, repo.NewInConditionForSubQuery("id", filterSubquery, args))
	}

	var solutionEnt Solution
	if err := r.singleGetter.Get(ctx, tenant, additionalConditions, repo.NoOrderBy, &solutionEnt); err != nil {
		return nil, err
	}

	return solutionEnt.ToModel(), nil
}

func (r *pgRepository) GetByFiltersGlobal(ctx context.Context, filter []*labelfilter.LabelFilter) (*model.Solution, error) {
	filterSubquery, args, err := label.FilterQueryGlobal(model.SolutionLabelableObject, label.IntersectSet, filter)
	if err != nil {
		return nil, errors.Wrap(err, "while building filter query")
	}

	var additionalConditions repo.Conditions
	if filterSubquery != "" {
		additionalConditions = append(additionalConditions, repo.NewInConditionForSubQuery("id", filterSubquery, args))
	}

	var solutionEnt Solution
	if err := r.singleGetterGlobal.GetGlobal(ctx, additionalConditions, repo.NoOrderBy, &solutionEnt); err != nil {
		return nil, err
	}

	return solutionEnt.ToModel(), nil
}

func (r *pgRepository) ListAll(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter) ([]*model.Solution, error) {
	var entities SolutionCollection

	tenantID, err := uuid.Parse(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing tenant as UUID")
	}

	filterSubquery, args, err := label.FilterQuery(model.SolutionLabelableObject, label.IntersectSet, tenantID, filter)
	if err != nil {
		return nil, errors.Wrap(err, "while building filter query")
	}

	var conditions repo.Conditions
	if filterSubquery != "" {
		conditions = append(conditions, repo.NewInConditionForSubQuery("id", filterSubquery, args))
	}

	err = r.lister.List(ctx, tenant, &entities, conditions...)
	if err != nil {
		return nil, err
	}

	return r.conv.MultipleFromEntities(entities), nil
}

type SolutionCollection []Solution

func (s SolutionCollection) Len() int {
	return len(s)
}

func (r *pgRepository) List(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.SolutionPage, error) {
	var solutionsCollection SolutionCollection
	tenantID, err := uuid.Parse(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing tenant as UUID")
	}
	filterSubquery, args, err := label.FilterQuery(model.SolutionLabelableObject, label.IntersectSet, tenantID, filter)
	if err != nil {
		return nil, errors.Wrap(err, "while building filter query")
	}

	var conditions repo.Conditions
	if filterSubquery != "" {
		conditions = append(conditions, repo.NewInConditionForSubQuery("id", filterSubquery, args))
	}

	page, totalCount, err := r.pageableQuerier.List(ctx, tenant, pageSize, cursor, "name", &solutionsCollection, conditions...)

	if err != nil {
		return nil, err
	}

	items := r.conv.MultipleFromEntities(solutionsCollection)

	return &model.SolutionPage{
		Data:       items,
		TotalCount: totalCount,
		PageInfo:   page}, nil
}

func (r *pgRepository) Create(ctx context.Context, item *model.Solution) error {
	if item == nil {
		return apperrors.NewInternalError("item can not be empty")
	}

	solutionEnt, err := EntityFromSolutionModel(item)
	if err != nil {
		return errors.Wrap(err, "while creating solution entity from model")
	}

	return r.creator.Create(ctx, solutionEnt)
}

func (r *pgRepository) Update(ctx context.Context, item *model.Solution) error {
	solutionEnt, err := EntityFromSolutionModel(item)
	if err != nil {
		return errors.Wrap(err, "while creating solution entity from model")
	}
	return r.updater.UpdateSingle(ctx, solutionEnt)
}

func (r *pgRepository) GetOldestForFilters(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter) (*model.Solution, error) {
	tenantID, err := uuid.Parse(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing tenant as UUID")
	}

	var additionalConditions repo.Conditions
	filterSubquery, args, err := label.FilterQuery(model.SolutionLabelableObject, label.IntersectSet, tenantID, filter)
	if err != nil {
		return nil, errors.Wrap(err, "while building filter query")
	}
	if filterSubquery != "" {
		additionalConditions = append(additionalConditions, repo.NewInConditionForSubQuery("id", filterSubquery, args))
	}

	orderByParams := repo.OrderByParams{repo.NewAscOrderBy("creation_timestamp")}

	var solutionEnt Solution
	if err := r.singleGetter.Get(ctx, tenant, additionalConditions, orderByParams, &solutionEnt); err != nil {
		return nil, err
	}

	return solutionEnt.ToModel(), nil
}
