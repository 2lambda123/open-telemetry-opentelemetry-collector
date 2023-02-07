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

package consumer // import "go.opentelemetry.io/collector/consumer"

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
)

// Logs is an interface that receives plog.Logs, processes it
// as needed, and sends it to the next processing node if any or to the destination.
type Logs interface {
	// ConsumeLogs receives plog.Logs for consumption.
	ConsumeLogs(ctx context.Context, ld plog.Logs) error
}

// ConsumeLogsFunc is a helper function that is similar to ConsumeLogs.
type ConsumeLogsFunc func(ctx context.Context, ld plog.Logs) error

// ConsumeLogs calls f(ctx, ld).
func (f ConsumeLogsFunc) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	return f(ctx, ld)
}

// NewLogs returns a Logs configured with the provided options.
func NewLogs(consume ConsumeLogsFunc, _ ...Option) (Logs, error) {
	if consume == nil {
		return nil, errNilFunc
	}
	return consume, nil
}
