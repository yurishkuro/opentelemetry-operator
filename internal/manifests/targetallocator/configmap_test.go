// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package targetallocator

import (
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-telemetry/opentelemetry-operator/internal/config"
	"github.com/open-telemetry/opentelemetry-operator/internal/manifests"
)

func TestDesiredConfigMap(t *testing.T) {
	expectedLables := map[string]string{
		"app.kubernetes.io/managed-by": "opentelemetry-operator",
		"app.kubernetes.io/instance":   "default.my-instance",
		"app.kubernetes.io/part-of":    "opentelemetry",
		"app.kubernetes.io/version":    "0.47.0",
	}

	t.Run("should return expected target allocator config map", func(t *testing.T) {
		expectedLables["app.kubernetes.io/component"] = "opentelemetry-targetallocator"
		expectedLables["app.kubernetes.io/name"] = "my-instance-targetallocator"

		expectedData := map[string]string{
			"targetallocator.yaml": `allocation_strategy: consistent-hashing
collector_selector:
  matchlabels:
    app.kubernetes.io/component: opentelemetry-collector
    app.kubernetes.io/instance: default.my-instance
    app.kubernetes.io/managed-by: opentelemetry-operator
    app.kubernetes.io/part-of: opentelemetry
config:
  scrape_configs:
  - job_name: otel-collector
    scrape_interval: 10s
    static_configs:
    - targets:
      - 0.0.0.0:8888
      - 0.0.0.0:9999
filter_strategy: relabel-config
label_selector:
  app.kubernetes.io/component: opentelemetry-collector
  app.kubernetes.io/instance: default.my-instance
  app.kubernetes.io/managed-by: opentelemetry-operator
  app.kubernetes.io/part-of: opentelemetry
prometheus_cr:
  pod_monitor_selector: null
  service_monitor_selector: null
`,
		}
		instance := collectorInstance()
		cfg := config.New()
		params := manifests.Params{
			OtelCol: instance,
			Config:  cfg,
			Log:     logr.Discard(),
		}
		actual, err := ConfigMap(params)
		assert.NoError(t, err)

		assert.Equal(t, "my-instance-targetallocator", actual.Name)
		assert.Equal(t, expectedLables, actual.Labels)
		assert.Equal(t, expectedData, actual.Data)

	})
	t.Run("should return expected target allocator config map with label selectors", func(t *testing.T) {
		expectedLables["app.kubernetes.io/component"] = "opentelemetry-targetallocator"
		expectedLables["app.kubernetes.io/name"] = "my-instance-targetallocator"

		expectedData := map[string]string{
			"targetallocator.yaml": `allocation_strategy: consistent-hashing
collector_selector:
  matchlabels:
    app.kubernetes.io/component: opentelemetry-collector
    app.kubernetes.io/instance: default.my-instance
    app.kubernetes.io/managed-by: opentelemetry-operator
    app.kubernetes.io/part-of: opentelemetry
config:
  scrape_configs:
  - job_name: otel-collector
    scrape_interval: 10s
    static_configs:
    - targets:
      - 0.0.0.0:8888
      - 0.0.0.0:9999
filter_strategy: relabel-config
label_selector:
  app.kubernetes.io/component: opentelemetry-collector
  app.kubernetes.io/instance: default.my-instance
  app.kubernetes.io/managed-by: opentelemetry-operator
  app.kubernetes.io/part-of: opentelemetry
pod_monitor_selector:
  release: my-instance
prometheus_cr:
  pod_monitor_selector:
    matchlabels:
      release: my-instance
    matchexpressions: []
  service_monitor_selector:
    matchlabels:
      release: my-instance
    matchexpressions: []
service_monitor_selector:
  release: my-instance
`,
		}
		instance := collectorInstance()
		instance.Spec.TargetAllocator.PrometheusCR.PodMonitorSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"release": "my-instance",
			},
		}
		instance.Spec.TargetAllocator.PrometheusCR.ServiceMonitorSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"release": "my-instance",
			},
		}
		cfg := config.New()
		params := manifests.Params{
			OtelCol: instance,
			Config:  cfg,
			Log:     logr.Discard(),
		}
		actual, err := ConfigMap(params)
		assert.NoError(t, err)

		assert.Equal(t, "my-instance-targetallocator", actual.Name)
		assert.Equal(t, expectedLables, actual.Labels)
		assert.Equal(t, expectedData, actual.Data)

	})
	t.Run("should return expected target allocator config map with scrape interval set", func(t *testing.T) {
		expectedLables["app.kubernetes.io/component"] = "opentelemetry-targetallocator"
		expectedLables["app.kubernetes.io/name"] = "my-instance-targetallocator"

		expectedData := map[string]string{
			"targetallocator.yaml": `allocation_strategy: consistent-hashing
collector_selector:
  matchlabels:
    app.kubernetes.io/component: opentelemetry-collector
    app.kubernetes.io/instance: default.my-instance
    app.kubernetes.io/managed-by: opentelemetry-operator
    app.kubernetes.io/part-of: opentelemetry
config:
  scrape_configs:
  - job_name: otel-collector
    scrape_interval: 10s
    static_configs:
    - targets:
      - 0.0.0.0:8888
      - 0.0.0.0:9999
filter_strategy: relabel-config
label_selector:
  app.kubernetes.io/component: opentelemetry-collector
  app.kubernetes.io/instance: default.my-instance
  app.kubernetes.io/managed-by: opentelemetry-operator
  app.kubernetes.io/part-of: opentelemetry
prometheus_cr:
  pod_monitor_selector: null
  scrape_interval: 30s
  service_monitor_selector: null
`,
		}

		collector := collectorInstance()
		collector.Spec.TargetAllocator.PrometheusCR.ScrapeInterval = &metav1.Duration{Duration: time.Second * 30}
		cfg := config.New()
		params := manifests.Params{
			OtelCol: collector,
			Config:  cfg,
			Log:     logr.Discard(),
		}
		actual, err := ConfigMap(params)
		assert.NoError(t, err)

		assert.Equal(t, "my-instance-targetallocator", actual.Name)
		assert.Equal(t, expectedLables, actual.Labels)
		assert.Equal(t, expectedData, actual.Data)

	})

}
