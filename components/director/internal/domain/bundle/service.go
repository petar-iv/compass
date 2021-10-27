package bundle

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/pkg/errors"
)

// BundleRepository missing godoc
//go:generate mockery --name=BundleRepository --output=automock --outpkg=automock --case=underscore
type BundleRepository interface {
	Create(ctx context.Context, item *model.Bundle) error
	Update(ctx context.Context, item *model.Bundle) error
	Delete(ctx context.Context, tenant, id string) error
	Exists(ctx context.Context, tenant, id string) (bool, error)
	GetByID(ctx context.Context, tenant, id string) (*model.Bundle, error)
	GetForApplication(ctx context.Context, tenant string, id string, applicationID string) (*model.Bundle, error)
	ListByApplicationIDNoPaging(ctx context.Context, tenantID, appID string) ([]*model.Bundle, error)
	ListByApplicationIDs(ctx context.Context, tenantID string, applicationIDs []string, pageSize int, cursor string) ([]*model.BundlePage, error)
	ListByApplicationIDsForScenariosNoPaging(ctx context.Context, tenantID string, applicationIDs []string, scenarios []string) ([]*model.Bundle, error)
	ListByApplicationIDsForScenarios(ctx context.Context, tenantID string, applicationIDs []string, scenarios []string, pageSize int, cursor string) ([]*model.BundlePage, error)
}

// LabelRepository is used to manage the storage of labels.
//go:generate mockery --name=LabelRepository --output=automock --outpkg=automock --case=underscore
type LabelRepository interface {
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
	ListForObject(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error)
	Delete(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) error
}

// LabelUpsertService is used to upsert labels for different entities.
//go:generate mockery --name=LabelUpsertService --output=automock --outpkg=automock --case=underscore
type LabelUpsertService interface {
	UpsertLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) error
}

// UIDService missing godoc
//go:generate mockery --name=UIDService --output=automock --outpkg=automock --case=underscore
type UIDService interface {
	Generate() string
}

type service struct {
	bndlRepo           BundleRepository
	apiSvc             APIService
	eventSvc           EventService
	documentSvc        DocumentService
	labelRepo          LabelRepository
	labelUpsertService LabelUpsertService

	uidService UIDService
}

// NewService missing godoc
func NewService(bndlRepo BundleRepository, apiSvc APIService, eventSvc EventService, documentSvc DocumentService, labelRepo LabelRepository, labelUpsertService LabelUpsertService, uidService UIDService) *service {
	return &service{
		bndlRepo:           bndlRepo,
		apiSvc:             apiSvc,
		eventSvc:           eventSvc,
		documentSvc:        documentSvc,
		labelRepo:          labelRepo,
		labelUpsertService: labelUpsertService,
		uidService:         uidService,
	}
}

// Create missing godoc
func (s *service) Create(ctx context.Context, applicationID string, in model.BundleCreateInput) (string, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return "", err
	}

	id := s.uidService.Generate()
	bndl := in.ToBundle(id, applicationID, tnt)

	err = s.bndlRepo.Create(ctx, bndl)
	if err != nil {
		return "", errors.Wrapf(err, "error occurred while creating a Bundle with id %s and name %s for Application with id %s", id, bndl.Name, applicationID)
	}
	log.C(ctx).Infof("Successfully created a Bundle with id %s and name %s for Application with id %s", id, bndl.Name, applicationID)

	log.C(ctx).Infof("Creating related resources in Bundle with id %s and name %s for Application with id %s", id, bndl.Name, applicationID)
	err = s.createRelatedResources(ctx, in, id, applicationID)
	if err != nil {
		return "", errors.Wrapf(err, "while creating related resources for Application with id %s", applicationID)
	}

	return id, nil
}

// CreateMultiple missing godoc
func (s *service) CreateMultiple(ctx context.Context, applicationID string, in []*model.BundleCreateInput) error {
	if in == nil {
		return nil
	}

	for _, bndl := range in {
		if bndl == nil {
			continue
		}

		_, err := s.Create(ctx, applicationID, *bndl)
		if err != nil {
			return errors.Wrapf(err, "while creating Bundle for Application with id %s", applicationID)
		}
	}

	return nil
}

// Update missing godoc
func (s *service) Update(ctx context.Context, id string, in model.BundleUpdateInput) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	bndl, err := s.bndlRepo.GetByID(ctx, tnt, id)
	if err != nil {
		return errors.Wrapf(err, "while getting Bundle with id %s", id)
	}

	bndl.SetFromUpdateInput(in)

	err = s.bndlRepo.Update(ctx, bndl)
	if err != nil {
		return errors.Wrapf(err, "while updating Bundle with id %s", id)
	}
	return nil
}

// Delete missing godoc
func (s *service) Delete(ctx context.Context, id string) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "while loading tenant from context")
	}

	err = s.bndlRepo.Delete(ctx, tnt, id)
	if err != nil {
		return errors.Wrapf(err, "while deleting Bundle with id %s", id)
	}

	return nil
}

// GetLabel retrieves a label by its key for a given Bundle ID.
func (s *service) GetLabel(ctx context.Context, bundleID string, key string) (*model.Label, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	exists, err := s.bndlRepo.Exists(ctx, tnt, bundleID)
	if err != nil {
		return nil, errors.Wrap(err, "while checking Bundle existence")
	}
	if !exists {
		return nil, fmt.Errorf("bundle with ID %s doesn't exist", bundleID)
	}

	label, err := s.labelRepo.GetByKey(ctx, tnt, model.BundleLabelableObject, bundleID, key)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Bundle")
	}

	return label, nil
}

