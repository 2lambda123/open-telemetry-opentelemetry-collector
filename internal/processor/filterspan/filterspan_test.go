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

package filterspan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/internal/data/testdata"
	"go.opentelemetry.io/collector/internal/processor/filterset"
	"go.opentelemetry.io/collector/translator/conventions"
)

func createConfig(matchType filterset.MatchType) *filterset.Config {
	return &filterset.Config{
		MatchType: matchType,
	}
}

func TestSpan_validateMatchesConfiguration_InvalidConfig(t *testing.T) {
	testcases := []struct {
		name        string
		property    MatchProperties
		errorString string
	}{
		{
			name:        "empty_property",
			property:    MatchProperties{},
			errorString: errAtLeastOneMatchFieldNeeded.Error(),
		},
		{
			name: "empty_service_span_names_and_attributes",
			property: MatchProperties{
				Services: []string{},
			},
			errorString: errAtLeastOneMatchFieldNeeded.Error(),
		},
		{
			name: "invalid_match_type",
			property: MatchProperties{
				Config:   *createConfig("wrong_match_type"),
				Services: []string{"abc"},
			},
			errorString: "error creating service name filters: unrecognized match_type: 'wrong_match_type', valid types are: [regexp strict]",
		},
		{
			name: "missing_match_type",
			property: MatchProperties{
				Services: []string{"abc"},
			},
			errorString: "error creating service name filters: unrecognized match_type: '', valid types are: [regexp strict]",
		},
		{
			name: "regexp_match_type_for_int_attribute",
			property: MatchProperties{
				Config: *createConfig(filterset.Regexp),
				Attributes: []Attribute{
					{Key: "key", Value: 1},
				},
			},
			errorString: `error creating attribute filters: match_type=regexp for "key" only supports STRING, but found INT`,
		},
		{
			name: "unknown_attribute_value",
			property: MatchProperties{
				Config: *createConfig(filterset.Strict),
				Attributes: []Attribute{
					{Key: "key", Value: []string{}},
				},
			},
			errorString: `error creating attribute filters: error unsupported value type "[]string"`,
		},
		{
			name: "invalid_regexp_pattern",
			property: MatchProperties{
				Config:   *createConfig(filterset.Regexp),
				Services: []string{"["},
			},
			errorString: "error creating service name filters: error parsing regexp: missing closing ]: `[`",
		},
		{
			name: "invalid_regexp_pattern2",
			property: MatchProperties{
				Config:    *createConfig(filterset.Regexp),
				SpanNames: []string{"["},
			},
			errorString: "error creating span name filters: error parsing regexp: missing closing ]: `[`",
		},
		{
			name: "invalid_regexp_pattern3",
			property: MatchProperties{
				Config:     *createConfig(filterset.Regexp),
				SpanNames:  []string{"["},
				Attributes: []Attribute{{Key: "key", Value: "["}},
			},
			errorString: "error creating attribute filters: error parsing regexp: missing closing ]: `[`",
		},
		{
			name: "invalid_regexp_pattern4",
			property: MatchProperties{
				Config:    *createConfig(filterset.Regexp),
				Resources: []Attribute{{Key: "key", Value: "["}},
			},
			errorString: "error creating resource filters: error parsing regexp: missing closing ]: `[`",
		},
		{
			name: "empty_key_name_in_attributes_list",
			property: MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{"a"},
				Attributes: []Attribute{
					{
						Key: "",
					},
				},
			},
			errorString: "error creating attribute filters: error creating processor. Can't have empty key in the list of attributes",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := NewMatcher(&tc.property)
			assert.Nil(t, output)
			assert.EqualError(t, err, tc.errorString)
		})
	}
}

