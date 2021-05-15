package bundleinstanceauth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
)

//go:generate mockery --name=Service --output=automock --outpkg=automock --case=underscore
type Service interface {
	RequestDeletion(ctx context.Context, instanceAuth *model.BundleInstanceAuth, defaultBundleInstanceAuth *model.Auth) (bool, error)
	Create(ctx context.Context, bundleID string, in model.BundleInstanceAuthRequestInput, defaultAuth *model.Auth, requestInputSchema *string) (string, error)
	Get(ctx context.Context, id string) (*model.BundleInstanceAuth, error)
	SetAuth(ctx context.Context, id string, in model.BundleInstanceAuthSetInput) error
	Delete(ctx context.Context, id string) error
}

//go:generate mockery --name=Converter --output=automock --outpkg=automock --case=underscore
type Converter interface {
	ToGraphQL(in *model.BundleInstanceAuth) (*graphql.BundleInstanceAuth, error)
	RequestInputFromGraphQL(in graphql.BundleInstanceAuthRequestInput) model.BundleInstanceAuthRequestInput
	SetInputFromGraphQL(in graphql.BundleInstanceAuthSetInput) (model.BundleInstanceAuthSetInput, error)
}

//go:generate mockery --name=BundleService --output=automock --outpkg=automock --case=underscore
type BundleService interface {
	Get(ctx context.Context, id string) (*model.Bundle, error)
	GetByInstanceAuthID(ctx context.Context, instanceAuthID string) (*model.Bundle, error)
	ListByApplicationIDNoPaging(ctx context.Context, appID string) ([]*model.Bundle, error)
}

//go:generate mockery --name=BundleService --output=automock --outpkg=automock --case=underscore
type ApplicationService interface {
	ListBySolutionIDNoPaging(ctx context.Context, solutionUUID uuid.UUID) ([]*model.Application, error)
}

//go:generate mockery --name=BundleConverter --output=automock --outpkg=automock --case=underscore
type BundleConverter interface {
	ToGraphQL(in *model.Bundle) (*graphql.Bundle, error)
}

type Resolver struct {
	transact persistence.Transactioner
	svc      Service
	bndlSvc  BundleService
	conv     Converter
	bndlConv BundleConverter
	appSvc   ApplicationService
}

func NewResolver(transact persistence.Transactioner, svc Service, appSvc ApplicationService, bndlSvc BundleService, conv Converter, bndlConv BundleConverter) *Resolver {
	return &Resolver{
		transact: transact,
		svc:      svc,
		appSvc:   appSvc,
		bndlSvc:  bndlSvc,
		conv:     conv,
		bndlConv: bndlConv,
	}
}

func (r *Resolver) BundleByInstanceAuth(ctx context.Context, authID string) (*graphql.Bundle, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	bndlInstanceAuth, err := r.svc.Get(ctx, authID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, tx.Commit()
		}
		return nil, err
	}

	pkg, err := r.bndlSvc.Get(ctx, bndlInstanceAuth.BundleID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.bndlConv.ToGraphQL(pkg)
}

func (r *Resolver) BundleInstanceAuth(ctx context.Context, id string) (*graphql.BundleInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	bndlInstanceAuth, err := r.svc.Get(ctx, id)
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

	return r.conv.ToGraphQL(bndlInstanceAuth)
}

func (r *Resolver) DeleteBundleInstanceAuth(ctx context.Context, authID string) (*graphql.BundleInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}

	defer r.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	instanceAuth, err := r.svc.Get(ctx, authID)
	if err != nil {
		return nil, err
	}

	err = r.svc.Delete(ctx, authID)
	if err != nil {
		return nil, err
	}

	log.C(ctx).Infof("BundleInstanceAuth with id %s successfully deleted", authID)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.conv.ToGraphQL(instanceAuth)
}

func (r *Resolver) SetBundleInstanceAuth(ctx context.Context, authID string, in graphql.BundleInstanceAuthSetInput) (*graphql.BundleInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}

	defer r.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	log.C(ctx).Infof("Setting credentials for BundleInstanceAuth with id %s", authID)

	convertedIn, err := r.conv.SetInputFromGraphQL(in)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting BundleInstanceAuth with id %s from GraphQL", authID)
	}

	err = r.svc.SetAuth(ctx, authID, convertedIn)
	if err != nil {
		return nil, err
	}

	instanceAuth, err := r.svc.Get(ctx, authID)
	if err != nil {
		return nil, err
	}

	log.C(ctx).Infof("Credentials successfully set for BundleInstanceAuth with id %s", authID)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.conv.ToGraphQL(instanceAuth)
}

func (r *Resolver) RequestBundleInstanceAuthCreation(ctx context.Context, bundleID string, in graphql.BundleInstanceAuthRequestInput) (*graphql.BundleInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}

	defer r.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	return r.requestBundleInstanceAuth(ctx, err, bundleID, in, tx)
}

