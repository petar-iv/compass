/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operation

import (
	"context"
	"fmt"
	"math"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
	"github.com/pkg/errors"
)

const ModeParam = "mode"

type Enhance func(ctx context.Context, tenantID string, operation *Operation, resp interface{}, webhookType *graphql.WebhookType) error
type directive struct {
	transact             persistence.Transactioner
	operationEnhancer    map[resource.Type]Enhance
	resourceFetcherFuncs map[resource.Type]ResourceFetcherFunc
	resourceUpdaterFuncs map[resource.Type]ResourceUpdaterFunc
	tenantLoaderFunc     TenantLoaderFunc
	scheduler            Scheduler
}

// NewDirective creates a new handler struct responsible for the Async directive business logic
func NewDirective(transact persistence.Transactioner, operationEnhancer map[resource.Type]Enhance, fesourceFetcherFunc map[resource.Type]ResourceFetcherFunc, resourceUpdaterFuncs map[resource.Type]ResourceUpdaterFunc, tenantLoaderFunc TenantLoaderFunc, scheduler Scheduler) *directive {
	return &directive{
		transact:             transact,
		operationEnhancer: operationEnhancer,
		resourceFetcherFuncs: fesourceFetcherFunc,
		resourceUpdaterFuncs: resourceUpdaterFuncs,
		tenantLoaderFunc:     tenantLoaderFunc,
		scheduler:            scheduler,
	}
}

// HandleOperation enriches the request with an Operation information when the requesting mutation is annotated with the Async directive
func (d *directive) HandleOperation(ctx context.Context, _ interface{}, next gqlgen.Resolver, operationType graphql.OperationType, webhookType *graphql.WebhookType, idField *string) (res interface{}, err error) {
	resCtx := gqlgen.GetFieldContext(ctx)
	mode, err := getOperationMode(resCtx)
	if err != nil {
		return nil, err
	}

	resourceType, err := getResourceType(webhookType)
	if err != nil {
		return nil, err
	}

	ctx = SaveModeToContext(ctx, *mode)

	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while opening database transaction: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to initialize database operation")
	}
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	tenant, err := d.tenantLoaderFunc(ctx)
	if err != nil {
		return nil, apperrors.NewTenantRequiredError()
	}

	if err := d.concurrencyCheck(ctx, tenant, operationType, resourceType, resCtx, idField); err != nil {
		return nil, err
	}

	if *mode == graphql.OperationModeSync {
		return executeSyncOperation(ctx, next, tx)
	}

	operation := &Operation{
		OperationType:     OperationType(str.Title(operationType.String())),
		OperationCategory: resCtx.Field.Name,
		CorrelationID:     log.C(ctx).Data[log.FieldRequestID].(string),
	}

	ctx = SaveToContext(ctx, &[]*Operation{operation})
	operationsArr, _ := FromCtx(ctx)

	committed := false
	defer func() {
		if !committed {
			lastIndex := int(math.Max(0, float64(len(*operationsArr)-1)))
			*operationsArr = (*operationsArr)[:lastIndex]
		}
	}()

	resp, err := next(ctx)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while processing operation: %s", err.Error())
		return nil, err
	}

	entity, ok := resp.(graphql.Entity)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Failed to cast the resource entity of type %T to type graphql.Entity", resp))
	}
	enhancerFunc, ok := d.operationEnhancer[resourceType]
	if !ok {
		return nil, errors.New("unable to find enhancer func for type " + string(resourceType))
	}
	if err := enhancerFunc(ctx, tenant, operation, resp, webhookType); err != nil {
		return nil, err
	}

	// TODO generalize - preSubmitFunc
	if resourceType == resource.Application {
		appConditionStatus, err := determineApplicationInProgressStatus(operationType)
		if err != nil {
			log.C(ctx).WithError(err).Errorf("While determining the application status condition: %v", err)
			return nil, err
		}

		updateFunc, ok := d.resourceUpdaterFuncs[operation.ResourceType]
		if !ok {
			log.C(ctx).Errorf("No update function provided for resource of type %s with id %s", entity.GetType(), entity.GetID())
			return nil, apperrors.NewInternalError("Unable to update resource %s with id %s", entity.GetType(), entity.GetID())
		}
		if err := updateFunc(ctx, entity.GetID(), false, nil, *appConditionStatus); err != nil {
			log.C(ctx).WithError(err).Errorf("While updating resource %s with id %s and status condition %v: %v", entity.GetType(), entity.GetID(), appConditionStatus, err)
			return nil, apperrors.NewInternalError("Unable to update resource %s with id %s", entity.GetType(), entity.GetID())
		}
	}

	operationID, err := d.scheduler.Schedule(ctx, operation)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while scheduling operation: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to schedule operation")
	}

	operation.OperationID = operationID

	err = tx.Commit()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while closing database transaction: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to finalize database operation")
	}
	committed = true

	return resp, nil
}

