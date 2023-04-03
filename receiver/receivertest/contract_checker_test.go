// Copyright The OpenTelemetry Authors
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

package receivertest

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
)

// This file is an example that demonstrates how to use the CheckConsumeContract() function.
// We declare a trivial example receiver, a data generator and then use them in TestConsumeContract().

type exampleReceiver struct {
	nextConsumer consumer.Logs
}

func (s *exampleReceiver) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (s *exampleReceiver) Shutdown(_ context.Context) error {
	return nil
}

func (s *exampleReceiver) Receive(data plog.Logs) {
	// This very simple implementation demonstrates how a single items receiving should happen.
	for {
		err := s.nextConsumer.ConsumeLogs(context.Background(), data)
		if err != nil {
			// The next consumer returned an error.
			if !consumererror.IsPermanent(err) {
				// It is not a permanent error, so we must retry sending it again. In network-based
				// receivers instead we can ask our sender to re-retry the same data again later.
				// We may also pause here a bit if we don't want to hammer the next consumer.
				continue
			}
		}
		// If we are hear either the ConsumeLogs returned success or it returned a permanent error.
		// In either case we don't need to retry the same data, we are done.
		return
	}
}

// A config for exampleReceiver.
type exampleReceiverConfig struct {
	generator *exampleGenerator
}

// A generator that can send data to exampleReceiver.
type exampleGenerator struct {
	t           *testing.T
	receiver    *exampleReceiver
	sequenceNum int64
}

func (g *exampleGenerator) Generate() IDSet {
	// Make sure the id is atomically incremented. Generate() may be called concurrently.
	id := atomic.AddInt64(&g.sequenceNum, 1)

	data := CreateOneLogWithID(UniqueIDAttrDataType(id))

	// Send the generated data to the recever.
	g.receiver.Receive(data)

	// And return the ids for bookkeeping by the test.
	ids, err := IDSetFromLogs(data)
	require.NoError(g.t, err)

	return ids
}

func newExampleFactory() receiver.Factory {
	return receiver.NewFactory(
		"example_receiver",
		func() component.Config {
			return &exampleReceiverConfig{}
		},
		receiver.WithLogs(createLog, component.StabilityLevelBeta),
	)
}

func createLog(
	_ context.Context,
	_ receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	rcv := &exampleReceiver{nextConsumer: consumer}
	cfg.(*exampleReceiverConfig).generator.receiver = rcv
	return rcv, nil
}

// TestConsumeContract is an example of testing of the receiver for the contract between the
// receiver and next consumer.
func TestConsumeContract(t *testing.T) {

	// Different scenarios to test for.
	decisionFuncs := []func(ids IDSet) error{
		// Always succeed. We expect all data to be delivered as is.
		func(ids IDSet) error { return nil },

		// Various scenarios with errors injected into the consumer's decision making.
		RandomNonPermanentErrorConsumeDecision,
		RandomPermanentErrorConsumeDecision,
		RandomErrorsConsumeDecision,
	}

	// Number of log records to send per scenario.
	const logsPerTest = 100

	for i, decisionFunc := range decisionFuncs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			generator := &exampleGenerator{t: t}
			cfg := &exampleReceiverConfig{generator: generator}

			params := CheckConsumeContractParams{
				T:                   t,
				Factory:             newExampleFactory(),
				Config:              cfg,
				Generator:           generator,
				GenerateCount:       logsPerTest,
				ConsumeDecisionFunc: decisionFunc,
			}

			// Run the contract checker. This will trigger test failures if any problems are found.
			CheckConsumeContract(params)
		})
	}
}
