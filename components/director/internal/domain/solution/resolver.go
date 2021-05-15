package solution

import (
	"context"
	"strings"

	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/inputvalidation"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"
	"k8s.io/apiserver/pkg/storage/names"
)

//go:generate mockery --name=SolutionService --output=automock --outpkg=automock --case=underscore
type SolutionService interface {
	Create(ctx context.Context, in model.SolutionInput) (string, error)
	Update(ctx context.Context, id string, in model.SolutionInput) error
	Get(ctx context.Context, id string) (*model.Solution, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.SolutionPage, error)
	SetLabel(ctx context.Context, label *model.LabelInput) error
	GetLabel(ctx context.Context, solutionID string, key string) (*model.Label, error)
	ListLabels(ctx context.Context, solutionID string) (map[string]*model.Label, error)
	DeleteLabel(ctx context.Context, solutionID string, key string) error
}

//go:generate mockery --name=ScenarioAssignmentService --output=automock --outpkg=automock --case=underscore
type ScenarioAssignmentService interface {
	GetForScenarioName(ctx context.Context, scenarioName string) (model.AutomaticScenarioAssignment, error)
	Delete(ctx context.Context, in model.AutomaticScenarioAssignment) error
}

//go:generate mockery --name=SolutionConverter --output=automock --outpkg=automock --case=underscore
type SolutionConverter interface {
	ToGraphQL(in *model.Solution) *graphql.Solution
	MultipleToGraphQL(in []*model.Solution) []*graphql.Solution
	InputFromGraphQL(in graphql.SolutionInput) (model.SolutionInput, error)
}

type ApplicationService interface {
	CreateFromTemplate(ctx context.Context, in model.ApplicationRegisterInput, appTemplateID *string) (string, error)
	SetLabel(ctx context.Context, labelInput *model.LabelInput) error
}

type ApplicationConverter interface {
	ToGraphQL(in *model.Application) *graphql.Application
	CreateInputJSONToGQL(in string) (graphql.ApplicationRegisterInput, error)
	CreateInputFromGraphQL(ctx context.Context, in graphql.ApplicationRegisterInput) (model.ApplicationRegisterInput, error)
}

type ApplicationTemplateService interface {
	GetByName(ctx context.Context, name string) (*model.ApplicationTemplate, error)
	PrepareApplicationCreateInputJSON(appTemplate *model.ApplicationTemplate, values model.ApplicationFromTemplateInputValues) (string, error)
}

type ApplicationTemplateConverter interface {
	ApplicationFromTemplateInputFromGraphQL(in graphql.ApplicationFromTemplateInput) model.ApplicationFromTemplateInput
}

type LabelDefinitionService interface {
	Get(ctx context.Context, tenant string, key string) (*model.LabelDefinition, error)
	Create(ctx context.Context, def model.LabelDefinition) (model.LabelDefinition, error)
	Update(ctx context.Context, def model.LabelDefinition) error
}

type Resolver struct {
	transact         persistence.Transactioner
	solutionSvc      SolutionService
	scenariosService ScenariosService
	converter        SolutionConverter

	appSvc          ApplicationService
	appConv         ApplicationConverter
	appTemplateSvc  ApplicationTemplateService
	appTemplateConv ApplicationTemplateConverter
	labelDefSvc     LabelDefinitionService

	bndlConverter BundleConverter
}

func NewResolver(transact persistence.Transactioner, solutionService SolutionService, scenariosService ScenariosService, conv SolutionConverter,
	appSvc ApplicationService, appConv ApplicationConverter, appTemplateSvc ApplicationTemplateService, appTemplateConv ApplicationTemplateConverter,
	bndlConv BundleConverter, labelDefSvc LabelDefinitionService) *Resolver {
	return &Resolver{
		transact:         transact,
		solutionSvc:      solutionService,
		scenariosService: scenariosService,
		converter:        conv,
		appSvc:           appSvc,
		appConv:          appConv,
		appTemplateSvc:   appTemplateSvc,
		appTemplateConv:  appTemplateConv,
		labelDefSvc:      labelDefSvc,
		bndlConverter:    bndlConv,
	}
}

