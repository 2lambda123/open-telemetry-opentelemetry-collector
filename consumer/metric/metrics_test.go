// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestDefaultMetrics(t *testing.T) {
	cp, err := NewMetrics(func(context.Context, pmetric.Metrics) error { return nil })
	assert.NoError(t, err)
	assert.NoError(t, cp.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
	assert.Equal(t, consumer.Capabilities{MutatesData: false}, cp.Capabilities())
}

func TestNilFuncMetrics(t *testing.T) {
	_, err := NewMetrics(nil)
	assert.Equal(t, errNilFunc, err)
}

func TestWithCapabilitiesMetrics(t *testing.T) {
	cp, err := NewMetrics(
		func(context.Context, pmetric.Metrics) error { return nil },
		WithCapabilities(consumer.Capabilities{MutatesData: true}))
	assert.NoError(t, err)
	assert.NoError(t, cp.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
	assert.Equal(t, consumer.Capabilities{MutatesData: true}, cp.Capabilities())
}

func TestConsumeMetrics(t *testing.T) {
	consumeCalled := false
	cp, err := NewMetrics(func(context.Context, pmetric.Metrics) error { consumeCalled = true; return nil })
	assert.NoError(t, err)
	assert.NoError(t, cp.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
	assert.True(t, consumeCalled)
}

func TestConsumeMetrics_ReturnError(t *testing.T) {
	want := errors.New("my_error")
	cp, err := NewMetrics(func(context.Context, pmetric.Metrics) error { return want })
	assert.NoError(t, err)
	assert.Equal(t, want, cp.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
}
