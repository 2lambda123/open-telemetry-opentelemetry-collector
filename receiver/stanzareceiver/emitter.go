package stanzareceiver

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

import (
	"context"
	"sync"

	"github.com/observiq/stanza/entry"
	"github.com/observiq/stanza/operator/helper"
	"go.uber.org/zap"
)

// LogEmitter is a stanza operator that emits log entries to a channel
type LogEmitter struct {
	helper.OutputOperator
	logChan  chan *entry.Entry
	stopOnce sync.Once
}

// NewLogEmitter creates a new receiver output
func NewLogEmitter(logger *zap.SugaredLogger) *LogEmitter {
	return &LogEmitter{
		OutputOperator: helper.OutputOperator{
			BasicOperator: helper.BasicOperator{
				OperatorID:    "log_emitter",
				OperatorType:  "log_emitter",
				SugaredLogger: logger,
			},
		},
		logChan: make(chan *entry.Entry),
	}
}

// Process will emit an entry to the output channel
func (e *LogEmitter) Process(ctx context.Context, entry *entry.Entry) error {
	e.logChan <- entry
	return nil
}

// Stop will close the log channel
func (e *LogEmitter) Stop() error {
	e.stopOnce.Do(func() {
		close(e.logChan)
	})
	return nil
}

func (e *LogEmitter) LogChan() chan *entry.Entry {
	return e.logChan
}
