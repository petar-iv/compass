package tenantmapping

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/authenticator"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/internal/consumer"
	"github.com/kyma-incubator/compass/components/director/internal/oathkeeper"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=RuntimeService --output=automock --outpkg=automock --case=underscore --exported=true
type runtimeService interface {
	ListByFiltersGlobal(context.Context, []*labelfilter.LabelFilter) ([]*model.Runtime, error)
}

// NewCertServiceContextProvider implements the ObjectContextProvider interface by looking for tenant information directly populated in the certificate.
func NewCertServiceContextProvider(tenantRepo TenantRepository, runtimeSvc runtimeService, scopesGetter ScopesGetter) *certServiceContextProvider {
	return &certServiceContextProvider{
		tenantRepo: tenantRepo,
		tenantKeys: KeysExtra{
			TenantKey:         ProviderTenantKey,
			ExternalTenantKey: ProviderExternalTenantKey,
		},
		scopesGetter: scopesGetter,
		runtimeSvc:   runtimeSvc,
	}
}

type certServiceContextProvider struct {
	tenantRepo   TenantRepository
	tenantKeys   KeysExtra
	scopesGetter ScopesGetter
	runtimeSvc   runtimeService
}

// GetObjectContext is the certServiceContextProvider implementation of the ObjectContextProvider interface
// By using trusted external certificate issuer we assume that we will receive the tenant information extracted from the certificate.
// There we should only convert the tenant identifier from external to internal. Additionally, we mark the consumer in this flow as Runtime.
func (p *certServiceContextProvider) GetObjectContext(ctx context.Context, reqData oathkeeper.ReqData, authDetails oathkeeper.AuthDetails) (ObjectContext, error) {
	logger := log.C(ctx).WithFields(logrus.Fields{
		"consumer_type": consumer.Runtime,
	})

	ctx = log.ContextWithLogger(ctx, logger)

	externalTenantID := authDetails.AuthID

	matchedComponentName, ok := reqData.Body.Header[authenticator.ComponentName]
	if !ok || len(matchedComponentName) == 0 {
		return ObjectContext{}, errors.New("empty matched component header")
	}

	scopes := ""
	var tc TenantContext
	// This if is needed to separate the director from ord flow because for the director flow we need to use the internal ID of the subaccount
	// whereas in the ord flow we expect external IDs in order ord views to work properly(using Automatic Scenario Assignments)
	log.C(ctx).Infof("Matched component name is %s", matchedComponentName[0])
	if matchedComponentName[0] == "director" { // Director Flow, do the tenant conversion
		log.C(ctx).Infof("Getting the tenant with external ID: %s", externalTenantID)
		tenantMapping, err := p.tenantRepo.GetByExternalTenant(ctx, externalTenantID)
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				log.C(ctx).Warningf("Could not find tenant with external ID: %s, error: %s", externalTenantID, err.Error())
				log.C(ctx).Infof("Returning tenant context with empty internal tenant ID and external ID %s", externalTenantID)
				return NewObjectContext(NewTenantContext(externalTenantID, ""), p.tenantKeys, scopes, authDetails.Region, "", authDetails.AuthID, authDetails.AuthFlow, consumer.Runtime, CertServiceObjectContextProvider), nil
			}
			return ObjectContext{}, errors.Wrapf(err, "while getting external tenant mapping [ExternalTenantID=%s]", externalTenantID)
		}

		tc = NewTenantContext(externalTenantID, tenantMapping.ID)
		scopes, err = p.directorScopes(ctx, authDetails)
		if err != nil {
			return ObjectContext{}, err
		}

		authDetails.AuthID, err = p.extractAuthIDFromRuntime(ctx, authDetails.AuthID)
		if err != nil {
			log.C(ctx).WithError(err).Errorf("Failed to retrieve auth ID from runtime: %v")
			return ObjectContext{}, err
		}
	} else {
		// ORD Flow, set the external tenant ID both for internal and external tenants
		tc = NewTenantContext(externalTenantID, externalTenantID)
	}

	objCtx := NewObjectContext(tc, p.tenantKeys, scopes, authDetails.Region, "", authDetails.AuthID, authDetails.AuthFlow, consumer.Runtime, CertServiceObjectContextProvider)
	log.C(ctx).Infof("Successfully got object context: %+v", objCtx)
	return objCtx, nil
}

// Match checks if there is "client-id-from-certificate" Header with nonempty value and "client-certificate-issuer" Header with value "certificate-service".
// If so AuthDetails object is build.
func (p *certServiceContextProvider) Match(_ context.Context, data oathkeeper.ReqData) (bool, *oathkeeper.AuthDetails, error) {
	idVal := data.Body.Header.Get(oathkeeper.ClientIDCertKey)
	certIssuer := data.Body.Header.Get(oathkeeper.ClientIDCertIssuer)

	if idVal != "" && certIssuer == oathkeeper.ExternalIssuer {
		return true, &oathkeeper.AuthDetails{AuthID: idVal, AuthFlow: oathkeeper.CertificateFlow, CertIssuer: certIssuer}, nil
	}

	return false, nil, nil
}

func (p *certServiceContextProvider) directorScopes(ctx context.Context, authDetails oathkeeper.AuthDetails) (string, error) {
	if authDetails.Authenticator != nil {
		log.C(ctx).Infof("")
	}
	declaredScopes, err := p.scopesGetter.GetRequiredScopes(buildPath(model.RuntimeReference))
	if err != nil {
		return "", errors.Wrap(err, "while fetching scopes")
	}
	return strings.Join(declaredScopes, " "), nil
}

func (p *certServiceContextProvider) extractAuthIDFromRuntime(ctx context.Context, currentAuthID string) (string, error) {
	runtimes, err := p.runtimeSvc.ListByFiltersGlobal(ctx, []*labelfilter.LabelFilter{})
	if err != nil {
		return "", err
	}

	runtimesCnt := len(runtimes)
	if runtimesCnt > 1 {
		return "", fmt.Errorf("expected 1 runtime, got %d", runtimesCnt)
	}
	if runtimesCnt == 1 {
		return runtimes[0].ID, nil
	}

	return currentAuthID, nil
}