func TestSpan_Matching_False(t *testing.T) {
	testcases := []struct {
		name       string
		properties *MatchProperties
	}{
		{
			name: "service_name_doesnt_match_regexp",
			properties: &MatchProperties{
				Config:     *createConfig(filterset.Regexp),
				Services:   []string{"svcA"},
				Attributes: []Attribute{},
			},
		},

		{
			name: "service_name_doesnt_match_strict",
			properties: &MatchProperties{
				Config:     *createConfig(filterset.Strict),
				Services:   []string{"svcA"},
				Attributes: []Attribute{},
			},
		},

		{
			name: "span_name_doesnt_match",
			properties: &MatchProperties{
				Config:     *createConfig(filterset.Regexp),
				SpanNames:  []string{"spanNo.*Name"},
				Attributes: []Attribute{},
			},
		},

		{
			name: "span_name_doesnt_match_any",
			properties: &MatchProperties{
				Config: *createConfig(filterset.Regexp),
				SpanNames: []string{
					"spanNo.*Name",
					"non-matching?pattern",
					"regular string",
				},
				Attributes: []Attribute{},
			},
		},

		{
			name: "wrong_attribute_value",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{},
				Attributes: []Attribute{
					{
						Key:   "keyInt",
						Value: 1234,
					},
				},
			},
		},
		{
			name: "wrong_resource_value",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{},
				Resources: []Attribute{
					{
						Key:   "keyInt",
						Value: 1234,
					},
				},
			},
		},
		{
			name: "incompatible_attribute_value",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{},
				Attributes: []Attribute{
					{
						Key:   "keyInt",
						Value: "123",
					},
				},
			},
		},
		{
			name: "unsupported_attribute_value",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Regexp),
				Services: []string{},
				Attributes: []Attribute{
					{
						Key:   "keyMap",
						Value: "123",
					},
				},
			},
		},
		{
			name: "property_key_does_not_exist",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{},
				Attributes: []Attribute{
					{
						Key:   "doesnotexist",
						Value: nil,
					},
				},
			},
		},
	}

	span := pdata.NewSpan()
	span.InitEmpty()
	span.SetName("spanName")
	span.Attributes().InitFromMap(map[string]pdata.AttributeValue{
		"keyInt": pdata.NewAttributeValueInt(123),
		"keyMap": pdata.NewAttributeValueMap(),
	})
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			matcher, err := NewMatcher(tc.properties)
			assert.Nil(t, err)
			assert.NotNil(t, matcher)

			assert.False(t, matcher.MatchSpan(span, resource("wrongSvc")))
		})
	}
}

func resource(service string) pdata.Resource {
	r := pdata.NewResource()
	r.InitEmpty()
	r.Attributes().InitFromMap(map[string]pdata.AttributeValue{conventions.AttributeServiceName: pdata.NewAttributeValueString(service)})
	return r
}

func TestSpan_MatchingCornerCases(t *testing.T) {
	cfg := &MatchProperties{
		Config:   *createConfig(filterset.Strict),
		Services: []string{"svcA"},
		Attributes: []Attribute{
			{
				Key:   "keyOne",
				Value: nil,
			},
		},
	}

	mp, err := NewMatcher(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, mp)

	emptySpan := pdata.NewSpan()
	emptySpan.InitEmpty()
	assert.False(t, mp.MatchSpan(emptySpan, resource("svcA")))
}

func TestSpan_MissingServiceName(t *testing.T) {
	cfg := &MatchProperties{
		Config:   *createConfig(filterset.Regexp),
		Services: []string{"svcA"},
	}

	mp, err := NewMatcher(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, mp)

	emptySpan := pdata.NewSpan()
	emptySpan.InitEmpty()
	assert.False(t, mp.MatchSpan(emptySpan, resource("")))
}

