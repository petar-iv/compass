package enhance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/header"
	"github.com/kyma-incubator/compass/components/director/pkg/operation"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/kyma-incubator/compass/components/director/pkg/webhook"
	"github.com/pkg/errors"
)

func NewAllpicationOperationEnhancer(whs WebhookService) *applicationOpEnhancer {
	return &applicationOpEnhancer{
		webhookSvc: whs,
	}
}

type WebhookService interface {
	ListAllApplicationWebhooks(ctx context.Context, applicationID string) ([]*model.Webhook, error)
}
type applicationOpEnhancer struct {
	webhookSvc WebhookService
}

func (e *applicationOpEnhancer) Enhance(ctx context.Context, tenantID string, operation *operation.Operation, resp interface{}, webhookType *graphql.WebhookType) error {
	entity, ok := resp.(graphql.Entity)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to cast the resource entity of type %T to type graphql.Entity", resp))
	}

	operation.ResourceID = entity.GetID()
	operation.ResourceType = entity.GetType()

	// TODO pre-submit func???
	//appConditionStatus, err := determineApplicationInProgressStatus(operation.OperationType)
	//if err != nil {
	//	return errors.Wrap(err, "while determining the application status condition")
	//}
	//
	//updateFunc, ok := d.resourceUpdaterFuncs[operation.ResourceType]
	//if !ok {
	//	return errors.Wrapf(err,"No update function provided for resource of type %s with id %s", entity.GetType(), entity.GetID())
	//}
	//if err := updateFunc(ctx, entity.GetID(), false, nil, *appConditionStatus); err != nil {
	//	return errors.Wrapf(err,"While updating resource %s with id %s and status condition %v: %v", entity.GetType(), entity.GetID(), appConditionStatus, err)
	//}

	if webhookType != nil {
		webhookIDs, err := e.prepareWebhookIDs(ctx, operation, *webhookType)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve webhooks")
		}
		operation.WebhookProviderID = operation.ResourceID
		operation.WebhookIDs = webhookIDs
	}

	requestObject, err := prepareRequestObject(ctx, resource.Application, tenantID, resp, nil)
	if err != nil {
		return errors.Wrap(err, "an error occurred while preparing request data")
	}

	operation.RequestObject = requestObject
	return nil
}

func prepareRequestObject(ctx context.Context, resType resource.Type, tenantID string, application interface{}, bndlInstanceAuth interface{}) (string, error) {
	appResource, ok := application.(webhook.Resource)
	if !ok {
		return "", errors.New("application entity is not a webhook provider")
	}
	bndlInstanceAuthResource, ok := bndlInstanceAuth.(webhook.Resource)
	if !ok {
		return "", errors.New("bndlInstanceAuth entity is not a webhook provider")
	}
	reqHeaders, ok := ctx.Value(header.ContextKey).(http.Header)
	if !ok {
		return "", errors.New("failed to retrieve request headers")
	}

	headers := make(map[string]string, 0)
	for key, value := range reqHeaders {
		headers[key] = value[0]
	}

	extTenant,err := tenant.LoadExternalFromContext(ctx)

	requestObject := &webhook.RequestObject{
		Application: appResource,
		BundleInstanceAuth: bndlInstanceAuthResource,
		Type:        resType,
		TenantID:    tenantID,
		ExternalTenantID: extTenant,
		Headers:     headers,
	}

	data, err := json.Marshal(requestObject)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (e *applicationOpEnhancer) prepareWebhookIDs(ctx context.Context, operation *operation.Operation, webhookType graphql.WebhookType) ([]string, error) {
	webhooks, err := e.webhookSvc.ListAllApplicationWebhooks(ctx, operation.ResourceID)
	if err != nil {
		return nil, err
	}

	return extractWebhookIDs(webhooks, webhookType)
}

func determineApplicationInProgressStatus(opType operation.OperationType) (*model.ApplicationStatusCondition, error) {
	var appStatusCondition model.ApplicationStatusCondition
	switch opType {
	case operation.OperationTypeCreate:
		appStatusCondition = model.ApplicationStatusConditionCreating
	case operation.OperationTypeUpdate:
		appStatusCondition = model.ApplicationStatusConditionUpdating
	case operation.OperationTypeDelete:
		appStatusCondition = model.ApplicationStatusConditionDeleting
	default:
		return nil, apperrors.NewInvalidStatusCondition(resource.Application)
	}

	return &appStatusCondition, nil
}

func extractWebhookIDs(webhooks []*model.Webhook, webhookType graphql.WebhookType) ([]string, error) {
	webhookIDs := make([]string, 0)
	for _, currWebhook := range webhooks {
		if graphql.WebhookType(currWebhook.Type) == webhookType {
			webhookIDs = append(webhookIDs, currWebhook.ID)
		}
	}

	if len(webhookIDs) > 1 {
		return nil, errors.New("multiple webhooks per operation are not supported")
	}

	return webhookIDs, nil
}
