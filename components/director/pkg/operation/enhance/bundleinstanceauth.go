package enhance

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/operation"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
)

type BundleService interface {
	Get(ctx context.Context, id string) (*model.Bundle, error)
}

type BundleInstanceAuthService interface {
	Get(ctx context.Context, id string) (*model.BundleInstanceAuth, error)
}

type ApplicationService interface {
	Get(ctx context.Context, id string) (*model.Application, error)
}

func NewIntanceAuthOperationEnhancer(webhookSvc WebhookService, bundleSvc BundleService, bundleInstanceAuth BundleInstanceAuthService, appSvc ApplicationService) instanceAuthOpEnhancer {
	return instanceAuthOpEnhancer{
		webhookSvc:         webhookSvc,
		bundleSvc:          bundleSvc,
		bundleInstanceAuth: bundleInstanceAuth,
		appSvc:             appSvc,
	}
}

type instanceAuthOpEnhancer struct {
	webhookSvc         WebhookService
	bundleSvc          BundleService
	bundleInstanceAuth BundleInstanceAuthService
	appSvc             ApplicationService
}

func (e *instanceAuthOpEnhancer) Enhance(ctx context.Context, tenantID string, operation *operation.Operation, resp interface{}, webhookType *graphql.WebhookType) error {
	entity, ok := resp.(*graphql.BundleInstanceAuth)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to cast the resource entity of type %T to type *graphql.BundleInstanceAuth", resp))
	}

	operation.ResourceID = entity.GetID()
	operation.ResourceType = entity.GetType()

	instanceAuth, err := e.bundleInstanceAuth.Get(ctx, operation.ResourceID)
	if err != nil {
		return err
	}

	bundle, err := e.bundleSvc.Get(ctx, instanceAuth.BundleID)
	if err != nil {
		return err
	}

	if webhookType != nil {
		webhooks, err := e.webhookSvc.ListAllApplicationWebhooks(ctx, bundle.ApplicationID)
		webhookIDs, err := extractWebhookIDs(webhooks, *webhookType)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve webhooks")
		}
		operation.WebhookIDs = webhookIDs
	}

	app, err := e.appSvc.Get(ctx, bundle.ApplicationID)
	operation.WebhookProviderID = app.ID

	requestObject, err := prepareRequestObject(ctx, resource.BundleInstanceAuth, tenantID, application.NewConverter(nil, nil).ToGraphQL(app), entity)
	if err != nil {
		return errors.Wrap(err, "an error occurred while preparing request data")
	}

	operation.RequestObject = requestObject
	return nil
}