func TestSpan_Matching_True(t *testing.T) {
	testcases := []struct {
		name       string
		properties *MatchProperties
	}{
		{
			name: "service_name_match_regexp",
			properties: &MatchProperties{
				Config:     *createConfig(filterset.Regexp),
				Services:   []string{"svcA"},
				Attributes: []Attribute{},
			},
		},
		{
			name: "service_name_match_strict",
			properties: &MatchProperties{
				Config:     *createConfig(filterset.Strict),
				Services:   []string{"svcA"},
				Attributes: []Attribute{},
			},
		},
		{
			name: "span_name_match",
			properties: &MatchProperties{
				Config:     *createConfig(filterset.Regexp),
				SpanNames:  []string{"span.*"},
				Attributes: []Attribute{},
			},
		},
		{
			name: "span_name_second_match",
			properties: &MatchProperties{
				Config: *createConfig(filterset.Regexp),
				SpanNames: []string{
					"wrong.*pattern",
					"span.*",
					"yet another?pattern",
					"regularstring",
				},
				Attributes: []Attribute{},
			},
		},
		{
			name: "attribute_exact_value_match",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{},
				Attributes: []Attribute{
					{
						Key:   "keyString",
						Value: "arithmetic",
					},
					{
						Key:   "keyInt",
						Value: 123,
					},
					{
						Key:   "keyDouble",
						Value: 3245.6,
					},
					{
						Key:   "keyBool",
						Value: true,
					},
				},
			},
		},
		{
			name: "attribute_regex_value_match",
			properties: &MatchProperties{
				Config: *createConfig(filterset.Regexp),
				Attributes: []Attribute{
					{
						Key:   "keyString",
						Value: "arith.*",
					},
					{
						Key:   "keyInt",
						Value: "12.*",
					},
					{
						Key:   "keyDouble",
						Value: "324.*",
					},
					{
						Key:   "keyBool",
						Value: "tr.*",
					},
				},
			},
		},
		{
			name: "resource_exact_value_match",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{},
				Resources: []Attribute{
					{
						Key:   "keyString",
						Value: "arithmetic",
					},
				},
			},
		},
		{
			name: "property_exists",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{"svcA"},
				Attributes: []Attribute{
					{
						Key:   "keyExists",
						Value: nil,
					},
				},
			},
		},
		{
			name: "match_all_settings_exists",
			properties: &MatchProperties{
				Config:   *createConfig(filterset.Strict),
				Services: []string{"svcA"},
				Attributes: []Attribute{
					{
						Key:   "keyExists",
						Value: nil,
					},
					{
						Key:   "keyString",
						Value: "arithmetic",
					},
				},
			},
		},
	}

	span := pdata.NewSpan()
	span.InitEmpty()
	span.SetName("spanName")
	span.Attributes().InitFromMap(map[string]pdata.AttributeValue{
		"keyString": pdata.NewAttributeValueString("arithmetic"),
		"keyInt":    pdata.NewAttributeValueInt(123),
		"keyDouble": pdata.NewAttributeValueDouble(3245.6),
		"keyBool":   pdata.NewAttributeValueBool(true),
		"keyExists": pdata.NewAttributeValueString("present"),
	})

	resource := pdata.NewResource()
	resource.InitEmpty()
	resource.Attributes().InitFromMap(map[string]pdata.AttributeValue{
		conventions.AttributeServiceName: pdata.NewAttributeValueString("svcA"),
		"keyString":                      pdata.NewAttributeValueString("arithmetic"),
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mp, err := NewMatcher(tc.properties)
			assert.Nil(t, err)
			assert.NotNil(t, mp)

			assert.NotNil(t, span)
			assert.True(t, mp.MatchSpan(span, resource))
		})
	}
}

func TestSpan_validateMatchesConfigurationForAttributes(t *testing.T) {
	testcase := []struct {
		name   string
		input  MatchProperties
		output Matcher
	}{
		{
			name: "attributes_build",
			input: MatchProperties{
				Config: *createConfig(filterset.Strict),
				Attributes: []Attribute{
					{
						Key: "key1",
					},
					{
						Key:   "key2",
						Value: 1234,
					},
				},
			},
			output: &propertiesMatcher{
				Attributes: []attributeMatcher{
					{
						Key: "key1",
					},
					{
						Key:            "key2",
						AttributeValue: newAttributeValueInt(1234),
					},
				},
			},
		},

		{
			name: "both_set_of_attributes",
			input: MatchProperties{
				Config: *createConfig(filterset.Strict),
				Attributes: []Attribute{
					{
						Key: "key1",
					},
					{
						Key:   "key2",
						Value: 1234,
					},
				},
			},
			output: &propertiesMatcher{
				Attributes: []attributeMatcher{
					{
						Key: "key1",
					},
					{
						Key:            "key2",
						AttributeValue: newAttributeValueInt(1234),
					},
				},
			},
		},
	}
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			output, err := NewMatcher(&tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.output, output)
		})
	}
}

func newAttributeValueInt(v int64) *pdata.AttributeValue {
	attr := pdata.NewAttributeValueInt(v)
	return &attr
}

func TestServiceNameForResource(t *testing.T) {
	td := testdata.GenerateTraceDataOneSpanNoResource()
	require.Equal(t, serviceNameForResource(td.ResourceSpans().At(0).Resource()), "<nil-resource>")

	td = testdata.GenerateTraceDataOneSpan()
	resource := td.ResourceSpans().At(0).Resource()
	require.Equal(t, serviceNameForResource(resource), "<nil-service-name>")

	resource.Attributes().InsertString(conventions.AttributeServiceName, "test-service")
	require.Equal(t, serviceNameForResource(resource), "test-service")
}
