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

package plog // import "go.opentelemetry.io/collector/pdata/plog"

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpcollectorlog "go.opentelemetry.io/collector/pdata/internal/data/protogen/collector/logs/v1"
	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
)

// Logs is the top-level struct that is propagated through the logs pipeline.
// Use NewLogs to create new instance, zero-initialized instance is not valid for use.
type Logs internal.Logs

func newLogs(orig *otlpcollectorlog.ExportLogsServiceRequest) Logs {
	return Logs(internal.NewLogs(orig))
}

func newLogsFromResourceLogsOrig(orig *[]*otlplogs.ResourceLogs) Logs {
	return Logs(internal.NewLogsFromResourceLogsOrig(orig))
}

func (ms Logs) getOrig() *otlpcollectorlog.ExportLogsServiceRequest {
	return internal.Logs(ms).GetOrig()
}

// NewLogs creates a new Logs struct.
func NewLogs() Logs {
	return newLogs(&otlpcollectorlog.ExportLogsServiceRequest{})
}

func (ms Logs) AsShared() Logs {
	return Logs(internal.Logs(ms).AsShared())
}

func (ms Logs) ensureMutability() {
	if *internal.Logs(ms).GetState() == internal.StateShared {
		*internal.Logs(ms).GetState() = internal.StateDirty
		newRL := newResourceLogsSliceFromOrig(&[]*otlplogs.ResourceLogs{})
		ms.ResourceLogs().CopyTo(newRL)
		newState := internal.StateExclusive
		internal.Logs(ms).SetState(&newState)
		internal.Logs(ms).SetOrig(&otlpcollectorlog.ExportLogsServiceRequest{
			ResourceLogs: *newRL.getOrig(),
		})
	}
}

func (ms Logs) getResourceLogsOrig() *[]*otlplogs.ResourceLogs {
	return &internal.Logs(ms).GetOrig().ResourceLogs
}

func (ms Logs) getState() *internal.State {
	return internal.Logs(ms).GetState()
}

func (ms Logs) refreshResourceLogsOrigState() (*[]*otlplogs.ResourceLogs, *internal.State) {
	return &internal.Logs(ms).GetOrig().ResourceLogs, ms.getState()
}

// CopyTo copies the Logs instance overriding the destination.
func (ms Logs) CopyTo(dest Logs) {
	ms.ResourceLogs().CopyTo(dest.ResourceLogs())
}

// LogRecordCount calculates the total number of log records.
func (ms Logs) LogRecordCount() int {
	logCount := 0
	rss := ms.ResourceLogs()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		ill := rs.ScopeLogs()
		for i := 0; i < ill.Len(); i++ {
			logs := ill.At(i)
			logCount += logs.LogRecords().Len()
		}
	}
	return logCount
}

// ResourceLogs returns the ResourceLogsSlice associated with this Logs.
func (ms Logs) ResourceLogs() ResourceLogsSlice {
	return newResourceLogsSliceFromParent(ms)
}
