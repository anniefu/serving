/*
Copyright 2019 The Knative Authors

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
	"log"

	"cloud.google.com/go/compute/metadata"
	"knative.dev/pkg/metrics/metricskey"
)

type gcpMetadata struct {
	project  string
	location string
	cluster  string
}

// Retrieves necessary GCP metadata for sending data to Stackdriver.
// Prioritizes explicitly set values from ConfigMap over values from GCE metadata server.
func retrieveGCPMetadata(config *metricsConfig) *gcpMetadata {
	gm := gcpMetadata{
		project:  metricskey.ValueUnknown,
		location: metricskey.ValueUnknown,
		cluster:  metricskey.ValueUnknown,
	}

	if metadata.OnGCE() {
		project, err := metadata.NumericProjectID()
		if err == nil && project != "" {
			gm.project = project
		}
		location, err := metadata.InstanceAttributeValue("cluster-location")
		if err == nil && location != "" {
			gm.location = location
		}
		cluster, err := metadata.InstanceAttributeValue("cluster-name")
		if err == nil && cluster != "" {
			gm.cluster = cluster
		}
	}

	if config.stackdriverConfig.ProjectID != "" {
		gm.project = config.stackdriverConfig.ProjectID
	}

	if config.stackdriverConfig.ProjectLocation != "" {
		gm.location = config.stackdriverConfig.ProjectLocation
	}

	if config.stackdriverConfig.ClusterName != "" {
		gm.cluster = config.stackdriverConfig.ClusterName
	}

	log.Printf("ANNIE: gcp metadata [%v]", gm)

	return &gm
}
