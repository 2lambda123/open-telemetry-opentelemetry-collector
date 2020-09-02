// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resourceprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/internal/data/testdata"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

var (
	processorSettings = configmodels.ProcessorSettings{
		TypeVal: "resource",
		NameVal: "resource",
	}

	cfg = &Config{
		ProcessorSettings: processorSettings,
		AttributesActions: []processorhelper.ActionKeyValue{
			{Key: "cloud.zone", Value: "zone-1", Action: processorhelper.UPSERT},
			{Key: "k8s.cluster.name", FromAttribute: "k8s-cluster", Action: processorhelper.INSERT},
			{Key: "redundant-attribute", Action: processorhelper.DELETE},
		},
	}
)

func TestResourceProcessorAttributesUpsert(t *testing.T) {
	tests := []struct {
		name             string
		config           *Config
		sourceAttributes map[string]string
		wantAttributes   map[string]string
	}{
		{
			name:             "config_with_attributes_applied_on_nil_resource",
			config:           cfg,
			sourceAttributes: nil,
			wantAttributes: map[string]string{
				"cloud.zone": "zone-1",
			},
		},
		{
			name:             "config_with_attributes_applied_on_empty_resource",
			config:           cfg,
			sourceAttributes: map[string]string{},
			wantAttributes: map[string]string{
				"cloud.zone": "zone-1",
			},
		},
		{
			name:   "config_attributes_applied_on_existing_resource_attributes",
			config: cfg,
			sourceAttributes: map[string]string{
				"cloud.zone":          "to-be-replaced",
				"k8s-cluster":         "test-cluster",
				"redundant-attribute": "to-be-removed",
			},
			wantAttributes: map[string]string{
				"cloud.zone":       "zone-1",
				"k8s-cluster":      "test-cluster",
				"k8s.cluster.name": "test-cluster",
			},
		},
		{
			name: "config_attributes_replacement",
			config: &Config{
				ProcessorSettings: processorSettings,
				AttributesActions: []processorhelper.ActionKeyValue{
					{Key: "k8s.cluster.name", FromAttribute: "k8s-cluster", Action: processorhelper.INSERT},
					{Key: "k8s-cluster", Action: processorhelper.DELETE},
				},
			},
			sourceAttributes: map[string]string{
				"k8s-cluster": "test-cluster",
			},
			wantAttributes: map[string]string{
				"k8s.cluster.name": "test-cluster",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test trace consumer
			ttn := &testTraceConsumer{}

			factory := NewFactory()
			rtp, err := factory.CreateTraceProcessor(context.Background(), component.ProcessorCreateParams{}, ttn, tt.config)
			require.NoError(t, err)
			assert.Equal(t, true, rtp.GetCapabilities().MutatesConsumedData)

			sourceTraceData := generateTraceData(tt.sourceAttributes)
			wantTraceData := generateTraceData(tt.wantAttributes)
			err = rtp.ConsumeTraces(context.Background(), sourceTraceData)
			require.NoError(t, err)
			assert.EqualValues(t, wantTraceData, ttn.td)

			// Test metrics consumer
			tmn := &testMetricsConsumer{}
			rmp, err := factory.CreateMetricsProcessor(context.Background(), component.ProcessorCreateParams{}, tmn, tt.config)
			require.NoError(t, err)
			assert.Equal(t, true, rtp.GetCapabilities().MutatesConsumedData)

			sourceMetricData := generateMetricData(tt.sourceAttributes)
			wantMetricData := generateMetricData(tt.wantAttributes)
			err = rmp.ConsumeMetrics(context.Background(), sourceMetricData)
			require.NoError(t, err)
			assert.EqualValues(t, wantMetricData, tmn.md)
		})
	}
}

func generateTraceData(attributes map[string]string) pdata.Traces {
	td := testdata.GenerateTraceDataOneSpanNoResource()
	if attributes == nil {
		return td
	}
	resource := td.ResourceSpans().At(0).Resource()
	resource.InitEmpty()
	for k, v := range attributes {
		resource.Attributes().InsertString(k, v)
	}
	resource.Attributes().Sort()
	return td
}

func generateMetricData(attributes map[string]string) pdata.Metrics {
	md := testdata.GenerateMetricsOneMetricNoResource()
	if attributes == nil {
		return md
	}
	resource := md.ResourceMetrics().At(0).Resource()
	resource.InitEmpty()
	for k, v := range attributes {
		resource.Attributes().InsertString(k, v)
	}
	resource.Attributes().Sort()
	return md
}

type testTraceConsumer struct {
	td pdata.Traces
}

func (ttn *testTraceConsumer) ConsumeTraces(_ context.Context, td pdata.Traces) error {
	// sort attributes to be able to compare traces
	for i := 0; i < td.ResourceSpans().Len(); i++ {
		sortResourceAttributes(td.ResourceSpans().At(i).Resource())
	}
	ttn.td = td
	return nil
}

type testMetricsConsumer struct {
	md pdata.Metrics
}

func (tmn *testMetricsConsumer) ConsumeMetrics(_ context.Context, md pdata.Metrics) error {
	// sort attributes to be able to compare traces
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		sortResourceAttributes(md.ResourceMetrics().At(i).Resource())
	}
	tmn.md = md
	return nil
}

func sortResourceAttributes(resource pdata.Resource) {
	if resource.IsNil() {
		return
	}
	resource.Attributes().Sort()
}
