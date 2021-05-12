package solution

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/pkg/errors"
)

//go:generate mockery --name=SolutionRepository --output=automock --outpkg=automock --case=underscore
type SolutionRepository interface {
	Exists(ctx context.Context, tenant, id string) (bool, error)
	GetByID(ctx context.Context, tenant, id string) (*model.Solution, error)
	GetByFiltersGlobal(ctx context.Context, filter []*labelfilter.LabelFilter) (*model.Solution, error)
	List(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.SolutionPage, error)
	ListAll(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter) ([]*model.Solution, error)
	Create(ctx context.Context, item *model.Solution) error
	Update(ctx context.Context, item *model.Solution) error
	Delete(ctx context.Context, tenant, id string) error
}

//go:generate mockery --name=LabelRepository --output=automock --outpkg=automock --case=underscore
type LabelRepository interface {
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
	ListForObject(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error)
	Delete(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) error
	DeleteAll(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) error
	DeleteByKeyNegationPattern(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, labelKeyPattern string) error
}

//go:generate mockery --name=LabelUpsertService --output=automock --outpkg=automock --case=underscore
type LabelUpsertService interface {
	UpsertMultipleLabels(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, labels map[string]interface{}) error
	UpsertLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) error
}

//go:generate mockery --name=ScenariosService --output=automock --outpkg=automock --case=underscore
type ScenariosService interface {
	AddScenarios(ctx context.Context, scenarios ...string) error
	EnsureScenariosLabelDefinitionExists(ctx context.Context, tenant string) error
}

//go:generate mockery --name=UIDService --output=automock --outpkg=automock --case=underscore
type UIDService interface {
	Generate() string
}

type service struct {
	repo      SolutionRepository
	labelRepo LabelRepository

	labelUpsertService LabelUpsertService
	uidService         UIDService
	scenariosService   ScenariosService

	protectedLabelPattern string
}

func NewService(repo SolutionRepository,
	labelRepo LabelRepository,
	scenariosService ScenariosService,
	labelUpsertService LabelUpsertService,
	uidService UIDService,
	protectedLabelPattern string) *service {
	return &service{
		repo:                  repo,
		labelRepo:             labelRepo,
		scenariosService:      scenariosService,
		labelUpsertService:    labelUpsertService,
		uidService:            uidService,
		protectedLabelPattern: protectedLabelPattern,
	}
}

func (s *service) List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.SolutionPage, error) {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.repo.List(ctx, tenant, filter, pageSize, cursor)
}

func (s *service) Get(ctx context.Context, id string) (*model.Solution, error) {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	solution, err := s.repo.GetByID(ctx, tenant, id)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting Solution with ID %s", id)
	}

	return solution, nil
}

func (s *service) Exist(ctx context.Context, id string) (bool, error) {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return false, errors.Wrapf(err, "while loading tenant from context")
	}

	exist, err := s.repo.Exists(ctx, tenant, id)
	if err != nil {
		return false, errors.Wrapf(err, "while getting Solution with ID %s", id)
	}

	return exist, nil
}

func (s *service) Create(ctx context.Context, in model.SolutionInput) (string, error) {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "while loading tenant from context")
	}

	id := s.uidService.Generate()
	solution := in.ToSolution(id, tenant)

	if err = s.repo.Create(ctx, solution); err != nil {
		return "", errors.Wrapf(err, "while creating Solution")
	}

	if err = s.scenariosService.EnsureScenariosLabelDefinitionExists(ctx, tenant); err != nil {
		return "", errors.Wrapf(err, "while ensuring Label Definition with key %s exists", model.ScenariosKey)
	}

	if in.Labels == nil {
		if in.Labels == nil {
			in.Labels = make(map[string]interface{}, 1)
		}
	}

	err = s.labelUpsertService.UpsertMultipleLabels(ctx, tenant, model.SolutionLabelableObject, id, in.Labels)
	if err != nil {
		return id, errors.Wrapf(err, "while creating multiple labels for Solution")
	}

	return id, nil
}

func (s *service) Update(ctx context.Context, id string, in model.SolutionInput) error {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	rtm, err := s.repo.GetByID(ctx, tenant, id)
	if err != nil {
		return errors.Wrapf(err, "while getting Solution with id %s", id)
	}

	rtm = in.ToSolution(id, rtm.Tenant)

	err = s.repo.Update(ctx, rtm)
	if err != nil {
		return errors.Wrap(err, "while updating Solution")
	}

	if in.Labels == nil {
		if in.Labels == nil {
			in.Labels = make(map[string]interface{}, 1)
		}
	}

	err = s.labelRepo.DeleteByKeyNegationPattern(ctx, tenant, model.SolutionLabelableObject, id, s.protectedLabelPattern)
	if err != nil {
		return errors.Wrapf(err, "while deleting all labels for Solution")
	}

	if in.Labels == nil {
		return nil
	}

	err = s.labelUpsertService.UpsertMultipleLabels(ctx, tenant, model.SolutionLabelableObject, id, in.Labels)
	if err != nil {
		return errors.Wrapf(err, "while creating multiple labels for Solution")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	err = s.repo.Delete(ctx, tenant, id)
	if err != nil {
		return errors.Wrapf(err, "while deleting Solution")
	}

	return nil
}

func (s *service) SetLabel(ctx context.Context, labelInput *model.LabelInput) error {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	if err := s.ensureSolutionExists(ctx, tenant, labelInput.ObjectID); err != nil {
		return err
	}

	err = s.labelUpsertService.UpsertLabel(ctx, tenant, labelInput)
	if err != nil {
		return errors.Wrapf(err, "while creating label for Solution")
	}

	return nil
}

func (s *service) GetLabel(ctx context.Context, solutionID string, key string) (*model.Label, error) {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	if err := s.ensureSolutionExists(ctx, tenant, solutionID); err != nil {
		return nil, err
	}

	label, err := s.labelRepo.GetByKey(ctx, tenant, model.SolutionLabelableObject, solutionID, key)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Solution")
	}

	return label, nil
}

func (s *service) ListLabels(ctx context.Context, solutionID string) (map[string]*model.Label, error) {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	if err := s.ensureSolutionExists(ctx, tenant, solutionID); err != nil {
		return nil, err
	}

	labels, err := s.labelRepo.ListForObject(ctx, tenant, model.SolutionLabelableObject, solutionID)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Solution")
	}

	return labels, nil
}

func (s *service) DeleteLabel(ctx context.Context, solutionID string, key string) error {
	tenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	if err := s.ensureSolutionExists(ctx, tenant, solutionID); err != nil {
		return err
	}

	err = s.labelRepo.Delete(ctx, tenant, model.SolutionLabelableObject, solutionID, key)
	if err != nil {
		return errors.Wrapf(err, "while deleting Solution label")
	}

	return nil
}
func (s *service) ensureSolutionExists(ctx context.Context, tnt string, solutionID string) error {
	exists, err := s.repo.Exists(ctx, tnt, solutionID)
	if err != nil {
		return errors.Wrap(err, "while checking Solution existence")
	}
	if !exists {
		return fmt.Errorf("Solution with ID %s doesn't exist", solutionID)
	}

	return nil
}
