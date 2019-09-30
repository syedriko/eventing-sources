/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"io/ioutil"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"knative.dev/eventing-contrib/prometheus/pkg/adapter"
	"knative.dev/eventing/pkg/kncloudevents"
	"knative.dev/pkg/signals"
)

type envConfig struct {
	Namespace              string `envconfig:"SYSTEM_NAMESPACE" default:"default"`
	EventSource            string `envconfig:"EVENT_SOURCE" required:"true"`
	SinkURI                string `split_words:"true" required:"true"`
// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
}

func main() {
	flag.Parse()

	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	dlogger, err := logCfg.Build()
	logger := dlogger.Sugar()

	var env envConfig
	err = envconfig.Process("", &env)
	if err != nil {
		logger.Fatalw("Error processing environment", zap.Error(err))
	}

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	logger = logger.With(zap.String("controller/prometheus_source", "adapter"))
	logger.Info("Starting the adapter")

	eventsClient, err := kncloudevents.NewDefaultClient(env.SinkURI)
	if err != nil {
		logger.Fatalw("Error building cloud event client", zap.Error(err))
	}

	opt := adapter.Options{
		Namespace:   env.Namespace,
		EventSource: env.EventSource,
	}

	a, err := adapter.New(eventsClient, logger, &opt)
	if err != nil {
		logger.Fatalf("Failed to create prometheus source adapter: %s", err.Error())
	}

	logger.Info("starting Prometheus source adapter")
	if err := a.Start(stopCh); err != nil {
		logger.Warn("start returned an error,", zap.Error(err))
	}
}
