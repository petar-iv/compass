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
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	mp_bundle "github.com/kyma-incubator/compass/components/director/internal2/domain/bundle"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/kyma-incubator/compass/components/director/pkg/normalizer"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/kyma-incubator/compass/components/formation-watcher/notifications"

	"github.com/kyma-incubator/compass/components/director/internal2/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/auth"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/document"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/fetchrequest"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/integrationsystem"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/labeldef"
	rt "github.com/kyma-incubator/compass/components/director/internal2/domain/runtime"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/scenarioassignment"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/spec"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/version"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/webhook"
	"github.com/kyma-incubator/compass/components/director/internal2/uid"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/formation-watcher/config"
	"github.com/kyma-incubator/compass/components/formation-watcher/pkg/destination"
	"github.com/kyma-incubator/compass/components/formation-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/formation-watcher/pkg/signal"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = path.Dir(b)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	cfg := config.DefaultConfig()
	// err := envconfig.InitWithPrefix(&cfg, "APP")
	// fatalOnError(err)

	err := cfg.Validate()
	fatalOnError(err)

	ctx, err = log.Configure(ctx, cfg.Log)
	fatalOnError(err)

	destinationsFile, err := os.Open(cfg.DestinationFilePath)
	fatalOnError(err)

	destinationsData, err := readDestinations(destinationsFile)
	fatalOnError(err)

	authConverter := auth.NewConverter()
	frConverter := fetchrequest.NewConverter(authConverter)
	versionConverter := version.NewConverter()
	specConverter := spec.NewConverter(frConverter)
	docConverter := document.NewConverter(frConverter)
	webhookConverter := webhook.NewConverter(authConverter)
	apiConverter := api.NewConverter(versionConverter, specConverter)
	eventAPIConverter := eventdef.NewConverter(versionConverter, specConverter)
	labelDefConverter := labeldef.NewConverter()
	labelConverter := label.NewConverter()
	intSysConverter := integrationsystem.NewConverter()
	bundleConverter := mp_bundle.NewConverter(authConverter, apiConverter, eventAPIConverter, docConverter)
	appConverter := application.NewConverter(webhookConverter, bundleConverter)
	assignmentConv := scenarioassignment.NewConverter()

	runtimeRepo := rt.NewRepository()
	applicationRepo := application.NewRepository(appConverter)
	labelRepo := label.NewRepository(labelConverter)
	labelDefRepo := labeldef.NewRepository(labelDefConverter)
	webhookRepo := webhook.NewRepository(webhookConverter)
	apiRepo := api.NewRepository(apiConverter)
	eventAPIRepo := eventdef.NewRepository(eventAPIConverter)
	docRepo := document.NewRepository(docConverter)
	fetchRequestRepo := fetchrequest.NewRepository(frConverter)
	intSysRepo := integrationsystem.NewRepository(intSysConverter)
	scenarioAssignmentRepo := scenarioassignment.NewRepository(assignmentConv)
	bndlRepo := mp_bundle.NewRepository(bundleConverter)
	uidSvc := uid.NewService()
	labelUpsertSvc := label.NewLabelUpsertService(labelRepo, labelDefRepo, uidSvc)
	scenariosSvc := labeldef.NewScenariosService(labelDefRepo, uidSvc, false)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	scenarioAssignmentEngine := scenarioassignment.NewEngine(labelUpsertSvc, labelRepo, scenarioAssignmentRepo)
	fetchRequestSvc := fetchrequest.NewService(fetchRequestRepo, httpClient)
	specRepo := spec.NewRepository(specConverter)
	specSvc := spec.NewService(specRepo, fetchRequestRepo, uidSvc, fetchRequestSvc)
	apiSvc := api.NewService(apiRepo, uidSvc, specSvc)
	eventSvc := eventdef.NewService(eventAPIRepo, uidSvc, specSvc)
	docSvc := document.NewService(docRepo, fetchRequestRepo, uidSvc)
	bundleSvc := mp_bundle.NewService(bndlRepo, apiSvc, eventSvc, docSvc, uidSvc)

	runtimeSvc := rt.NewService(runtimeRepo, labelRepo, scenariosSvc, labelUpsertSvc, uidSvc, scenarioAssignmentEngine, "")
	normalizer := &normalizer.DefaultNormalizator{}
	appSvc := application.NewService(normalizer, &DummyApplicationHideCfgProvider{}, applicationRepo, webhookRepo, runtimeRepo, labelRepo, intSysRepo, labelUpsertSvc, scenariosSvc, bundleSvc, uidSvc)

	transact, closeFunc, err := persistence.Configure(ctx, cfg.Database)
	fatalOnError(err)

	defer func() {
		err := closeFunc()
		fatalOnError(err)
	}()

	ccConfig := clientcredentials.Config{
		ClientID:     "sb-clonef2ff560788634eb5ba8211a840330d59!b19366|destination-xsappname!b433",
		ClientSecret: "1746ef83-6c6e-4c81-af86-58e4ffd63f77$tmGGI7C6zT3zEV-FHSZ_cQN-CsoUmUDlaCmfrDnMIXo=",
		TokenURL:     "https://graph-crystal-demo.authentication.sap.hana.ondemand.com/oauth/token",
	}
	destHttpClient := ccConfig.Client(ctx)
	destClient := destination.NewClient("https://destination-configuration.cfapps.sap.hana.ondemand.com/destination-configuration/v1", destHttpClient)

	appLabelsHandler := &notifications.AppLabelNotificationHandler{
		RuntimeLister:      runtimeSvc,
		BundleGetter:       bundleSvc,
		AppLister:          appSvc,
		AppLabelGetter:     appSvc,
		RuntimeLabelGetter: runtimeSvc,
		Transact:           transact,
		DestinationCient:   destClient,
		DestinationsData:   destinationsData,
	}

	rtLabelsHandler := &notifications.RuntimeLabelNotificationHandler{
		RuntimeGetter:  runtimeSvc,
		AppLister:      appSvc,
		AppLabelGetter: appSvc,
		Transact:       transact,
	}

	labelsHandler := &notifications.LabelNotificationHandler{
		Handlers: map[resource.Type]notifications.NotificationLabelHandler{
			resource.Application: appLabelsHandler,
			resource.Runtime:     rtLabelsHandler,
		},
	}

	// appHandler := &notifications.AppNotificationHandler{
	// 	ScriptRunner: runner,
	// }

	// rtHandler := &notifications.RtNotificationHandler{
	// 	ScriptRunner: runner,
	// }

	processor := notifications.NewNotificationProcessor(cfg.Database, map[notifications.HandlerKey]notifications.NotificationHandler{
		{
			NotificationChannel: "events",
			ResourceType:        resource.Label,
		}: labelsHandler,
		// {
		// 	NotificationChannel: "events",
		// 	ResourceType:        resource.Application,
		// }: appHandler,
		// {
		// 	NotificationChannel: "events",
		// 	ResourceType:        resource.Runtime,
		// }: rtHandler,
	})

	if err := processor.Run(ctx); err != nil {
		fatalOnError(err)
	}
}

func fatalOnError(err error) {
	if err != nil {
		log.D().Fatal(err.Error())
	}
}

type DummyApplicationHideCfgProvider struct {
}

func (d *DummyApplicationHideCfgProvider) GetApplicationHideSelectors() (map[string][]string, error) {
	return map[string][]string{}, nil
}

func readDestinations(destinationsReader io.Reader) ([]destination.Destination, error) {
	bytes, err := ioutil.ReadAll(destinationsReader)
	if err != nil {
		return nil, err
	}
	var result []destination.Destination
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