func (r *Resolver) Solutions(ctx context.Context, filter []*graphql.LabelFilter, first *int, after *graphql.PageCursor) (*graphql.SolutionPage, error) {
	labelFilter := labelfilter.MultipleFromGraphQL(filter)

	var cursor string
	if after != nil {
		cursor = string(*after)
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	if first == nil {
		return nil, apperrors.NewInvalidDataError("missing required parameter 'first'")
	}

	solutionsPage, err := r.solutionSvc.List(ctx, labelFilter, *first, cursor)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlSolutions := r.converter.MultipleToGraphQL(solutionsPage.Data)

	return &graphql.SolutionPage{
		Data:       gqlSolutions,
		TotalCount: solutionsPage.TotalCount,
		PageInfo: &graphql.PageInfo{
			StartCursor: graphql.PageCursor(solutionsPage.PageInfo.StartCursor),
			EndCursor:   graphql.PageCursor(solutionsPage.PageInfo.EndCursor),
			HasNextPage: solutionsPage.PageInfo.HasNextPage,
		},
	}, nil
}

func (r *Resolver) Solution(ctx context.Context, id string) (*graphql.Solution, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	solution, err := r.solutionSvc.Get(ctx, id)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.converter.ToGraphQL(solution), nil
}

func (r *Resolver) RegisterSolution(ctx context.Context, in graphql.SolutionInput) (*graphql.Solution, error) {
	convertedIn, err := r.converter.InputFromGraphQL(in)
	if err != nil {
		return nil, err
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	id, err := r.solutionSvc.Create(ctx, convertedIn)
	if err != nil {
		return nil, err
	}
	solution, err := r.solutionSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.scenariosService.AddScenarios(ctx, in.Name); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to add new scenario for solution with name %s: %v", in.Name, err)
		return nil, apperrors.NewInternalError("Failed to create scenario")
	}

	err = r.solutionSvc.SetLabel(ctx, &model.LabelInput{
		Key:        model.ScenariosKey,
		Value:      []string{in.Name},
		ObjectType: model.SolutionLabelableObject,
		ObjectID:   id,
	})
	if err != nil {
		return nil, err
	}

	if err := r.createRelatedResources(ctx, in); err != nil {
		log.C(ctx).Errorf("Failed to create solution dependency: %v", err)
		return nil, apperrors.NewInternalError("failed to create dependencies")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.converter.ToGraphQL(solution), nil
}

func (r *Resolver) createRelatedResources(ctx context.Context, in graphql.SolutionInput) error {
	for _, dep := range in.Dependencies {
		if dep == nil {
			continue
		}

		if err := dep.Application.Validate(); err != nil {
			return err
		}

		bundlesIn, err := r.bndlConverter.MultipleCreateInputFromGraphQL(dep.Bundles)
		if err != nil {
			return err
		}

		id, err := r.registerAppFromTemplate(ctx, *dep.Application, bundlesIn)
		if err != nil {
			return errors.Wrap(err, "failed to create solutiuon dependency of type application")
		}

		if err := r.appSvc.SetLabel(ctx, &model.LabelInput{
			Key:        model.ScenariosKey,
			Value:      []string{in.Name},
			ObjectType: model.ApplicationLabelableObject,
			ObjectID:   id,
		}); err != nil {
			return errors.Wrapf(err, "failed to assign scenario to application with ID %s", id)
		}
	}
	return nil
}

func (r *Resolver) registerAppFromTemplate(ctx context.Context, in graphql.ApplicationFromTemplateInput, bundlesIn []*model.BundleCreateInput) (string, error) {
	log.C(ctx).Infof("Registering an Application from Application Template with name %s", in.TemplateName)
	convertedIn := r.appTemplateConv.ApplicationFromTemplateInputFromGraphQL(in)

	log.C(ctx).Debugf("Extracting Application Template with name %s from GraphQL input", in.TemplateName)
	appTemplate, err := r.appTemplateSvc.GetByName(ctx, convertedIn.TemplateName)
	if err != nil {
		return "", err
	}
	baseName := strings.ToLower(strings.ReplaceAll(appTemplate.Name, " ", "-"))
	applicationName := names.SimpleNameGenerator.GenerateName(baseName)
	log.C(ctx).Infof("Registering an Application with name %s from Application Template with name %s", applicationName, in.TemplateName)

	// TODO think think think
	convertedIn.Values = append(convertedIn.Values, &model.ApplicationTemplateValueInput{
		Placeholder: "name",
		Value:       applicationName,
	})
	log.C(ctx).Debugf("Preparing ApplicationCreateInput JSON from Application Template with name %s", in.TemplateName)
	appCreateInputJSON, err := r.appTemplateSvc.PrepareApplicationCreateInputJSON(appTemplate, convertedIn.Values)
	if err != nil {
		return "", errors.Wrapf(err, "while preparing ApplicationCreateInput JSON from Application Template with name %s", in.TemplateName)
	}

	log.C(ctx).Debugf("Converting ApplicationCreateInput JSON to GraphQL ApplicationRegistrationInput from Application Template with name %s", in.TemplateName)
	appCreateInputGQL, err := r.appConv.CreateInputJSONToGQL(appCreateInputJSON)
	if err != nil {
		return "", errors.Wrapf(err, "while converting ApplicationCreateInput JSON to GraphQL ApplicationRegistrationInput from Application Template with name %s", in.TemplateName)
	}

	log.C(ctx).Infof("Validating GraphQL ApplicationRegistrationInput from Application Template with name %s", convertedIn.TemplateName)
	if err := inputvalidation.Validate(appCreateInputGQL); err != nil {
		return "", errors.Wrapf(err, "while validating application input from Application Template with name %s", convertedIn.TemplateName)
	}

	appCreateInputModel, err := r.appConv.CreateInputFromGraphQL(ctx, appCreateInputGQL)

	appCreateInputModel.Bundles = bundlesIn
	if err != nil {
		return "", errors.Wrap(err, "while converting ApplicationFromTemplate input")
	}

	log.C(ctx).Infof("Creating an Application with name %s from Application Template with name %s", applicationName, in.TemplateName)
	id, err := r.appSvc.CreateFromTemplate(ctx, appCreateInputModel, &appTemplate.ID)
	if err != nil {
		return "", errors.Wrapf(err, "while creating an Application with name %s from Application Template with name %s", applicationName, in.TemplateName)
	}
	log.C(ctx).Infof("Application with name %s and id %s successfully created from Application Template with name %s", applicationName, id, in.TemplateName)
	return id, nil
}

func (r *Resolver) UpdateSolution(ctx context.Context, id string, in graphql.SolutionInput) (*graphql.Solution, error) {
	convertedIn, err := r.converter.InputFromGraphQL(in)
	if err != nil {
		return nil, err
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	err = r.solutionSvc.Update(ctx, id, convertedIn)
	if err != nil {
		return nil, err
	}

	solution, err := r.solutionSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.converter.ToGraphQL(solution), nil
}

func (r *Resolver) DeleteSolution(ctx context.Context, id string) (*graphql.Solution, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	solution, err := r.solutionSvc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	deletedSolution := r.converter.ToGraphQL(solution)

	err = r.solutionSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return deletedSolution, nil
}

func (r *Resolver) GetLabel(ctx context.Context, solutionID string, key string) (*graphql.Labels, error) {
	if solutionID == "" {
		return nil, apperrors.NewInternalError("Solution cannot be empty")
	}
	if key == "" {
		return nil, apperrors.NewInternalError("Solution label key cannot be empty")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	label, err := r.solutionSvc.GetLabel(ctx, solutionID, key)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	resultLabels := make(map[string]interface{})
	resultLabels[key] = label.Value

	var gqlLabels graphql.Labels = resultLabels
	return &gqlLabels, nil
}

func (r *Resolver) SetSolutionLabel(ctx context.Context, solutionID string, key string, value interface{}) (*graphql.Label, error) {
	// TODO: Use @validation directive on input type instead, after resolving https://github.com/kyma-incubator/compass/issues/515
	gqlLabel := graphql.LabelInput{Key: key, Value: value}
	if err := inputvalidation.Validate(&gqlLabel); err != nil {
		return nil, errors.Wrap(err, "validation error for type LabelInput")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	if err := r.solutionSvc.SetLabel(ctx, &model.LabelInput{
		Key:        key,
		Value:      value,
		ObjectType: model.SolutionLabelableObject,
		ObjectID:   solutionID,
	}); err != nil {
		return nil, err
	}

	label, err := r.solutionSvc.GetLabel(ctx, solutionID, key)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting label with key: [%s]", key)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &graphql.Label{
		Key:   label.Key,
		Value: label.Value,
	}, nil
}

func (r *Resolver) DeleteSolutionLabel(ctx context.Context, solutionID string, key string) (*graphql.Label, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	label, err := r.solutionSvc.GetLabel(ctx, solutionID, key)
	if err != nil {
		return nil, err
	}

	err = r.solutionSvc.DeleteLabel(ctx, solutionID, key)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &graphql.Label{
		Key:   key,
		Value: label.Value,
	}, nil
}

func (r *Resolver) Labels(ctx context.Context, obj *graphql.Solution, key *string) (graphql.Labels, error) {
	if obj == nil {
		return nil, apperrors.NewInternalError("Solution cannot be empty")
	}

	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	itemMap, err := r.solutionSvc.ListLabels(ctx, obj.ID)
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") { // TODO: Use custom error and check its type
			return nil, tx.Commit()
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	resultLabels := make(map[string]interface{})

	for _, label := range itemMap {
		resultLabels[label.Key] = label.Value
	}

	var gqlLabels graphql.Labels = resultLabels
	return gqlLabels, nil
}
