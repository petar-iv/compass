package proxy

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/kyma-incubator/compass/components/gateway/internal/jwtutil"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/correlation"
	"github.com/kyma-incubator/compass/components/gateway/pkg/httpcommon"
	"github.com/pkg/errors"
)

//TODO: Create a single CommonTransport struct and provide implementations of
//RoundTrip for tenant-fetcher and director/connector structs which embed the CommonTransport
type TenantFetcherTransport struct {
	http.RoundTripper
	auditlogSink AuditlogService
	auditlogSvc  AuditlogService
}

func NewTenantFetcherTransport(sink AuditlogService, svc AuditlogService, trip RoundTrip) *Transport {
	return &Transport{
		RoundTripper: trip,
		auditlogSink: sink,
		auditlogSvc:  svc,
	}
}

func (t *TenantFetcherTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	tenant := gjson.GetBytes(requestBody, "globalAccountGUID").String()

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

	ctx := context.WithValue(req.Context(), correlation.RequestIDHeaderKey, correlationHeaders)
	err = preAuditLogger.PreLog(ctx, AuditlogMessage{
		CorrelationIDHeaders: correlationHeaders,
		Request:              string(requestBody),
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

	err = t.auditlogSink.Log(req.Context(), AuditlogMessage{
		CorrelationIDHeaders: correlationHeaders,
		Request:              string(requestBody),
		Response:             string(responseBody),
		Claims:               claims,
	})
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