func (r *Resolver) RequestBundleInstanceAuthDeletion(ctx context.Context, authID string) (*graphql.BundleInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}

	defer r.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	log.C(ctx).Infof("Requesting BundleInstanceAuth deletion for BundleInstanceAuth with id %s", authID)

	instanceAuth, err := r.svc.Get(ctx, authID)
	if err != nil {
		return nil, err
	}

	bndl, err := r.bndlSvc.GetByInstanceAuthID(ctx, authID)
	if err != nil {
		return nil, err
	}

	deleted, err := r.svc.RequestDeletion(ctx, instanceAuth, bndl.DefaultInstanceAuth)
	if err != nil {
		return nil, err
	}

	if !deleted {
		instanceAuth, err = r.svc.Get(ctx, authID) // get InstanceAuth once again for new status
		if err != nil {
			return nil, err
		}
	}

	log.C(ctx).Infof("BundleInstanceAuth with id %s successfully deleted.", authID)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.conv.ToGraphQL(instanceAuth)
}

func (r *Resolver) RequestBundleInstanceAuthCreationForSolutionApplications(ctx context.Context, solutionID string, in []*graphql.BundleInstanceAuthRequestInputByOrdID) ([]*graphql.BundleInstanceAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}

	defer r.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	solutionUUID, err := uuid.Parse(solutionID)
	if err != nil {
		return nil, errors.Wrap(err, "while converting runtimeID to UUID")
	}

	apps, err := r.appSvc.ListBySolutionIDNoPaging(ctx, solutionUUID)
	if err != nil {
		return nil, err
	}
	if len(apps) < 1 {
		return nil, errors.New(fmt.Sprintf("No apps were found for solution with ID %s", solutionID))
	}
	result := make([]*graphql.BundleInstanceAuth, 0)
	// TODO fix mess
	for _, app := range apps {
		appID := app.ID
		bndls, err := r.bndlSvc.ListByApplicationIDNoPaging(ctx, appID)
		if err != nil {
			log.C(ctx).Errorf("Failed to list bundles for application with ID %s: %v", appID, err)
			return nil, errors.Wrap(err, fmt.Sprintf("failed to list bundles for application with ID %s", appID))
		}
		if len(bndls) == 0 {
			log.C(ctx).Errorf("No bundles found for application with ID %s", appID)
			return nil, errors.New(fmt.Sprintf("No bundles found for application with ID %s", appID))
		}

		bundlesIds := bundlesIDMapping(bndls)
		for _, b := range in {
			id, ok := bundlesIds[b.ID]
			if !ok {
				log.C(ctx).Errorf("Bundle with ORD ID %s not found in application with ID %s", b.ID)
			}
			authIn := graphql.BundleInstanceAuthRequestInput{
				ID:          str.Ptr(id),
				Context:     b.Context,
				InputParams: b.InputParams,
			}
			instanceAuth, err := r.requestBundleInstanceAuth(ctx, err, b.ID, authIn, tx)
			if err != nil {
				log.C(ctx).WithError(err).Errorf("Failed to request bundle instance auth for bundle with ID %s of application with ID %s", id, appID)
				return nil, errors.Wrapf(err, "failed to create bundle instance auth for bundle %s", b.ID)
			}
			result = append(result, instanceAuth)
		}
	}

	return result, nil
}

func (r *Resolver) requestBundleInstanceAuth(ctx context.Context, err error, bundleID string, in graphql.BundleInstanceAuthRequestInput, tx persistence.PersistenceTx) (*graphql.BundleInstanceAuth, error) {
	log.C(ctx).Infof("Requesting BundleInstanceAuth creation for Bundle with id %s", bundleID)
	bndl, err := r.bndlSvc.Get(ctx, bundleID)
	if err != nil {
		return nil, err
	}

	convertedIn := r.conv.RequestInputFromGraphQL(in)

	instanceAuthID, err := r.svc.Create(ctx, bundleID, convertedIn, bndl.DefaultInstanceAuth, bndl.InstanceAuthRequestInputSchema)
	if err != nil {
		return nil, err
	}

	instanceAuth, err := r.svc.Get(ctx, instanceAuthID)
	if err != nil {
		return nil, err
	}
	log.C(ctx).Infof("Successfully created BundleInstanceAuth with id %s for Bundle with id %s", instanceAuthID, bundleID)

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.conv.ToGraphQL(instanceAuth)
}

func bundlesIDMapping(bndls []*model.Bundle) map[string]string {
	bndlsMap := make(map[string]string, 0)
	for _, b := range bndls {
		if b.OrdID != nil {
			bndlsMap[*b.OrdID] = b.ID
		}
	}

	return bndlsMap
}
