// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package processortest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestNewNopFactory(t *testing.T) {
	factory := NewNopFactory()
	require.NotNil(t, factory)
	assert.Equal(t, component.MustNewType("nop"), factory.Type())
	cfg := factory.CreateDefaultConfig()
	assert.Equal(t, &nopConfig{}, cfg)

	traces, err := factory.CreateTracesProcessor(context.Background(), NewNopSettings(), cfg, consumertest.NewNop())
	require.NoError(t, err)
	assert.Equal(t, consumer.Capabilities{MutatesData: false}, traces.Capabilities())
	assert.NoError(t, traces.Start(context.Background(), componenttest.NewNopHost()))
	assert.NoError(t, traces.ConsumeTraces(context.Background(), ptrace.NewTraces()))
	assert.NoError(t, traces.Shutdown(context.Background()))

	metrics, err := factory.CreateMetricsProcessor(context.Background(), NewNopSettings(), cfg, consumertest.NewNop())
	require.NoError(t, err)
	assert.Equal(t, consumer.Capabilities{MutatesData: false}, metrics.Capabilities())
	assert.NoError(t, metrics.Start(context.Background(), componenttest.NewNopHost()))
	assert.NoError(t, metrics.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
	assert.NoError(t, metrics.Shutdown(context.Background()))

	logs, err := factory.CreateLogsProcessor(context.Background(), NewNopSettings(), cfg, consumertest.NewNop())
	require.NoError(t, err)
	assert.Equal(t, consumer.Capabilities{MutatesData: false}, logs.Capabilities())
	assert.NoError(t, logs.Start(context.Background(), componenttest.NewNopHost()))
	assert.NoError(t, logs.ConsumeLogs(context.Background(), plog.NewLogs()))
	assert.NoError(t, logs.Shutdown(context.Background()))
}
