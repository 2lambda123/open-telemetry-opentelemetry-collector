// Copyright The OpenTelemetry Authors
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

package internal

var resourceFile = &File{
	Name: "resource",
	imports: []string{
		`otlpresource "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/resource/v1"`,
	},
	testImports: []string{
		`"testing"`,
		``,
		`"github.com/stretchr/testify/assert"`,
	},
	structs: []baseStruct{
		resource,
	},
}

var resource = &messageStruct{
	structName:     "Resource",
	description:    "// Resource information.",
	originFullName: "otlpresource.Resource",
	fields: []baseField{
		attributes,
	},
}

var resourceField = &messageField{
	fieldName:       "Resource",
	originFieldName: "Resource",
	returnMessage:   resource,
}
