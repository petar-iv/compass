package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strings"

	"github.com/kyma-incubator/compass/components/atom-tenants-controller/pkg/fetcher_client"

	"github.com/kyma-incubator/compass/components/atom-tenants-controller/internal/config"
	"github.com/kyma-incubator/compass/components/atom-tenants-controller/reconcilers"
	"github.com/kyma-incubator/compass/components/director/pkg/signal"
	rmlogger "github.tools.sap/unified-resource-manager/api/pkg/apis/logger"
	"github.tools.sap/unified-resource-manager/controller-utils/pkg/manager"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2/clientcredentials"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	settings, err := config.NewControllerManagerSettings()
	if err != nil {
		initLogger(nil).Error(err, "failed to initialize controller manager settings")
		os.Exit(1)
	}

	ctx = rmlogger.NewContextWithLogger(ctx, initLogger(settings))
	log := rmlogger.FromContext(ctx, "tenants-aggregator")

	options := &manager.Options{
		ClientConfigFile: settings.RMConfig,
	}

	oauth20Config := clientcredentials.Config{
		ClientID:     settings.OAuth20.ClientKey,
		ClientSecret: settings.OAuth20.ClientSecret,
		TokenURL:     settings.OAuth20.TokenURL,
	}

	oauth20Client := oauth20Config.Client(ctx)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: settings.SkipSSLValidation,
		},
	}
	oauth20Client.Transport = fetcher_client.NewOAuth20Transport(transport, oauth20Config)
	creator := fetcher_client.NewCreator(oauth20Client, settings.TenantFetcherURL)

	mgr, err := manager.NewManager(ctx, options)
	exitOnError(log, err, "failed to create manager")

	err = registerControllers(ctx, mgr, creator, settings)
	exitOnError(log, err, "failed to register reconcilers")

	err = mgr.Start(ctx)
	exitOnError(log, err, "failed to start manager")
}

func registerControllers(ctx context.Context, mgr *manager.ControllerManager, creator fetcher_client.TenantCreator, settings *config.ControllerManagerSettings) error {
	var err error
	log := rmlogger.FromContext(ctx, "ControllerManager")
	rmClient := mgr.GetClient()

	if err = (&reconcilers.OrganizationController{
		Client:  rmClient,
		Creator: creator,
		Log:     rmlogger.FromContext(ctx, "OrganizationController"),
	}).ControllerWithManager(mgr, settings.MaxConcurrentThreads.OrganizationReconcilerThreads); err != nil {
		log.Error(err, "failed to create Organization controller")
		return err
	}

	if err = (&reconcilers.FolderController{
		Client:  rmClient,
		Creator: creator,
		Log:     rmlogger.FromContext(ctx, "FolderReconciler"),
	}).ControllerWithManager(mgr, settings.MaxConcurrentThreads.FolderReconcilerThreads); err != nil {
		log.Error(err, "failed to create Folder controller")
		return err
	}

	if err = (&reconcilers.ResourceGroupController{
		Client:  rmClient,
		Creator: creator,
		Log:     rmlogger.FromContext(ctx, "ResourceGroupReconciler"),
	}).ControllerWithManager(mgr, settings.MaxConcurrentThreads.ResourceGroupReconcilerThreads); err != nil {
		log.Error(err, "failed to create ResourceGroup controller")
		return err
	}

	return nil
}

func initLogger(settings *config.ControllerManagerSettings) rmlogger.Logger {
	var logLevel zapcore.LevelEnabler
	if settings == nil {
		logLevel = zapcore.InfoLevel
	} else {
		logLevel = getLogLevel(settings.Log)
	}
	log := zap.New(zap.UseFlagOptions(&zap.Options{
		Level: logLevel,
	}))
	return rmlogger.NewLogger(log)
}

func getLogLevel(logSettings config.LogSettings) zapcore.LevelEnabler {
	if !logSettings.Enabled {
		return zapcore.PanicLevel
	}

	level := strings.ToLower(logSettings.Level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.ErrorLevel
	}
}

func exitOnError(log rmlogger.Logger, err error, msg string) {
	if err != nil {
		log.Error(err, msg)
		os.Exit(1)
	}
}
