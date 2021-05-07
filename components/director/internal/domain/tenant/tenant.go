package tenant

import (
	"context"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/header"
)

type key int

const (
	TenantContextKey key = iota
)

type TenantCtx struct {
	InternalID string
	ExternalID string
}

func LoadFromContext(ctx context.Context) (string, error) {
	tenant, ok := ctx.Value(TenantContextKey).(TenantCtx)

	if !ok {
		return "", apperrors.NewCannotReadTenantError()
	}

	if tenant.InternalID == "" {
		return "", apperrors.NewTenantRequiredError()
	}

	return tenant.InternalID, nil
}

// TODO [trap] external tenant is not part of request
func LoadExternalFromContext(ctx context.Context) (string, error) {
	tenant, ok := ctx.Value(TenantContextKey).(TenantCtx)

	if !ok {
		return "", apperrors.NewCannotReadTenantError()
	}

	if tenant.InternalID == "" {
		return "", apperrors.NewTenantRequiredError()
	}
	headers := ctx.Value(header.ContextKey)
	headers2, ok := (headers).(http.Header)
	if !ok {
		return "", nil
	}
	extTenant := headers2["Tenant"]
	if len(extTenant) == 0 {
		return "", nil
	}
	return extTenant[0], nil
}

func SaveToContext(ctx context.Context, internalID, externalID string) context.Context {
	tenantCtx := TenantCtx{InternalID: internalID, ExternalID: externalID}
	return context.WithValue(ctx, TenantContextKey, tenantCtx)
}
