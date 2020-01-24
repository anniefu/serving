/*
Copyright 2020 The Knative Authors
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

package metrics

import (
	"context"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/resource"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

func newOpenCensusExporter(config *metricsConfig, logger *zap.SugaredLogger) (view.Exporter, error) {
	opts := []ocagent.ExporterOption{ocagent.WithServiceName(config.component)}
	opts = append(opts, ocagent.WithResourceDetector(func(context.Context) (*resource.Resource, error) {
		return &resource.Resource{
			Type: "knative_revision",
			Labels: map[string]string{
				"project_id":         "anniefu-knative-dev",
				"service_name":       "helloworld-go",
				"revision_name":      "helloworld-go-hfc7j",
				"location":           "us-central1-c",
				"configuration_name": "helloworld-go",
				"cluster_name":       "red",
				"namespace_name":     "default",
			},
		}, nil
	}))
	if config.collectorAddress != "" {
		opts = append(opts, ocagent.WithAddress(config.collectorAddress))
	}
	if !config.requireSecure {
		opts = append(opts, ocagent.WithInsecure())
	}
	e, err := ocagent.NewExporter(opts...)
	if err != nil {
		logger.Errorw("Failed to create the OpenCensus exporter.", zap.Error(err))
		return nil, err
	}
	logger.Infof("Created OpenCensus exporter with config: %+v.", *config)
	// Start the server for Prometheus scraping
	view.RegisterExporter(e)
	return e, nil
}
