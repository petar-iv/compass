package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kyma-incubator/compass/components/director/pkg/cert"
	timeouthandler "github.com/kyma-incubator/compass/components/director/pkg/handler"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
	"net/http"
	"sync"
	"time"
)

func PrepareHydratorServer(cfg Config, revokedCertsCache Cache, middlewares ...mux.MiddlewareFunc) (*http.Server, error) {
	subjectProcessor, err := NewProcessor(cfg.SubjectConsumerMappingConfig, cfg.ExternalIssuerSubject.OrganizationalUnitPattern)
	if err != nil {
		return nil, err
	}

	//TODO these no longer accept "constants" but just config
	externalCertHeaderParser := NewHeaderParser(cfg.CertificateDataHeader, ExternalIssuer,
		ExternalCertIssuerSubjectMatcher(cfg), subjectProcessor.AuthIDFromSubjectFunc(), subjectProcessor.AuthSessionExtraFromSubjectFunc())
	connectorCertHeaderParser := NewHeaderParser(cfg.CertificateDataHeader, ConnectorIssuer,
		ConnectorCertificateSubjectMatcher(cfg), cert.GetCommonName, subjectProcessor.EmptyAuthSessionExtraFunc())

	//TODO move logic of repository .Contains into validation hydrator
	validationHydrator := NewValidationHydrator(revokedCertsCache, connectorCertHeaderParser, externalCertHeaderParser)

	router := mux.NewRouter()
	router.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Use(middlewares...)

	v1Router := router.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/certificate/data/resolve", validationHydrator.ResolveIstioCertHeader)

	handlerWithTimeout, err := timeouthandler.WithTimeout(router, cfg.ServerTimeout)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              cfg.HydratorAddress,
		Handler:           handlerWithTimeout,
		ReadHeaderTimeout: cfg.ServerTimeout,
	}, nil
}

func startServer(parentCtx context.Context, server *http.Server, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		stopServer(server)
	}()

	log.C(ctx).Infof("Starting and listening on %s://%s", "http", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.C(ctx).Fatalf("Could not listen on %s://%s: %v\n", "http", server.Addr, err)
	}
}

func stopServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	go func(ctx context.Context) {
		<-ctx.Done()

		if ctx.Err() == context.Canceled {
			return
		} else if ctx.Err() == context.DeadlineExceeded {
			log.C(ctx).Panic("Timeout while stopping the server, killing instance!")
		}
	}(ctx)

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		log.C(ctx).Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}

func exitOnError(err error, context string) {
	if err != nil {
		wrappedError := errors.Wrap(err, context)
		log.D().Fatal(wrappedError)
	}
}
