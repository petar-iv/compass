package main

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/pkg/correlation"
	"github.com/kyma-incubator/compass/components/director/pkg/kubernetes"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/namespacedname"
	"github.com/kyma-incubator/compass/components/director/pkg/signal"
	"github.com/vrischmann/envconfig"
	"os"
	"sync"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	//TODO make its own Config struct for the new component
	cfg := Config{}
	err := envconfig.InitWithPrefix(&cfg, "APP")
	exitOnError(err, "Error while loading app Config")

	ctx, err = log.Configure(ctx, &cfg.Log)
	exitOnError(err, "Filed to configure logger")

	log.C(ctx).Info("Starting Hydrators Service")
	log.C(ctx).Infof("Config: %+v", cfg)

	//TODO this is picked up from director utils now, its already there
	k8sClientSet, appErr := kubernetes.NewKubernetesClientSet(ctx, cfg.KubernetesClient.PollInteval, cfg.KubernetesClient.PollTimeout, cfg.KubernetesClient.Timeout)
	exitOnError(appErr, "Failed to initialize Kubernetes client.")

	revokedCertsCache := NewCache()

	//TODO this is picked up from director utils now, its already there
	revokedCertsConfigMap, err := namespacedname.Parse(cfg.RevocationConfigMapName)
	exitOnError(err, "Failed to initialize revokedCertsConfigMap.")

	revokedCertsLoader := NewRevokedCertificatesLoader(revokedCertsCache,
		k8sClientSet.CoreV1().ConfigMaps(revokedCertsConfigMap.Namespace),
		revokedCertsConfigMap.Name,
		time.Second,
	)

	hydratorServer, err := PrepareHydratorServer(cfg, revokedCertsCache, correlation.AttachCorrelationIDToContext(), log.RequestLogger())
	exitOnError(err, "Failed configuring hydrator handler")

	//TODO rm wg
	wg := &sync.WaitGroup{}
	go revokedCertsLoader.Run(ctx)
	go startServer(ctx, hydratorServer, wg)

}
