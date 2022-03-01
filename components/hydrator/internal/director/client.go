package director

import (
	"context"
	"time"

	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/retry"
	schema "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/machinebox/graphql"
)

//go:generate mockery --name=Client --output=automock --outpkg=automock --case=underscore
type Client interface {
	GetTenantByExternalID(ctx context.Context, tenantID string) (schema.Tenant, apperrors.AppError)
	GetTenantByInternalID(ctx context.Context, tenantID string) (schema.Tenant, apperrors.AppError)
	GetTenantByLowestOwnerForResource(ctx context.Context, resourceID, resourceType string) (string, apperrors.AppError)
	GetSystemAuthByID(ctx context.Context, authID string) (schema.SystemAuth, apperrors.AppError)
	GetSystemAuthByToken(ctx context.Context, token string) (schema.SystemAuth, apperrors.AppError)
	UpdateSystemAuth(ctx context.Context, authID string, auth schema.Auth) (UpdateAuthResult, apperrors.AppError)
	InvalidateSystemAuthOneTimeToken(ctx context.Context, authID string) (schema.SystemAuth, apperrors.AppError)
	GetRuntimeByTokenIssuer(ctx context.Context, issuer string) (schema.Runtime, apperrors.AppError)
}

type Config struct {
	DirectorURL   string        `envconfig:"default=http://127.0.0.1:3000/graphql"`
	ClientTimeout time.Duration `envconfig:"default=115s"`
}

func NewClient(gqlClient *graphql.Client) Client {
	return &client{
		gqlClient: gqlClient,
		timeout:   30 * time.Second,
	}
}

type client struct {
	gqlClient *graphql.Client
	timeout   time.Duration
}

type TenantResponse struct {
	Result schema.Tenant `json:"result"`
}

type TenantByLowestOwnerForResourceResponse struct {
	Result string `json:"result"`
}

type SystemAuthResponse struct {
	Result schema.SystemAuth `json:"result"`
}

type RuntimeResponse struct {
	Result schema.Runtime `json:"result"`
}

type UpdateAuthResult struct {
	ID string `json:"id"`
}
type UpdateSystemAuthResponse struct {
	Result UpdateAuthResult `json:"result"`
}

func (c *client) GetTenantByExternalID(ctx context.Context, tenantID string) (schema.Tenant, apperrors.AppError) {
	query := tenantByExternalIDQuery(tenantID)
	var response TenantResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return schema.Tenant{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) GetTenantByInternalID(ctx context.Context, tenantID string) (schema.Tenant, apperrors.AppError) {
	query := tenantByInternalIDQuery(tenantID)
	var response TenantResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return schema.Tenant{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) GetTenantByLowestOwnerForResource(ctx context.Context, resourceID, resourceType string) (string, apperrors.AppError) {
	query := tenantByLowestOwnerForResourceQuery(resourceID, resourceType)
	var response TenantByLowestOwnerForResourceResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return "", apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) GetSystemAuthByID(ctx context.Context, authID string) (schema.SystemAuth, apperrors.AppError) {
	query := systemAuthQuery(authID)

	var response SystemAuthResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return schema.AppSystemAuth{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) GetSystemAuthByToken(ctx context.Context, token string) (schema.SystemAuth, apperrors.AppError) {
	query := systemAuthByTokenQuery(token)

	var response SystemAuthResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return schema.AppSystemAuth{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) UpdateSystemAuth(ctx context.Context, authID string, auth schema.Auth) (UpdateAuthResult, apperrors.AppError) {
	query, err := updateSystemAuthQuery(authID, auth)
	if err != nil {
		return UpdateAuthResult{}, apperrors.Internal(err.Error())
	}

	var response UpdateSystemAuthResponse

	err = c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return UpdateAuthResult{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) InvalidateSystemAuthOneTimeToken(ctx context.Context, authID string) (schema.SystemAuth, apperrors.AppError) {
	query := invalidateSystemAuthOneTimeTokenQuery(authID)

	var response SystemAuthResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return schema.AppSystemAuth{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) GetRuntimeByTokenIssuer(ctx context.Context, issuer string) (schema.Runtime, apperrors.AppError) {
	query := runtimeByTokenIssuerQuery(issuer)
	var response RuntimeResponse

	err := c.execute(ctx, c.gqlClient, query, &response)
	if err != nil {
		return schema.Runtime{}, apperrors.Internal(err.Error())
	}

	return response.Result, nil
}

func (c *client) execute(ctx context.Context, client *graphql.Client, query string, res interface{}) error {
	req := graphql.NewRequest(query)

	newCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return retry.GQLRun(client.Run, newCtx, req, res)
}