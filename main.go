/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"time"

	"nais/replicator/internal/logger"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	naisiov1 "nais/replicator/api/v1"
	"nais/replicator/controllers"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(naisiov1.AddToScheme(scheme))
	logger.SetupLogrus()
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var enableWebhooks bool
	var debug bool
	var interval time.Duration
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&enableWebhooks, "enable-webhooks", true, "Enable webhooks")
	flag.BoolVar(&debug, "debug", os.Getenv("DEBUG") == "true", "Enable debug logging")
	flag.DurationVar(&interval, "sync-interval", 15*time.Minute, "Synchronization interval for reconciliation")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if os.Getenv("POD_NAMESPACE") == "" {
		log.Error("POD_NAMESPACE environment variable must be set")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "9046ff70.nais.io",
		Cache: cache.Options{
			SyncPeriod: &interval,
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: 9443,
		}),
	})
	if err != nil {
		log.Errorf("unable to start manager %v", err)
		os.Exit(1)
	}

	if err = (&controllers.ReplicationConfigReconciler{
		Client:       mgr.GetClient(),
		Scheme:       mgr.GetScheme(),
		Recorder:     mgr.GetEventRecorderFor("replicator"),
		SyncInterval: interval,
	}).SetupWithManager(mgr); err != nil {
		log.Errorf("unable to create controller %v", err)
		os.Exit(1)
	}

	if enableWebhooks {
		log.Infof("webhooks enabled, registering webhook server at /validate-replicationconfig")
		ctrl := controllers.NewReplicatorValidator(mgr)
		mgr.GetWebhookServer().Register("/validate-replicationconfig", &webhook.Admission{Handler: ctrl})
	}

	//+kubebuilder:scaffold:builder
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Errorf("unable to set up health check %v", err)
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Errorf("unable to set up ready check %v", err)
		os.Exit(1)
	}

	log.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Errorf("problem running manager %v", err)
		os.Exit(1)
	}
}
