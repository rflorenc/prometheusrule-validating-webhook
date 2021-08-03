package main

import (
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	webhookadmission "github.com/rflorenc/prometheusrule-validating-webhook/admission"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	setupLog := log.Log.WithName("entrypoint")

	// setup a manager
	setupLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		setupLog.Error(err, "unable to setup controller manager")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	setupLog.Info("registering PrometheusRule validating webhook endpoint")
	hookServer.Register("/validate-v1-prometheusrule", &webhook.Admission{Handler: &webhookadmission.PrometheusRuleValidator{Client: mgr.GetClient()}})

	setupLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
