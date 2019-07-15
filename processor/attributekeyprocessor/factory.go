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

package attributekeyprocessor

import (
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-service/config/configerror"
	"github.com/open-telemetry/opentelemetry-service/config/configmodels"
	"github.com/open-telemetry/opentelemetry-service/consumer"
	"github.com/open-telemetry/opentelemetry-service/processor"
)

var _ = processor.RegisterFactory(&factory{})

const (
	// The value of "type" Attribute Key in configuration.
	typeStr = "attributes"
)

// factory is the factory for Attribute Key processor.
type factory struct {
}

// Type gets the type of the config created by this factory.
func (f *factory) Type() string {
	return typeStr
}

// CreateDefaultConfig creates the default configuration for processor.
func (f *factory) CreateDefaultConfig() configmodels.Processor {
	return &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		KeyReplacements: make(map[string]NewKeyProperties, 0),
	}
}

// CreateTraceProcessor creates a trace processor based on this config.
func (f *factory) CreateTraceProcessor(
	logger *zap.Logger,
	nextConsumer consumer.TraceConsumer,
	cfg configmodels.Processor,
) (processor.TraceProcessor, error) {
	oCfg := cfg.(*Config)
	return NewTraceProcessor(nextConsumer, convertToKeyReplacements(&oCfg.KeyReplacements)...)
}

// CreateMetricsProcessor creates a metrics processor based on this config.
func (f *factory) CreateMetricsProcessor(
	logger *zap.Logger,
	nextConsumer consumer.MetricsConsumer,
	cfg configmodels.Processor,
) (processor.MetricsProcessor, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}

// convert key replacments' "map" to KeyReplacement
func convertToKeyReplacements(keyMap *map[string]NewKeyProperties) []KeyReplacement {
	var replacements []KeyReplacement
	for key, val := range *keyMap {
		replacements = append(replacements, KeyReplacement{Key: key, NewKey: val.NewKey, Overwrite: val.Overwrite, KeepOriginal: val.KeepOriginal})
	}
	return replacements
}
