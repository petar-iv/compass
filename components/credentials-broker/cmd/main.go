/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"credentials-broker/internal/config"
	credentials_provider "credentials-broker/internal/credentialsprovider"
	"credentials-broker/internal/osb"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"os"

	"github.com/kyma-incubator/compass/components/director/pkg/signal"

	correlation "github.com/kyma-incubator/compass/components/director/pkg/http"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/env"
	httputil "github.com/kyma-incubator/compass/components/system-broker/pkg/http"
	lagger_adapter "github.com/kyma-incubator/compass/components/system-broker/pkg/log"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/oauth"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/server"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/uuid"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	env, err := env.Default(ctx, config.AddPFlags)
	fatalOnError(err)

	cfg, err := config.New(env)
	fatalOnError(err)

	err = cfg.Validate()
	fatalOnError(err)

	ctx, err = log.Configure(ctx, cfg.Log)
	fatalOnError(err)

	uuidSrv := uuid.NewService()

	credentialsProviderClient, err := prepareCredentialsProviderClient(cfg)
	fatalOnError(err)

	credentialsBroker := osb.NewCredentialsBroker(credentialsProviderClient)
	osbApi := osb.API(cfg.Server.RootAPI, credentialsBroker, lagger_adapter.NewDefaultLagerAdapter())
	srv := server.New(cfg.Server, uuidSrv, osbApi)

	srv.Start(ctx)
}

func fatalOnError(err error) {
	if err != nil {
		log.D().Fatal(err.Error())
	}
}

func prepareCredentialsProviderClient(cfg *config.Config) (credentials_provider.Client, error) {
	// prepare raw http transport and http client based on cfg
	httpTransport := correlation.NewCorrelationIDTransport(httputil.NewErrorHandlerTransport(httputil.NewHTTPTransport(cfg.HttpClient)))
	httpClient := httputil.NewClient(cfg.HttpClient.Timeout, httpTransport)

	oauthTokenProvider, err := oauth.NewTokenProvider(cfg.OAuthProvider, httpClient)
	if err != nil {
		return nil, err
	}

	securedTransport := httputil.NewSecuredTransport(cfg.HttpClient.Timeout, httpTransport, oauthTokenProvider)
	securedClient := httputil.NewClient(cfg.HttpClient.Timeout, securedTransport)

	return credentials_provider.NewClient(cfg.CredentialsProvider, securedClient), nil
}
