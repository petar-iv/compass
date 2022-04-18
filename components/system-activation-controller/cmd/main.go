package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/kyma-incubator/compass/components/system-activation-controller/controllers"
	"github.com/kyma-incubator/compass/components/system-activation-controller/internal/config"
	"net/http"
	"time"

	"github.tools.sap/unified-resource-manager/api/pkg/apis/logger"
	"github.tools.sap/unified-resource-manager/controller-utils/pkg/manager"
	"go.uber.org/zap/zapcore"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var setupLog = ctrl.Log.WithName("setup")

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	ctx := logger.NewContextWithLogger(context.Background(), initLogger())

	mgr, err := NewManager(ctx)
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Register the SystemActivationReconciler in the controllers manager
	if err := (&controllers.SystemActivationReconciler{
		HttpClient: httpClient,
		Client:     mgr.GetClient(),
		Log:        ctrl.Log.WithName("controllers").WithName("SystemActivation"),
	}).ControllerWithManager(mgr); err != nil {
		setupLog.Error(err, "failed register SystemActivation controller")
		panic(err)
	}

	if err := mgr.Start(ctx); err != nil {
		setupLog.Info(fmt.Sprintf("mgr.Start %+v", err))
		panic(err)
	}
}

func NewManager(ctx context.Context) (*manager.ControllerManager, error) {
	managerSettings, err := config.NewControllerManagerSettings("config")
	if managerSettings == nil || err != nil {
		setupLog.Error(err, "failed to get manager setting")
		panic(err)
	}

	options := &manager.Options{
		ClientConfigFile: managerSettings.APIConfig,
		WebhookPort:      managerSettings.Webhooks.Admission.Port,
		CertDirectory:    managerSettings.Webhooks.X509.Path,
		CertFileName:     managerSettings.Webhooks.X509.Certificate,
		KeyFileName:      managerSettings.Webhooks.X509.Key,
	}

	return manager.NewManager(ctx, options)
}

func initLogger() logger.Logger {
	log := zap.New(zap.UseFlagOptions(&zap.Options{
		Level: zapcore.InfoLevel,
	}))
	return logger.NewLogger(log)
}
