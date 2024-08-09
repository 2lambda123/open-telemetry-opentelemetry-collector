// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package consumer // import "go.opentelemetry.io/collector/consumer"

import (
	"context"

	"go.opentelemetry.io/collector/consumer/internal"
	"go.opentelemetry.io/collector/pdata/plog"
)

// Logs is an interface that receives plog.Logs, processes it
// as needed, and sends it to the next processing node if any or to the destination.
type Logs interface {
	internal.BaseConsumer
	// ConsumeLogs receives plog.Logs for consumption.
	ConsumeLogs(ctx context.Context, ld plog.Logs) error
}

// ConsumeLogsFunc is a helper function that is similar to ConsumeLogs.
type ConsumeLogsFunc func(ctx context.Context, ld plog.Logs) error

// ConsumeLogs calls f(ctx, ld).
func (f ConsumeLogsFunc) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	return f(ctx, ld)
}

type baseLogs struct {
	*internal.BaseImpl
	ConsumeLogsFunc
}

// NewLogs returns a Logs configured with the provided options.
func NewLogs(consume ConsumeLogsFunc, options ...Option) (Logs, error) {
	if consume == nil {
		return nil, errNilFunc
	}
	baseImpl := internal.NewBaseImpl(options...)
	fn := func(ctx context.Context, ld plog.Logs) error {
		baseImpl.ObsReport.StartTracesOp(ctx)
		logRecordCount := ld.LogRecordCount()
		err := consume(ctx, ld)
		baseImpl.ObsReport.EndTracesOp(ctx, logRecordCount, err)
		return err
	}

	return &baseLogs{
		BaseImpl:        baseImpl,
		ConsumeLogsFunc: fn,
	}, nil
}