// ListLabels retrieves all labels for a given Bundle.
func (s *service) ListLabels(ctx context.Context, bundleID string) (map[string]*model.Label, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	exists, err := s.bndlRepo.Exists(ctx, tnt, bundleID)
	if err != nil {
		return nil, errors.Wrap(err, "while checking Bundle existence")
	}
	if !exists {
		return nil, fmt.Errorf("bundle with ID %s doesn't exist", bundleID)
	}

	labels, err := s.labelRepo.ListForObject(ctx, tnt, model.BundleLabelableObject, bundleID)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Bundle")
	}

	return labels, nil
}

// SetLabel labels a given bundle with a key-value pair.
func (s *service) SetLabel(ctx context.Context, in *model.LabelInput) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	bndlExists, err := s.bndlRepo.Exists(ctx, tnt, in.ObjectID)
	if err != nil {
		return errors.Wrap(err, "while checking Bundle existence")
	}
	if !bndlExists {
		return apperrors.NewNotFoundError(resource.Bundle, in.ObjectID)
	}

	err = s.labelUpsertService.UpsertLabel(ctx, tnt, in)
	if err != nil {
		return errors.Wrapf(err, "while creating label for Bundle")
	}

	return nil
}

// DeleteLabel deletes a label from a bundle for a given key.
func (s *service) DeleteLabel(ctx context.Context, bundleID string, key string) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	exists, err := s.bndlRepo.Exists(ctx, tnt, bundleID)
	if err != nil {
		return errors.Wrap(err, "while checking Bundle existence")
	}
	if !exists {
		return fmt.Errorf("bundle with ID %s doesn't exist", bundleID)
	}

	err = s.labelRepo.Delete(ctx, tnt, model.BundleLabelableObject, bundleID, key)
	if err != nil {
		return errors.Wrapf(err, "while deleting Bundle label")
	}

	return nil
}

// Exist missing godoc
func (s *service) Exist(ctx context.Context, id string) (bool, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return false, errors.Wrap(err, "while loading tenant from context")
	}

	exist, err := s.bndlRepo.Exists(ctx, tnt, id)
	if err != nil {
		return false, errors.Wrapf(err, "while getting Bundle with ID: [%s]", id)
	}

	return exist, nil
}

// Get missing godoc
func (s *service) Get(ctx context.Context, id string) (*model.Bundle, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "while loading tenant from context")
	}

	bndl, err := s.bndlRepo.GetByID(ctx, tnt, id)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting Bundle with ID: [%s]", id)
	}

	return bndl, nil
}

// GetForApplication missing godoc
func (s *service) GetForApplication(ctx context.Context, id string, applicationID string) (*model.Bundle, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	bndl, err := s.bndlRepo.GetForApplication(ctx, tnt, id, applicationID)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting Bundle with ID: [%s]", id)
	}

	return bndl, nil
}

// ListByApplicationIDNoPaging missing godoc
func (s *service) ListByApplicationIDNoPaging(ctx context.Context, appID string) ([]*model.Bundle, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.bndlRepo.ListByApplicationIDNoPaging(ctx, tnt, appID)
}

// ListByApplicationIDs missing godoc
func (s *service) ListByApplicationIDs(ctx context.Context, applicationIDs []string, pageSize int, cursor string) ([]*model.BundlePage, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.bndlRepo.ListByApplicationIDs(ctx, tnt, applicationIDs, pageSize, cursor)
}

// ListByApplicationIDsForScenariosNoPaging retrieves all bundles for a given number of application IDs that are part of given scenarios.
func (s *service) ListByApplicationIDsForScenariosNoPaging(ctx context.Context, applicationIDs []string, scenarios []string) ([]*model.Bundle, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.bndlRepo.ListByApplicationIDsForScenariosNoPaging(ctx, tnt, applicationIDs, scenarios)
}

// ListByApplicationIDsForScenarios retrieves all bundles for a given number of application IDs that are part of given scenarios using paging.
func (s *service) ListByApplicationIDsForScenarios(ctx context.Context, applicationIDs []string, scenarios []string, pageSize int, cursor string) ([]*model.BundlePage, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.bndlRepo.ListByApplicationIDsForScenarios(ctx, tnt, applicationIDs, scenarios, pageSize, cursor)
}

func (s *service) createRelatedResources(ctx context.Context, in model.BundleCreateInput, bundleID, appID string) error {
	for i := range in.APIDefinitions {
		_, err := s.apiSvc.CreateInBundle(ctx, appID, bundleID, *in.APIDefinitions[i], in.APISpecs[i])
		if err != nil {
			return errors.Wrapf(err, "while creating APIs for bundle with id %q", bundleID)
		}
	}

	for i := range in.EventDefinitions {
		_, err := s.eventSvc.CreateInBundle(ctx, appID, bundleID, *in.EventDefinitions[i], in.EventSpecs[i])
		if err != nil {
			return errors.Wrapf(err, "while creating Event for bundle with id %q", bundleID)
		}
	}

	for _, document := range in.Documents {
		_, err := s.documentSvc.CreateInBundle(ctx, bundleID, *document)
		if err != nil {
			return errors.Wrapf(err, "while creating Document for bundle with id %q", bundleID)
		}
	}

	return nil
}
