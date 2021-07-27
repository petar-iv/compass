package authenticator

import (
	"context"
	"encoding/json"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"sync"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
)

const (
	JwksKeyIDKey = "kid"
)

type Error struct {
	Message string `json:"message"`
}

type Authenticator struct {
	mux                        sync.RWMutex
	cachedJWKS                 jwk.Set
	jwksEndpoints              []string
	zoneId                     string
	trustedClaimPrefixes       []string
	subscriptionCallbacksScope string
	allowJWTSigningNone        bool
}

func New(jwksEndpoints []string, zoneId, subscriptionCallbacksScope string, trustedClaimPrefixes []string, allowJWTSigningNone bool) *Authenticator {
	return &Authenticator{
		jwksEndpoints:              jwksEndpoints,
		zoneId:                     zoneId,
		trustedClaimPrefixes:       trustedClaimPrefixes,
		subscriptionCallbacksScope: subscriptionCallbacksScope,
		allowJWTSigningNone:        allowJWTSigningNone,
	}
}

func (a *Authenticator) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tokenString, err := a.getBearerToken(r)
			if err != nil {
				log.C(ctx).WithError(err).Errorf("An error has occurred while extracting the JWT token. Error code: %d: %v", http.StatusUnauthorized, err)
				a.writeAppError(ctx, w, err, http.StatusBadRequest)
				return
			}

			log.C(ctx).Warn(tokenString)

			zidClaim := gjson.Get(tokenString, "zid").String()
			if zidClaim != a.zoneId {
				log.C(ctx).Errorf(`Zone id "%s" from user token does not match the trusted zone %s`, zidClaim, a.zoneId)
				a.writeAppError(ctx, w, errors.Errorf(`Zone id "%s" from user token is not trusted`, zidClaim), http.StatusUnauthorized)
				return
			}

			scopes := gjson.Get(tokenString, "scope").Array()
			for _, scope := range scopes {
				if scope.String() == a.subscriptionCallbacksScope {
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}

			log.C(ctx).Errorf(`Scopes from token "%v" does not match expected scope %s`, scopes, a.zoneId)
			a.writeAppError(ctx, w, errors.Errorf(`Scopes from token "%v" does not match expected scope %s`, scopes, a.zoneId), http.StatusUnauthorized)
			return
		})
	}
}

func (a *Authenticator) getBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("X-Authorization")
	if reqToken == "" {
		return "", apperrors.NewUnauthorizedError("invalid bearer token")
	}

	reqToken = strings.TrimPrefix(reqToken, "Bearer ")
	return reqToken, nil
}

func (a *Authenticator) writeAppError(ctx context.Context, w http.ResponseWriter, appErr error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(Error{Message: appErr.Error()})
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while encoding data: %v", err)
	}
}
