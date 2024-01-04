// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ballastextension

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/internal/iruntime"
)

func TestMemoryBallast(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		getTotalMem func() (uint64, error)
		expect      int
		expectErr   bool
	}{
		{
			name: "test_abs_ballast",
			config: &Config{
				SizeMiB: 13,
			},
			getTotalMem: iruntime.TotalMemory,
			expect:      13 * megaBytes,
		},
		{
			name: "test_abs_ballast_priority",
			config: &Config{
				SizeMiB:          13,
				SizeInPercentage: 20,
			},
			getTotalMem: iruntime.TotalMemory,
			expect:      13 * megaBytes,
		},
		{
			name:        "test_ballast_zero_val",
			config:      &Config{},
			getTotalMem: iruntime.TotalMemory,
			expect:      0,
		},
		{
			name: "test_ballast_in_percentage",
			config: &Config{
				SizeInPercentage: 20,
			},
			getTotalMem: func() (uint64, error) { return uint64(100 * megaBytes), nil },
			expect:      20 * megaBytes,
		},
		{
			name: "failing_get_total_mem",
			config: &Config{
				SizeInPercentage: 20,
			},
			getTotalMem: func() (uint64, error) { return 0, errors.New("failed") },
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mbExt, err := newMemoryBallast(tt.config, zap.NewNop(), tt.getTotalMem)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			require.NotNil(t, mbExt)
			require.NoError(t, err)
			assert.Equal(t, tt.expect, len(mbExt.ballast))

			assert.NoError(t, mbExt.Start(context.Background(), componenttest.NewNopHost()))
			assert.NoError(t, mbExt.Shutdown(context.Background()))
			assert.Nil(t, mbExt.ballast)
		})
	}
}
