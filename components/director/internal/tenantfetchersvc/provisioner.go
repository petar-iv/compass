package tenantfetchersvc

import (
	"context"
	"database/sql"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	tenantEntity "github.com/kyma-incubator/compass/components/director/pkg/tenant"
	"github.com/pkg/errors"
)

// TenantService missing godoc
//go:generate mockery --name=TenantService --output=automock --outpkg=automock --case=underscore --unroll-variadic=False
type TenantService interface {
	GetInternalTenant(ctx context.Context, externalTenant string) (string, error)
	CreateManyIfNotExists(ctx context.Context, tenantInputs ...model.BusinessTenantMappingInput) error
}

type provisioner struct {
	tenantSvc TenantService
}

// NewTenantProvisioner missing godoc
func NewTenantProvisioner(tenantSvc TenantService) *provisioner {
	return &provisioner{
		tenantSvc: tenantSvc,
	}
}

// ProvisionTenant missing godoc
func (p *provisioner) ProvisionTenant(ctx context.Context, tenant model.BusinessTenantMappingInput) error {
	externalTenantID := tenant.ExternalTenant
	parentExternalID := tenant.Parent
	if len(parentExternalID) > 0 {
		parentInternalID, err := p.ensureParentExists(ctx, parentExternalID, externalTenantID)
		if err != nil {
			return errors.Wrapf(err, "while ensuring parent tenant with ID %s exists", parentExternalID)
		}
		tenant.Parent = parentInternalID
	}

	if err := p.tenantSvc.CreateManyIfNotExists(ctx, tenant); err != nil {
		if !apperrors.IsNotUniqueError(err) {
			return errors.Wrapf(err, tenantCreationFailureMsgFmt, externalTenantID)
		}
	}

	return nil
}

func (p *provisioner) ensureParentExists(ctx context.Context, parentTenantID, childTenantID string) (string, error) {
	log.C(ctx).Infof("Ensuring parent tenant with external ID %s for tenant with external ID %s exists", parentTenantID, childTenantID)
	id, err := p.tenantSvc.GetInternalTenant(ctx, parentTenantID)
	if err != nil && !apperrors.IsNotFoundError(err) && err != sql.ErrNoRows {
		return "", errors.Wrapf(err, "failed to retrieve internal ID of parent with external ID %s", parentTenantID)
	}
	if id != "" {
		log.C(ctx).Infof("Parent tenant with external ID %s already exists", parentTenantID)
		return id, nil
	}

	log.C(ctx).Infof("Creating parent tenant with external ID %s", parentTenantID)
	err = p.tenantSvc.CreateManyIfNotExists(ctx, p.customerTenant(parentTenantID))
	if err != nil && apperrors.IsNotUniqueError(err) {
		log.C(ctx).Infof("Parent tenant with external ID %s already exists", parentTenantID)
		return p.tenantSvc.GetInternalTenant(ctx, parentTenantID)
	} else if err != nil {
		return "", errors.Wrapf(err, "failed to create parent tenant with ID %s", parentTenantID)
	}

	internalID, err := p.tenantSvc.GetInternalTenant(ctx, parentTenantID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve internal ID of parent with external ID %s", parentTenantID)
	}

	log.C(ctx).Infof("Successfully created parent tenant with external ID %s and internal ID %s", parentTenantID, internalID)
	return internalID, nil
}

func (p *provisioner) customerTenant(tenantID string) model.BusinessTenantMappingInput {
	return model.BusinessTenantMappingInput{
		Name:           tenantID,
		ExternalTenant: tenantID,
		Parent:         "",
		Subdomain:      "",
		Type:           tenantEntity.TypeToStr(tenantEntity.Customer),
		Provider:       autogeneratedTenantProvider,
	}
}