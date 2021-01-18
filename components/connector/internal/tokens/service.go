package tokens

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/connector/internal/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type TokenType string

const (
	ApplicationToken TokenType = "Application"
	RuntimeToken     TokenType = "Runtime"
	CSRToken         TokenType = "Certificate"
)

const clientCredentialScopesPrefix = "clientCredentialsRegistrationScopes"
const applicationJSONType = "application/json"

var defaultGrantTypes = []string{"client_credentials"}

//go:generate mockery -name=Service
type Service interface {
	CreateToken(ctx context.Context, clientId string, tokenType TokenType) (string, apperrors.AppError)
	Resolve(ctx context.Context, token string) (TokenData, string, apperrors.AppError)
	Delete(ctx context.Context, clientID, token string)
}

type tokenService struct {
	publicBaseURL string
	adminBaseURL  string
	httpCli       *http.Client

	applicationTokenTTL time.Duration
	runtimeTokenTTL     time.Duration
	csrTokenTTL         time.Duration
}

func NewTokenService(adminBaseURL, publicBaseURL string,
	httpCli *http.Client,
	applicationTokenTTL time.Duration,
	runtimeTokenTTL time.Duration,
	csrTokenTTL time.Duration,
) *tokenService {
	return &tokenService{
		adminBaseURL:        adminBaseURL,
		publicBaseURL:       publicBaseURL,
		httpCli:             httpCli,
		applicationTokenTTL: applicationTokenTTL,
		runtimeTokenTTL:     runtimeTokenTTL,
		csrTokenTTL:         csrTokenTTL,
	}
}

func (svc *tokenService) CreateToken(ctx context.Context, clientID string, tokenType TokenType) (string, apperrors.AppError) {
	scopes := []string{string(tokenType), clientID}

	oauthClientID := uuid.New().String()
	clientSecret, err := svc.registerClient(ctx, oauthClientID, scopes)
	if err != nil {
		return "", apperrors.Internal("while registering client credentials in Hydra %s", err)
	}

	ccConfig := clientcredentials.Config{
		ClientID:     oauthClientID,
		ClientSecret: clientSecret,
		TokenURL:     svc.publicBaseURL + "/oauth2/token",
		Scopes:       scopes,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, svc.httpCli)
	token, err := ccConfig.Token(ctx)
	if err != nil {
		return "", apperrors.Internal("could not get token: %s", err)
	}

	tokenData := TokenData{
		Type:     tokenType,
		ClientId: clientID,
	}
	log.C(ctx).Debugf("Storing token for %s with id %s in the cache", tokenData.Type, tokenData.ClientId)

	return token.AccessToken, nil
}

func (svc *tokenService) Resolve(ctx context.Context, token string) (TokenData, string, apperrors.AppError) {
	introspectData, err := svc.tokenIntrospect(ctx, token)
	if err != nil {
		return TokenData{}, "", apperrors.Internal("failed to resolve token %s", err)
	}
	if !introspectData.Active {
		return TokenData{}, "", apperrors.Internal("token is not active")
	}
	scopes := strings.Split(introspectData.Scopes, " ")
	// TODO: Check which scope is for token type
	tokenType := TokenType(scopes[0])
	var tokenTTL time.Duration
	switch tokenType {
	case RuntimeToken:
		tokenTTL = svc.runtimeTokenTTL
	case ApplicationToken:
		tokenTTL = svc.applicationTokenTTL
	case CSRToken:
		tokenTTL = svc.csrTokenTTL
	}

	if time.Now().Add(-tokenTTL).Unix() > introspectData.IssuedAt {
		return TokenData{}, "", apperrors.Internal("token has expired")
	}

	return TokenData{
		Type:     ApplicationToken,
		ClientId: scopes[1],
	}, introspectData.ClientID, nil
}

func (svc *tokenService) Delete(ctx context.Context, clientID, token string) {
	svc.tokenRevoke(ctx, clientID)
}

type clientCredentialsRegistrationBody struct {
	GrantTypes []string `json:"grant_types"`
	ClientID   string   `json:"client_id"`
	Scope      string   `json:"scope"`
}

