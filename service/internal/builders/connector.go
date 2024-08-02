// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package builders // import "go.opentelemetry.io/collector/service/internal/builders"

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/connector/connectortest"
	"go.opentelemetry.io/collector/consumer"
)

var (
	errNilNextConsumer = errors.New("nil next Consumer")
	nopType            = component.MustNewType("nop")
)

// ConnectorBuilder is a helper struct that given a set of Configs and Factories helps with creating connectors.
type ConnectorBuilder struct {
	cfgs      map[component.ID]component.Config
	factories map[component.Type]connector.Factory
}

// NewConnector creates a new ConnectorBuilder to help with creating components form a set of configs and factories.
func NewConnector(cfgs map[component.ID]component.Config, factories map[component.Type]connector.Factory) *ConnectorBuilder {
	return &ConnectorBuilder{cfgs: cfgs, factories: factories}
}

// CreateTracesToTraces creates a Traces connector based on the settings and config.
func (b *ConnectorBuilder) CreateTracesToTraces(ctx context.Context, set connector.Settings, next consumer.Traces) (connector.Traces, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.TracesToTracesStability())
	return f.CreateTracesToTraces(ctx, set, cfg, next)
}

// CreateTracesToMetrics creates a Traces connector based on the settings and config.
func (b *ConnectorBuilder) CreateTracesToMetrics(ctx context.Context, set connector.Settings, next consumer.Metrics) (connector.Traces, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.TracesToMetricsStability())
	return f.CreateTracesToMetrics(ctx, set, cfg, next)
}

// CreateTracesToLogs creates a Traces connector based on the settings and config.
func (b *ConnectorBuilder) CreateTracesToLogs(ctx context.Context, set connector.Settings, next consumer.Logs) (connector.Traces, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.TracesToLogsStability())
	return f.CreateTracesToLogs(ctx, set, cfg, next)
}

// CreateMetricsToTraces creates a Metrics connector based on the settings and config.
func (b *ConnectorBuilder) CreateMetricsToTraces(ctx context.Context, set connector.Settings, next consumer.Traces) (connector.Metrics, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.MetricsToTracesStability())
	return f.CreateMetricsToTraces(ctx, set, cfg, next)
}

// CreateMetricsToMetrics creates a Metrics connector based on the settings and config.
func (b *ConnectorBuilder) CreateMetricsToMetrics(ctx context.Context, set connector.Settings, next consumer.Metrics) (connector.Metrics, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.MetricsToMetricsStability())
	return f.CreateMetricsToMetrics(ctx, set, cfg, next)
}

// CreateMetricsToLogs creates a Metrics connector based on the settings and config.
func (b *ConnectorBuilder) CreateMetricsToLogs(ctx context.Context, set connector.Settings, next consumer.Logs) (connector.Metrics, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.MetricsToLogsStability())
	return f.CreateMetricsToLogs(ctx, set, cfg, next)
}

// CreateLogsToTraces creates a Logs connector based on the settings and config.
func (b *ConnectorBuilder) CreateLogsToTraces(ctx context.Context, set connector.Settings, next consumer.Traces) (connector.Logs, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.LogsToTracesStability())
	return f.CreateLogsToTraces(ctx, set, cfg, next)
}

// CreateLogsToMetrics creates a Logs connector based on the settings and config.
func (b *ConnectorBuilder) CreateLogsToMetrics(ctx context.Context, set connector.Settings, next consumer.Metrics) (connector.Logs, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.LogsToMetricsStability())
	return f.CreateLogsToMetrics(ctx, set, cfg, next)
}

// CreateLogsToLogs creates a Logs connector based on the settings and config.
func (b *ConnectorBuilder) CreateLogsToLogs(ctx context.Context, set connector.Settings, next consumer.Logs) (connector.Logs, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("connector %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("connector factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.LogsToLogsStability())
	return f.CreateLogsToLogs(ctx, set, cfg, next)
}

func (b *ConnectorBuilder) IsConfigured(componentID component.ID) bool {
	_, ok := b.cfgs[componentID]
	return ok
}

func (b *ConnectorBuilder) Factory(componentType component.Type) component.Factory {
	return b.factories[componentType]
}

// NewNopConnectorConfigsAndFactories returns a configuration and factories that allows building a new nop connector.
func NewNopConnectorConfigsAndFactories() (map[component.ID]component.Config, map[component.Type]connector.Factory) {
	nopFactory := connectortest.NewNopFactory()
	// Use a different ID than receivertest and exportertest to avoid ambiguous
	// configuration scenarios. Ambiguous IDs are detected in the 'otelcol' package,
	// but lower level packages such as 'service' assume that IDs are disambiguated.
	connID := component.NewIDWithName(nopType, "conn")

	configs := map[component.ID]component.Config{
		connID: nopFactory.CreateDefaultConfig(),
	}
	factories := map[component.Type]connector.Factory{
		nopType: nopFactory,
	}

	return configs, factories
}

// logStabilityLevel logs the stability level of a component. The log level is set to info for
// undefined, unmaintained, deprecated and development. The log level is set to debug
// for alpha, beta and stable.
func logStabilityLevel(logger *zap.Logger, sl component.StabilityLevel) {
	if sl >= component.StabilityLevelAlpha {
		logger.Debug(sl.LogMessage())
	} else {
		logger.Info(sl.LogMessage())
	}
}