func getResourceType(webhookType *graphql.WebhookType) (resource.Type, error) {
	switch *webhookType {
	case graphql.WebhookTypeUnregisterApplication, graphql.WebhookTypeRegisterApplication:
		return resource.Application, nil
	case graphql.WebhookTypeBundleInstanceAuthCreation, graphql.WebhookTypeBundleInstanceAuthDeletion:
		return resource.BundleInstanceAuth, nil
	default:
		return "", errors.New(fmt.Sprintf("webhook type %s is not bound to any type", webhookType))
	}
}

func (d *directive) concurrencyCheck(ctx context.Context, tenant string, op graphql.OperationType, resourceType resource.Type, resCtx *gqlgen.FieldContext, idField *string) error {
	if op == graphql.OperationTypeCreate {
		return nil
	}

	if idField == nil {
		return apperrors.NewInternalError("idField from context should not be empty")
	}

	resourceID, ok := resCtx.Args[*idField].(string)
	if !ok {
		return apperrors.NewInternalError(fmt.Sprintf("could not get idField: %q from request context", *idField))
	}

	app, err := d.resourceFetcherFuncs[resourceType](ctx, tenant, resourceID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return err
		}

		return apperrors.NewInternalError("failed to fetch resource with id %s", resourceID)
	}

	if app.GetDeletedAt().IsZero() && app.GetUpdatedAt().IsZero() && !app.GetReady() && isErrored(app) { // CREATING
		return apperrors.NewConcurrentOperationInProgressError("create operation is in progress")
	}
	if !app.GetDeletedAt().IsZero() && isErrored(app) { // DELETING
		return apperrors.NewConcurrentOperationInProgressError("delete operation is in progress")
	}
	// Note: This will be needed when there is async UPDATE supported
	// if app.DeletedAt.IsZero() && app.UpdatedAt.After(app.CreatedAt) && !app.Ready && *app.Error == "" { // UPDATING
	// 	return nil, apperrors.NewInvalidData	Error("another operation is in progress")
	// }

	return nil
}

func getOperationMode(resCtx *gqlgen.ResolverContext) (*graphql.OperationMode, error) {
	var mode graphql.OperationMode
	if _, found := resCtx.Args[ModeParam]; !found {
		mode = graphql.OperationModeSync
	} else {
		modePointer, ok := resCtx.Args[ModeParam].(*graphql.OperationMode)
		if !ok {
			return nil, apperrors.NewInternalError(fmt.Sprintf("could not get %s parameter", ModeParam))
		}
		mode = *modePointer
	}

	return &mode, nil
}

func executeSyncOperation(ctx context.Context, next gqlgen.Resolver, tx persistence.PersistenceTx) (interface{}, error) {
	resp, err := next(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while closing database transaction: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to finalize database operation")
	}

	return resp, nil
}

func isErrored(app model.Entity) bool {
	return app.GetError() == nil || *app.GetError() == ""
}

func determineApplicationInProgressStatus(opType graphql.OperationType) (*model.ApplicationStatusCondition, error) {
	var appStatusCondition model.ApplicationStatusCondition
	switch opType {
	case graphql.OperationTypeCreate:
		appStatusCondition = model.ApplicationStatusConditionCreating
	case graphql.OperationTypeUpdate:
		appStatusCondition = model.ApplicationStatusConditionUpdating
	case graphql.OperationTypeDelete:
		appStatusCondition = model.ApplicationStatusConditionDeleting
	default:
		return nil, apperrors.NewInvalidStatusCondition(resource.Application)
	}

	return &appStatusCondition, nil
}
