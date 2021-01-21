package tokens

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/connector/internal/apperrors"
	gcli "github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

const requestForCSRToken = `
		mutation { generateCSRToken (authID:"%s")
		  {
			token
		  }
		}`

//go:generate mockery -name=Service
type Service interface {
	CreateToken(ctx context.Context, clientId string, tokenType TokenType) (string, apperrors.AppError)
	Resolve(token string) (TokenData, apperrors.AppError)
	Delete(token string)
}

//go:generate mockery -name=GraphQLClient -output=automock -outpkg=automock -case=underscore
type GraphQLClient interface {
	Run(ctx context.Context, req *gcli.Request, resp interface{}) error
}

type tokenService struct {
	cli GraphQLClient
}

func NewTokenService(cli GraphQLClient, store Cache, generator TokenGenerator) *tokenService {
	return &tokenService{
		cli: cli,
	}
}

func (svc *tokenService) CreateToken(ctx context.Context, clientId string, tokenType TokenType) (string, apperrors.AppError) {
	token, err := svc.getOneTimeToken(ctx, clientId)
	if err != nil {
		return "", apperrors.Internal("could not get one time token: %s", err)
	}
	return token, nil
}

func (svc *tokenService) Resolve(token string) (TokenData, apperrors.AppError) {
	tokenData, err := svc.store.Get(token)
	if err != nil {
		return TokenData{}, err.Append("Failed to resolve token")
	}

	return tokenData, nil
}

func (svc *tokenService) Delete(token string) {
	svc.store.Delete(token)
}

func (s *tokenService) getOneTimeToken(ctx context.Context, id string) (string, error) {
	req := gcli.NewRequest(fmt.Sprintf(requestForCSRToken, id))

	var resp map[string]map[string]interface{}
	err := s.cli.Run(ctx, req, &resp)
	if err != nil {
		return "", errors.Wrapf(err, "while calling director for CSR one time token")
	}

	return resp["generateCSRToken"]["token"].(string), nil
}
