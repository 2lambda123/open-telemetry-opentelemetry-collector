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

package fanoutconsumer // import "go.opentelemetry.io/collector/service/internal/fanoutconsumer"

import (
	"context"

	"go.uber.org/multierr"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// NewTraces wraps multiple trace consumers in a single one sending the data marked as shared.
func NewTraces(tcs []consumer.Traces) consumer.Traces {
	if len(tcs) == 1 {
		// Don't wrap if no need to do it.
		return tcs[0]
	}
	return &tracesConsumer{consumers: tcs}
}

type tracesConsumer struct {
	consumers []consumer.Traces
}

// ConsumeTraces exports the ptrace.Traces to all consumers wrapped by the current one.
func (tsc *tracesConsumer) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	var errs error
	for _, tc := range tsc.consumers {
		// Send traces marked as shared so that they are cloned if mutation is needed.
		errs = multierr.Append(errs, tc.ConsumeTraces(ctx, td.AsShared()))
	}
	return errs
}
