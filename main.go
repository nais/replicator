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
	"strconv"
	"time"

	"nais/replicator/internal/logger"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
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
	var syncInterval string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&enableWebhooks, "enable-webhooks", true, "Enable webhooks")
	flag.BoolVar(&debug, "debug", os.Getenv("DEBUG") == "true", "Enable debug logging")
	flag.StringVar(&syncInterval, "sync-interval", os.Getenv("SYNC_INTERVAL_MINUTES"), "Synchronization interval for reconciliation in minutes")

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

	interval := 15 * time.Minute
	if syncInterval != "" {
		syncIntervalInt, err := strconv.Atoi(syncInterval)
		if err != nil {
			log.Errorf("unable to convert sync interval to number %v", err)
			os.Exit(1)
		}
		interval = time.Duration(syncIntervalInt) * time.Minute
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "9046ff70.nais.io",
		SyncPeriod:             &interval,
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
		mgr.GetWebhookServer().Register("/validate-replicationconfig", &webhook.Admission{Handler: &controllers.ReplicatorValidator{Client: mgr.GetClient()}})
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
