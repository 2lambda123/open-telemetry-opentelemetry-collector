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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
)

const (
	// The value of "type" key in configuration.
	typeStr = "logging"
)

// Factory is the factory for logging exporter.
type Factory struct {
}

// Type gets the type of the Exporter config created by this factory.
func (f *Factory) Type() string {
	return typeStr
}

// CreateDefaultConfig creates the default configuration for exporter.
func (f *Factory) CreateDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		LogLevel:           "info",
		SamplingInitial:    5,
		SamplingThereafter: 100,
	}
}

// CreateTraceExporter creates a trace exporter based on this config.
func (f *Factory) CreateTraceExporter(logger *zap.Logger, config configmodels.Exporter) (exporter.TraceExporter, error) {
	cfg := config.(*Config)

	exporterLogger, err := f.createLogger(cfg.LogLevel, cfg.SamplingInitial, cfg.SamplingThereafter)
	if err != nil {
		return nil, err
	}

	lexp, err := NewTraceExporter(config, cfg.LogLevel, exporterLogger)
	if err != nil {
		return nil, err
	}
	return lexp, nil
}

func (f *Factory) createLogger(logLevel string, initial int, thereafter int) (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}

	encoderConf := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		MessageKey:     "M",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	conf := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    initial,
			Thereafter: thereafter,
		},
		Encoding:         "console",
		EncoderConfig:    encoderConf,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logginglogger, err := conf.Build()
	if err != nil {
		return nil, err
	}
	return logginglogger, nil
}

// CreateMetricsExporter creates a metrics exporter based on this config.
func (f *Factory) CreateMetricsExporter(logger *zap.Logger, config configmodels.Exporter) (exporter.MetricsExporter, error) {
	cfg := config.(*Config)

	exporterLogger, err := f.createLogger(cfg.LogLevel, cfg.SamplingInitial, cfg.SamplingThereafter)
	if err != nil {
		return nil, err
	}

	lexp, err := NewMetricsExporter(config, exporterLogger)
	if err != nil {
		return nil, err
	}
	return lexp, nil
}
