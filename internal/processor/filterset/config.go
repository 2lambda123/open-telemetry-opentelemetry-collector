// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filterset

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector/internal/processor/filterset/regexp"
	"github.com/open-telemetry/opentelemetry-collector/internal/processor/filterset/strict"
)

// MatchType describes the type of pattern matching a FilterSet uses to filter strings.
type MatchType string

const (
	// Regexp is the FilterType for filtering by regexp string matches.
	Regexp MatchType = "regexp"
	// Strict is the FilterType for filtering by exact string matches.
	Strict MatchType = "strict"
)

// Config configures the matching behavior of a FilterSet.
type Config struct {
	MatchType    MatchType      `mapstructure:"match_type"`
	RegexpConfig *regexp.Config `mapstructure:"regexp"`
}

// CreateFilterSet creates a FilterSet from yaml config.
func CreateFilterSet(filters []string, cfg *Config) (FilterSet, error) {
	switch cfg.MatchType {
	case Regexp:
		return regexp.CreateRegexpFilterSet(filters, cfg.RegexpConfig)
	case Strict:
		// Strict FilterSets do not have any extra configuration options, so call the constructor directly.
		return strict.NewStrictFilterSet(filters)
	default:
		return nil, fmt.Errorf("unrecognized filter type: %v", cfg.MatchType)
	}
}
