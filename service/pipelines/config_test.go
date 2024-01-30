// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pipelines

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/component"
)

func TestConfigValidate(t *testing.T) {
	var testCases = []struct {
		name     string // test case name (also file name containing config yaml)
		cfgFn    func() Config
		expected error
	}{
		{
			name:     "valid",
			cfgFn:    generateConfig,
			expected: nil,
		},
		{
			name: "duplicate-processor-reference",
			cfgFn: func() Config {
				cfg := generateConfig()
				pipe := cfg[component.NewID(component.MustNewType("traces"))]
				pipe.Processors = append(pipe.Processors, pipe.Processors...)
				return cfg
			},
			expected: fmt.Errorf(`pipeline "traces": %w`, errors.New(`references processor "nop" multiple times`)),
		},
		{
			name: "missing-pipeline-receivers",
			cfgFn: func() Config {
				cfg := generateConfig()
				cfg[component.NewID(component.MustNewType("traces"))].Receivers = nil
				return cfg
			},
			expected: fmt.Errorf(`pipeline "traces": %w`, errMissingServicePipelineReceivers),
		},
		{
			name: "missing-pipeline-exporters",
			cfgFn: func() Config {
				cfg := generateConfig()
				cfg[component.NewID(component.MustNewType("traces"))].Exporters = nil
				return cfg
			},
			expected: fmt.Errorf(`pipeline "traces": %w`, errMissingServicePipelineExporters),
		},
		{
			name: "missing-pipelines",
			cfgFn: func() Config {
				return nil
			},
			expected: errMissingServicePipelines,
		},
		{
			name: "invalid-service-pipeline-type",
			cfgFn: func() Config {
				cfg := generateConfig()
				cfg[component.NewID(component.MustNewType("wrongtype"))] = &PipelineConfig{
					Receivers:  []component.ID{component.NewID(component.MustNewType("nop"))},
					Processors: []component.ID{component.NewID(component.MustNewType("nop"))},
					Exporters:  []component.ID{component.NewID(component.MustNewType("nop"))},
				}
				return cfg
			},
			expected: errors.New(`pipeline "wrongtype": unknown datatype "wrongtype"`),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			cfg := test.cfgFn()
			assert.Equal(t, test.expected, cfg.Validate())
		})
	}
}

func generateConfig() Config {
	return map[component.ID]*PipelineConfig{
		component.NewID(component.MustNewType("traces")): {
			Receivers:  []component.ID{component.NewID(component.MustNewType("nop"))},
			Processors: []component.ID{component.NewID(component.MustNewType("nop"))},
			Exporters:  []component.ID{component.NewID(component.MustNewType("nop"))},
		},
	}
}
