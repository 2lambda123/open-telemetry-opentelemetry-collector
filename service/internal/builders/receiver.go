// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package builders // import "go.opentelemetry.io/collector/service/internal/builders"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumerprofiles"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverprofiles"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

// ReceiverBuilder receiver is a helper struct that given a set of Configs and
// Factories helps with creating receivers.
type ReceiverBuilder struct {
	cfgs      map[component.ID]component.Config
	factories map[component.Type]receiver.Factory
}

// NewReceiver creates a new ReceiverBuilder to help with creating
// components form a set of configs and factories.
func NewReceiver(cfgs map[component.ID]component.Config, factories map[component.Type]receiver.Factory) *ReceiverBuilder {
	return &ReceiverBuilder{cfgs: cfgs, factories: factories}
}

// CreateTraces creates a Traces receiver based on the settings and config.
func (b *ReceiverBuilder) CreateTraces(ctx context.Context, set receiver.Settings, next consumer.Traces) (receiver.Traces, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("receiver %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("receiver factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.TracesReceiverStability())
	return f.CreateTracesReceiver(ctx, set, cfg, next)
}

// CreateMetrics creates a Metrics receiver based on the settings and config.
func (b *ReceiverBuilder) CreateMetrics(ctx context.Context, set receiver.Settings, next consumer.Metrics) (receiver.Metrics, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("receiver %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("receiver factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.MetricsReceiverStability())
	return f.CreateMetricsReceiver(ctx, set, cfg, next)
}

// CreateLogs creates a Logs receiver based on the settings and config.
func (b *ReceiverBuilder) CreateLogs(ctx context.Context, set receiver.Settings, next consumer.Logs) (receiver.Logs, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("receiver %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("receiver factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.LogsReceiverStability())
	return f.CreateLogsReceiver(ctx, set, cfg, next)
}

// CreateProfiles creates a Profiles receiver based on the settings and config.
func (b *ReceiverBuilder) CreateProfiles(ctx context.Context, set receiver.Settings, next consumerprofiles.Profiles) (receiverprofiles.Profiles, error) {
	if next == nil {
		return nil, errNilNextConsumer
	}
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("receiver %q is not configured", set.ID)
	}

	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("receiver factory not available for: %q", set.ID)
	}

	logStabilityLevel(set.Logger, f.ProfilesReceiverStability())
	return f.CreateProfilesReceiver(ctx, set, cfg, next)
}

func (b *ReceiverBuilder) Factory(componentType component.Type) component.Factory {
	return b.factories[componentType]
}

// NewNopReceiverConfigsAndFactories returns a configuration and factories that allows building a new nop receiver.
func NewNopReceiverConfigsAndFactories() (map[component.ID]component.Config, map[component.Type]receiver.Factory) {
	nopFactory := receivertest.NewNopFactory()
	configs := map[component.ID]component.Config{
		component.NewID(nopType): nopFactory.CreateDefaultConfig(),
	}
	factories := map[component.Type]receiver.Factory{
		nopType: nopFactory,
	}

	return configs, factories
}
