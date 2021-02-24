package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	tenantpkg "github.com/kyma-incubator/compass/components/director/pkg/tenant"

	"github.com/tidwall/gjson"

	"github.com/kyma-incubator/compass/components/gateway/internal/jwtutil"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/correlation"
	"github.com/kyma-incubator/compass/components/gateway/pkg/httpcommon"
	"github.com/pkg/errors"
)

type RestAuditlogService interface {
	AuditlogService
	LogRest(ctx context.Context, msg AuditlogMessage, responseStatus int, ip *net.IP) error
}

//TODO: Create a single CommonTransport struct and provide implementations of
//RoundTrip for tenant-fetcher and director/connector structs which embed the CommonTransport
type TenantFetcherTransport struct {
	http.RoundTripper
	auditlogSink                   AuditlogService
	auditlogSvc                    AuditlogService
	TenantProviderTenantIdProperty string
	TenantProvider                 string
}

func NewTenantFetcherTransport(sink AuditlogService, svc AuditlogService, trip RoundTrip, tenantProvider, tenantProviderProperty string) *TenantFetcherTransport {
	return &TenantFetcherTransport{
		RoundTripper:                   trip,
		auditlogSink:                   sink,
		auditlogSvc:                    svc,
		TenantProviderTenantIdProperty: tenantProviderProperty,
		TenantProvider:                 tenantProvider,
	}
}

type TenantBody struct {
	Method         string                 `json:"method"`
	TenantId       string                 `json:"tenant_id"`
	TenantProvider string                 `json:"tenant_provider"`
	Status         tenantpkg.TenantStatus `json:"tenant_status"`
}

func (t *TenantFetcherTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	tenant := gjson.GetBytes(requestBody, t.TenantProviderTenantIdProperty).String()

	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	defer httpcommon.CloseBody(req.Context(), req.Body)

	correlationHeaders := correlation.HeadersForRequest(req)

	preAuditLogger, ok := t.auditlogSvc.(PreAuditlogService)
	if !ok {
		return nil, errors.New("Failed to type cast PreAuditlogService")
	}

	claims, err := t.getClaims(req.Header, tenant)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing JWT")
	}

	body := TenantBody{
		Method:         req.Method,
		TenantId:       tenant,
		TenantProvider: t.TenantProvider,
		Status:         tenantpkg.Active,
	}

	marshaledBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "while marshaling tenant body")
	}

	ctx := context.WithValue(req.Context(), correlation.RequestIDHeaderKey, correlationHeaders)
	err = preAuditLogger.PreLog(ctx, AuditlogMessage{
		CorrelationIDHeaders: correlationHeaders,
		Request:              string(marshaledBody),
		Response:             "",
		Claims:               claims,
	})
	if err != nil {
		return nil, errors.Wrap(err, "while sending pre-change auditlog message to auditlog service")
	}

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, errors.Wrap(err, "on request round trip")
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(responseBody))
	defer httpcommon.CloseBody(req.Context(), resp.Body)

	restAuditLogger, ok := t.auditlogSvc.(RestAuditlogService)
	if !ok {
		return nil, errors.New("Failed to type cast PreAuditlogService")
	}

	ipString := req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]
	ip := net.ParseIP(ipString)

	err = restAuditLogger.LogRest(req.Context(), AuditlogMessage{
		CorrelationIDHeaders: correlationHeaders,
		Request:              string(marshaledBody),
		Response:             string(responseBody),
		Claims:               claims,
	}, resp.StatusCode, &ip)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("failed to send a post-change auditlog message to auditlog service")
	}

	return resp, nil
}

func (t *TenantFetcherTransport) getClaims(headers http.Header, tenant string) (Claims, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return Claims{}, errors.New("no bearer token")
	}
	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := jwtutil.NewClaims(token)

	return Claims{
		Tenant:       tenant,
		Scopes:       claims.Scopes,
		ConsumerID:   claims.ClientId,
		ConsumerType: "XSUAA",
	}, err
}
