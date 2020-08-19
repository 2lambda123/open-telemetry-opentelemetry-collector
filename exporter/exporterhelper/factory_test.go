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

package exporterhelper

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exportertest"
)

const typeStr = "test"

var (
	defaultCfg = &configmodels.ExporterSettings{
		TypeVal: typeStr,
		NameVal: typeStr,
	}
	nopTracesExporter  = exportertest.NewNopTraceExporter()
	nopMetricsExporter = exportertest.NewNopMetricsExporter()
	nopLogsExporter    = exportertest.NewNopLogsExporter()
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory(
		typeStr,
		defaultConfig)
	assert.EqualValues(t, typeStr, factory.Type())
	assert.EqualValues(t, defaultCfg, factory.CreateDefaultConfig())
	_, ok := factory.(component.ConfigUnmarshaler)
	assert.False(t, ok)
	_, err := factory.CreateTraceExporter(context.Background(), component.ExporterCreateParams{}, defaultCfg)
	assert.Equal(t, configerror.ErrDataTypeIsNotSupported, err)
	_, err = factory.CreateMetricsExporter(context.Background(), component.ExporterCreateParams{}, defaultCfg)
	assert.Equal(t, configerror.ErrDataTypeIsNotSupported, err)
	lfactory := factory.(component.LogsExporterFactory)
	_, err = lfactory.CreateLogsExporter(context.Background(), component.ExporterCreateParams{}, defaultCfg)
	assert.Equal(t, configerror.ErrDataTypeIsNotSupported, err)
}

func TestNewFactory_WithConstructors(t *testing.T) {
	factory := NewFactory(
		typeStr,
		defaultConfig,
		WithTraces(createTraceExporter),
		WithMetrics(createMetricsExporter),
		WithLogs(createLogsExporter),
		WithCustomUnmarshaler(customUnmarshaler))
	assert.EqualValues(t, typeStr, factory.Type())
	assert.EqualValues(t, defaultCfg, factory.CreateDefaultConfig())

	fu, ok := factory.(component.ConfigUnmarshaler)
	assert.True(t, ok)
	assert.Equal(t, errors.New("my error"), fu.Unmarshal(nil, nil))

	te, err := factory.CreateTraceExporter(context.Background(), component.ExporterCreateParams{}, defaultCfg)
	assert.NoError(t, err)
	assert.Same(t, nopTracesExporter, te)

	me, err := factory.CreateMetricsExporter(context.Background(), component.ExporterCreateParams{}, defaultCfg)
	assert.NoError(t, err)
	assert.Same(t, nopMetricsExporter, me)

	lfactory := factory.(component.LogsExporterFactory)
	le, err := lfactory.CreateLogsExporter(context.Background(), component.ExporterCreateParams{}, defaultCfg)
	assert.NoError(t, err)
	assert.Same(t, nopLogsExporter, le)
}

func defaultConfig() configmodels.Exporter {
	return defaultCfg
}

func createTraceExporter(context.Context, component.ExporterCreateParams, configmodels.Exporter) (component.TraceExporter, error) {
	return nopTracesExporter, nil
}

func createMetricsExporter(context.Context, component.ExporterCreateParams, configmodels.Exporter) (component.MetricsExporter, error) {
	return nopMetricsExporter, nil
}

func createLogsExporter(context.Context, component.ExporterCreateParams, configmodels.Exporter) (component.LogsExporter, error) {
	return nopLogsExporter, nil
}

func customUnmarshaler(*viper.Viper, interface{}) error {
	return errors.New("my error")
}
