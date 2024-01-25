// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otelcoltest

import (
	"testing"

	"go.uber.org/goleak"
)

// The IgnoreTopFunction call prevents catching the leak generated by opencensus
// defaultWorker.Start which at this time is part of the package's init call.
// See https://github.com/open-telemetry/opentelemetry-collector/issues/9165#issuecomment-1874836336 for more context.
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
		goleak.IgnoreTopFunction("go.opentelemetry.io/collector/service/internal/proctelemetry.InitPrometheusServer"),
	)
}
