package appmetadatavalidation

import (
	"context"
	"fmt"
	gqlgen "github.com/99designs/gqlgen/graphql"
	tnt "github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/consumer"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"
)

// LabelService is responsible to service-related label operations
//go:generate mockery --name=LabelService --output=automock --outpkg=automock --case=underscore --disable-version-string
type LabelService interface {
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
}

// TenantService is responsible to service-related tenant operations
//go:generate mockery --name=TenantService --output=automock --outpkg=automock --case=underscore --disable-version-string
type TenantService interface {
	GetTenantByExternalID(ctx context.Context, id string) (*model.BusinessTenantMapping, error)
}

// Handler is an object with dependencies
type directive struct {
	transact  persistence.Transactioner
	tenantSvc TenantService
	labelSvc  LabelService
}

// NewHandler is a constructor for Handler object
func NewDirective(transact persistence.Transactioner, tenantSvc TenantService, labelSvc LabelService) *directive {
	return &directive{
		transact:  transact,
		tenantSvc: tenantSvc,
		labelSvc:  labelSvc,
	}
}

// Validate is a middleware that checks if the flow is oathkeeper.CertificateFlow and consumer type is consumer.ExternalCertificate.
// If it's not - continue with next handler. If it is, get the consumer tenant's Region label and the Region label of the tenant header.
// They have to match in order to continue with the next handler, otherwise fail the request
func (d *directive) Validate(ctx context.Context, obj interface{}, next gqlgen.Resolver) (interface{}, error) {
	logger := log.C(ctx)

	consumerInfo, err := consumer.LoadFromContext(ctx)
	if err != nil {
		logger.WithError(err).Errorf("An error has occurred while fetching consumer info from context: %v", err)
		return next(ctx)
	}

	if !d.isExternalCertificateFlow(consumerInfo) {
		logger.Infof("Flow is not external certificate. Continue with next handler.")
		return next(ctx)

	}

	tenantID, ok := ctx.Value(TenantHeader).(string)
	if !ok {
		return "", apperrors.NewCannotReadTenantError()
	}

	if tenantID == "" {
		logger.Infof("Missing tenant header. Continue with next handler.")
		return next(ctx)
	}

	logger.Infof("Flow is external certificate")

	tx, err := d.transact.Begin()
	if err != nil {
		logger.WithError(err).Errorf("An error has occurred while opening transaction: %v", err)
		return nil, apperrors.InternalErrorFrom(err, "while opening db transaction")

	}
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	consumerRegionLabel, err := d.getTenantRegionLabelValue(ctx, consumerInfo.ConsumerID)
	if err != nil {
		return nil, apperrors.InternalErrorFrom(err, "while fetching consumer tenant %q region label", consumerInfo.ConsumerID)
	}

	headerTenantRegionLabel, err := d.getTenantRegionLabelValue(ctx, tenantID)
	if err != nil {
		return nil, apperrors.InternalErrorFrom(err, "while fetching tenant header %q region label", tenantID)
	}

	if consumerRegionLabel != headerTenantRegionLabel {
		err = errors.New(fmt.Sprintf("region labels mismatch: %q and %q for tenants %q and %q", consumerRegionLabel, headerTenantRegionLabel, consumerInfo.ConsumerID, tenantID))
		logger.WithError(err).Errorf("Regions for external tenants %q and %q do not match: %q and %q", consumerInfo.ConsumerID, tenantID, consumerRegionLabel, headerTenantRegionLabel)
		return nil, apperrors.InternalErrorFrom(err, "region labels for tenants do not match")
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Errorf("An error has occurred while committing transaction: %v", err)
		return nil, apperrors.InternalErrorFrom(err, "error has occurred while committing transaction")
	}

	logger.Infof("Regions for tenants with external ID %q and %q matched. Continuing with request", consumerInfo.ConsumerID, tenantID)

	return next(ctx)
}

func (h *directive) getTenantRegionLabelValue(ctx context.Context, tenantID string) (interface{}, error) {
	logger := log.C(ctx)

	tenantModel, err := h.tenantSvc.GetTenantByExternalID(ctx, tenantID)
	if err != nil {
		logger.WithError(err).Errorf("An error has occurred while fetching tenant by external ID %q: %v", tenantID, err)
		return nil, errors.Wrapf(err, "while fetching tenant by external ID %q", tenantID)
	}

	tenantRegionLabel, err := h.labelSvc.GetByKey(ctx, tenantModel.ID, model.TenantLabelableObject, tenantModel.ID, tnt.RegionLabelKey)
	if err != nil {
		logger.WithError(err).Errorf("An error has occurred while fetching %q label for tenant ID %q: %v", tnt.RegionLabelKey, tenantID, err)
		return nil, errors.Wrapf(err, "while fetching %q label tenant by external ID %q", tnt.RegionLabelKey, tenantID)
	}

	return tenantRegionLabel.Value, nil
}

func (h *directive) isExternalCertificateFlow(consumerInfo consumer.Consumer) bool {
	return consumerInfo.Flow.IsCertFlow() && consumerInfo.ConsumerType == consumer.ExternalCertificate
}
