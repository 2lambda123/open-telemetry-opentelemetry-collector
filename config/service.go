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

package config // import "go.opentelemetry.io/collector/config"

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.uber.org/zap/zapcore"
)

// Service defines the configurable components of the service.
type Service struct {
	// Telemetry is the configuration for collector's own telemetry.
	Telemetry ServiceTelemetry `mapstructure:"telemetry"`

	// Extensions are the ordered list of extensions configured for the service.
	Extensions []ComponentID `mapstructure:"extensions"`

	// Pipelines are the set of data pipelines configured for the service.
	Pipelines Pipelines `mapstructure:"pipelines"`
}

// ServiceTelemetry defines the configurable settings for service telemetry.
type ServiceTelemetry struct {
	Logs    ServiceTelemetryLogs    `mapstructure:"logs"`
	Metrics ServiceTelemetryMetrics `mapstructure:"metrics"`
}

func (srvT *ServiceTelemetry) validate() error {
	if err := srvT.Logs.validate(); err != nil {
		return err
	}

	return srvT.Metrics.validate()
}

// ServiceTelemetryLogs defines the configurable settings for service telemetry logs.
// This MUST be compatible with zap.Config. Cannot use directly zap.Config because
// the collector uses mapstructure and not yaml tags.
type ServiceTelemetryLogs struct {
	// Level is the minimum enabled logging level.
	Level zapcore.Level `mapstructure:"level"`

	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `mapstructure:"development"`

	// Encoding sets the logger's encoding.
	// Valid values are "json" and "console".
	Encoding string `mapstructure:"encoding"`
}

func (srvTL *ServiceTelemetryLogs) validate() error {
	if srvTL.Encoding != "json" && srvTL.Encoding != "console" {
		return fmt.Errorf(`service telemetry logs invalid encoding: %q, valid values are "json" and "console"`, srvTL.Encoding)
	}
	return nil
}

// ServiceTelemetryMetrics exposes the common Telemetry configuration for one component.
type ServiceTelemetryMetrics struct {
	// Level is the level of telemetry metrics, the possible values are:
	//  - "none" indicates that no telemetry data should be collected;
	//  - "basic" is the recommended and covers the basics of the service telemetry.
	//  - "normal" adds some other indicators on top of basic.
	//  - "detailed" adds dimensions and views to the previous levels.
	Level configtelemetry.Level `mapstructure:"level"`

	// Address is the [address]:port that metrics exposition should be bound to.
	Address string `mapstructure:"address"`
}

func (s ServiceTelemetryMetrics) validate() error {
	return nil
}

// DataType is a special Type that represents the data types supported by the collector. We currently support
// collecting metrics, traces and logs, this can expand in the future.
type DataType = Type

// Currently supported data types. Add new data types here when new types are supported in the future.
const (
	// TracesDataType is the data type tag for traces.
	TracesDataType DataType = "traces"

	// MetricsDataType is the data type tag for metrics.
	MetricsDataType DataType = "metrics"

	// LogsDataType is the data type tag for logs.
	LogsDataType DataType = "logs"
)

// Pipeline defines a single pipeline.
type Pipeline struct {
	Receivers  []ComponentID `mapstructure:"receivers"`
	Processors []ComponentID `mapstructure:"processors"`
	Exporters  []ComponentID `mapstructure:"exporters"`
}

// Pipelines is a map of names to Pipelines.
type Pipelines map[ComponentID]*Pipeline