type clientCredentialsRegistrationResponse struct {
	ClientSecret string `json:"client_secret"`
}

type TokenIntrospectData struct {
	Active   bool   `json:"active"`
	ClientID string `json:"client_id"`
	IssuedAt int64  `json:"iat"`
	Scopes   string `json:"scope"`
}

func (s *tokenService) tokenRevoke(ctx context.Context, clientID string) error {
	headers := map[string][]string{
		"Accept": {"application/json"},
	}

	log.C(ctx).Info("Revoking token")

	resp, closeBody, err := s.doRequest(ctx, http.MethodDelete, s.adminBaseURL+"/clients/"+clientID, headers, nil)
	if err != nil {
		return err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var respData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		if err != nil {
			return errors.Wrap(err, "while decoding response body")
		}
		log.C(ctx).Errorf("Response from hydra: %+v", respData)
		return fmt.Errorf("invalid HTTP status code: received: %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (s *tokenService) tokenIntrospect(ctx context.Context, token string) (*TokenIntrospectData, error) {
	headers := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"application/json"},
	}

	v := url.Values{}
	v.Add("token", token)
	body := []byte(v.Encode())

	log.C(ctx).Info("Introspecting token", token, "body:", string(body))
	// req, err := http.NewRequest("POST", "/oauth2/introspect", bytes.NewBuffer(body))
	// req.Header = headers

	resp, closeBody, err := s.doRequest(ctx, http.MethodPost, s.adminBaseURL+"/oauth2/introspect", headers, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var respData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		if err != nil {
			return nil, errors.Wrap(err, "while decoding response body")
		}
		log.C(ctx).Errorf("Response from hydra: %+v", respData)
		return nil, fmt.Errorf("invalid HTTP status code: received: %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	var introspectResponse TokenIntrospectData
	err = json.NewDecoder(resp.Body).Decode(&introspectResponse)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding response body")
	}

	return &introspectResponse, nil
}

func (s *tokenService) registerClient(ctx context.Context, clientID string, scopes []string) (string, error) {
	log.C(ctx).Debugf("Registering client_id %s and client_secret in Hydra with scopes: %s", clientID, scopes)
	reqBody := &clientCredentialsRegistrationBody{
		GrantTypes: defaultGrantTypes,
		ClientID:   clientID,
		Scope:      strings.Join(scopes, " "),
	}

	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(&reqBody)
	if err != nil {
		return "", errors.Wrap(err, "while encoding body")
	}

	resp, closeBody, err := s.doRequest(ctx, http.MethodPost, s.adminBaseURL+"/clients", nil, buffer)
	if err != nil {
		return "", err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("invalid HTTP status code: received: %d, expected %d", resp.StatusCode, http.StatusCreated)
	}

	var registrationResp clientCredentialsRegistrationResponse
	err = json.NewDecoder(resp.Body).Decode(&registrationResp)
	if err != nil {
		return "", errors.Wrap(err, "while decoding response body")
	}

	log.C(ctx).Debugf("client_id %s and client_secret successfully registered in Hydra", clientID)
	return registrationResp.ClientSecret, nil
}

func (s *tokenService) doRequest(ctx context.Context, method string, endpoint string, headers http.Header, body io.Reader) (*http.Response, func(body io.ReadCloser), error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "while creating new request")
	}

	if headers == nil {
		headers = http.Header{}
	}
	req.Header = headers
	req.Header.Set("Accept", applicationJSONType)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", applicationJSONType)
	}

	resp, err := s.httpCli.Do(req)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "while doing request to %s", endpoint)
	}

	closeBodyFn := func(body io.ReadCloser) {
		if body == nil {
			return
		}
		_, err = io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			log.C(ctx).WithError(err).Error("An error has occurred while copying response body.")
		}

		err := body.Close()
		if err != nil {
			log.C(ctx).WithError(err).Error("An error has occurred while closing body.")
		}
	}

	return resp, closeBodyFn, nil
}
