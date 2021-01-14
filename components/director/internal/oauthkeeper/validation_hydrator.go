package oauthkeeper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/internal/httputils"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/pkg/errors"
)

type ValidationHydrator interface {
	ResolveConnectorTokenHeader(w http.ResponseWriter, r *http.Request)
}

type Service interface {
	GetByToken(ctx context.Context, token string) (*model.SystemAuth, error)
}

type validationHydrator struct {
	tokenService Service
	transact     persistence.Transactioner
}

func NewValidationHydrator(tokenService Service, transact persistence.Transactioner) ValidationHydrator {
	return &validationHydrator{
		tokenService: tokenService,
		transact:     transact,
	}
}

func (tvh *validationHydrator) ResolveConnectorTokenHeader(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tx, err := tvh.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to open db transaction")
		httputils.RespondWithError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error occured while resolving one time token"))
		return
	}
	defer tvh.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	var authSession AuthenticationSession
	err = json.NewDecoder(r.Body).Decode(&authSession)
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to decode request body")
		httputils.RespondWithError(ctx, w, http.StatusBadRequest, errors.Wrap(err, "failed to decode Authentication Session from body"))
		return
	}
	defer httputils.Close(ctx, r.Body)

	connectorToken := r.Header.Get(ConnectorTokenHeader)
	if connectorToken == "" {
		connectorToken = r.URL.Query().Get(ConnectorTokenQueryParam)
	}

	if connectorToken == "" {
		log.C(ctx).Info("Token not provided")
		respondWithAuthSession(ctx, w, authSession)
		return
	}

	log.C(ctx).Info("Trying to resolve token...")

	tokenData, err := tvh.tokenService.GetByToken(ctx, connectorToken)
	// TODO: Check if token has not expired
	if err != nil {
		log.C(ctx).Infof("Invalid token provided: %s", err.Error())
		respondWithAuthSession(ctx, w, authSession)
		return
	}

	if authSession.Header == nil {
		authSession.Header = map[string][]string{}
	}

	// TODO: Is this the clientID?
	authSession.Header.Add(ClientIdFromTokenHeader, tokenData.ID)

	// TODO: Implement the invalidation
	// if err := tvh.tokenService.Invalidate(connectorToken); err != nil {
	// 	httputils.RespondWithError(ctx, w, http.StatusInternalServerError, errors.New("could not invalidate token"))
	// 	return
	// }

	err = tx.Commit()
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to commit db transaction")
		httputils.RespondWithError(ctx, w, http.StatusInternalServerError, errors.New("unexpected error occured while resolving one time token"))
		return
	}

	log.C(ctx).Infof("Token for %s resolved successfully", tokenData.ID)
	respondWithAuthSession(ctx, w, authSession)
}

func respondWithAuthSession(ctx context.Context, w http.ResponseWriter, authSession AuthenticationSession) {
	httputils.RespondWithBody(ctx, w, http.StatusOK, authSession)
}
