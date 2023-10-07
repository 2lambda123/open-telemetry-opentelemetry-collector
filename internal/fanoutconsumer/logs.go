// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package fanoutconsumer contains implementations of Traces/Metrics/Logs consumers
// that fan out the data to multiple other consumers.
package fanoutconsumer // import "go.opentelemetry.io/collector/internal/fanoutconsumer"

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
)

// NewLogs wraps multiple log consumers in a single one.
// It fanouts the incoming data to all the consumers, and does smart routing:
//   - Clones only to the consumer that needs to mutate the data.
//   - If all consumers needs to mutate the data one will get the original mutable data.
func NewLogs(lcs []consumer.Logs) consumer.Logs {
	lc := &logsConsumer{}
	for i := 0; i < len(lcs); i++ {
		if lcs[i].Capabilities().MutatesData {
			lc.mutable = append(lc.mutable, lcs[i])
		} else {
			lc.readonly = append(lc.readonly, lcs[i])
		}
	}
	return lc
}

type logsConsumer struct {
	mutable  []consumer.Logs
	readonly []consumer.Logs
}

func (lsc *logsConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeLogs exports the plog.Logs to all consumers wrapped by the current one.
func (lsc *logsConsumer) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	var errs error

	// Clone the data before sending to mutable consumers.
	// The only exception is the last consumer which is allowed to mutate the data only if there are no
	// other non-mutating consumers and the data is mutable. Never share the same data between
	// a mutating and a non-mutating consumer since the non-mutating consumer may process
	// data async and the mutating consumer may change the data before that.
	for i, lc := range lsc.mutable {
		if i < len(lsc.mutable)-1 || ld.IsReadOnly() || len(lsc.readonly) > 0 {
			clonedLogs := plog.NewLogs()
			ld.CopyTo(clonedLogs)
			errs = multierr.Append(errs, lc.ConsumeLogs(ctx, clonedLogs))
		} else {
			errs = multierr.Append(errs, lc.ConsumeLogs(ctx, ld))
		}
	}

	// Mark the data as read-only if it will be sent to more than one read-only consumer.
	if len(lsc.readonly) > 1 && !ld.IsReadOnly() {
		ld.MarkReadOnly()
	}
	for _, lc := range lsc.readonly {
		errs = multierr.Append(errs, lc.ConsumeLogs(ctx, ld))
	}

	return errs
}

var _ connector.LogsRouter = (*logsRouter)(nil)

type logsRouter struct {
	consumer.Logs
	consumers map[component.ID]consumer.Logs
}

func NewLogsRouter(cm map[component.ID]consumer.Logs) consumer.Logs {
	consumers := make([]consumer.Logs, 0, len(cm))
	for _, consumer := range cm {
		consumers = append(consumers, consumer)
	}
	return &logsRouter{
		Logs:      NewLogs(consumers),
		consumers: cm,
	}
}

func (r *logsRouter) PipelineIDs() []component.ID {
	ids := make([]component.ID, 0, len(r.consumers))
	for id := range r.consumers {
		ids = append(ids, id)
	}
	return ids
}

func (r *logsRouter) Consumer(pipelineIDs ...component.ID) (consumer.Logs, error) {
	if len(pipelineIDs) == 0 {
		return nil, fmt.Errorf("missing consumers")
	}
	consumers := make([]consumer.Logs, 0, len(pipelineIDs))
	var errors error
	for _, pipelineID := range pipelineIDs {
		c, ok := r.consumers[pipelineID]
		if ok {
			consumers = append(consumers, c)
		} else {
			errors = multierr.Append(errors, fmt.Errorf("missing consumer: %q", pipelineID))
		}
	}
	if errors != nil {
		// TODO potentially this could return a NewLogs with the valid consumers
		return nil, errors
	}
	return NewLogs(consumers), nil
}
