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

package addattributesprocessor

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-service/config"
	"github.com/open-telemetry/opentelemetry-service/config/configmodels"
	"github.com/open-telemetry/opentelemetry-service/processor"
)

var _ = config.RegisterTestFactories()

func TestLoadConfig(t *testing.T) {

	factory := processor.GetFactory(typeStr)

	cfg, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"))

	require.Nil(t, err)
	require.NotNil(t, cfg)

	p0 := cfg.Processors["attributes"]
	assert.Equal(t, p0, factory.CreateDefaultConfig())

	p1 := cfg.Processors["attributes/2"]
	assert.Equal(t, p1,
		&Config{
			ProcessorSettings: configmodels.ProcessorSettings{
				TypeVal: "attributes",
				NameVal: "attributes/2",
			},
			Values: map[string]interface{}{
				"attribute1":         123,
				"string attribute":   "string value",
				"attribute.with.dot": "another value",
			},
		})
}
