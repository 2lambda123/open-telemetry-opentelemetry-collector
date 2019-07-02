// Copyright 2019, OpenTelemetry Authors
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

package loggingexporter

import (
	"github.com/open-telemetry/opentelemetry-service/consumer"
	"github.com/open-telemetry/opentelemetry-service/factories"
	"github.com/open-telemetry/opentelemetry-service/models"
	"go.uber.org/zap"
)

var _ = factories.RegisterExporterFactory(&exporterFactory{})

const (
	// The value of "type" key in configuration.
	typeStr = "logging"
)

// exporterFactory is the factory for logging exporter.
type exporterFactory struct {
}

// Type gets the type of the Exporter config created by this factory.
func (f *exporterFactory) Type() string {
	return typeStr
}

// CreateDefaultConfig creates the default configuration for exporter.
func (f *exporterFactory) CreateDefaultConfig() models.Exporter {
	return &ConfigV2{
		ExporterSettings: models.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
	}
}

func noopStopFunc() error {
	return nil
}

// CreateTraceExporter creates a trace exporter based on this config.
func (f *exporterFactory) CreateTraceExporter(logger *zap.Logger, config models.Exporter) (consumer.TraceConsumer, factories.StopFunc, error) {

	lexp, err := NewTraceExporter(logger)
	if err != nil {
		return nil, nil, err
	}
	return lexp, noopStopFunc, nil
}

// CreateMetricsExporter creates a metrics exporter based on this config.
func (f *exporterFactory) CreateMetricsExporter(logger *zap.Logger, cfg models.Exporter) (consumer.MetricsConsumer, factories.StopFunc, error) {
	lexp, err := NewMetricsExporter(logger)
	if err != nil {
		return nil, nil, err
	}
	return lexp, noopStopFunc, nil
}
