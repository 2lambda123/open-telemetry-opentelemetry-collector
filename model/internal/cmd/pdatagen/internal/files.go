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

package internal // import "go.opentelemetry.io/collector/model/internal/cmd/pdatagen/internal"

import "strings"

const header = `// Copyright The OpenTelemetry Authors
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

// Code generated by "model/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "go run model/internal/cmd/pdatagen/main.go".

package pdata`

// AllFiles is a list of all files that needs to be generated.
var AllFiles = []*File{
	commonFile,
	metricsFile,
	resourceFile,
	traceFile,
	logFile,
}

// File represents the struct for one generated file.
type File struct {
	Name        string
	imports     []string
	testImports []string
	// Can be any of sliceOfPtrs, sliceOfValues, messageValueStruct, or messagePtrStruct
	structs []baseStruct
}

// GenerateFile generates the configured data structures for this File.
func (f *File) GenerateFile() string {
	var sb strings.Builder

	// Write headers
	sb.WriteString(header)
	sb.WriteString(newLine + newLine)
	// Add imports
	sb.WriteString("import (" + newLine)
	for _, imp := range f.imports {
		if imp != "" {
			sb.WriteString("\t" + imp + newLine)
		} else {
			sb.WriteString(newLine)
		}
	}
	sb.WriteString(")")
	// Write all structs
	for _, s := range f.structs {
		sb.WriteString(newLine + newLine)
		s.generateStruct(&sb)
	}
	sb.WriteString(newLine)
	return sb.String()
}

// GenerateTestFile generates tests for the configured data structures for this File.
func (f *File) GenerateTestFile() string {
	var sb strings.Builder

	// Write headers
	sb.WriteString(header)
	sb.WriteString(newLine + newLine)
	// Add imports
	sb.WriteString("import (" + newLine)
	for _, imp := range f.testImports {
		if imp != "" {
			sb.WriteString("\t" + imp + newLine)
		} else {
			sb.WriteString(newLine)
		}
	}
	sb.WriteString(")")
	// Write all tests
	for _, s := range f.structs {
		sb.WriteString(newLine + newLine)
		s.generateTests(&sb)
	}
	// Write all tests generate value
	for _, s := range f.structs {
		sb.WriteString(newLine + newLine)
		s.generateTestValueHelpers(&sb)
	}
	sb.WriteString(newLine)
	return sb.String()
}

// GenerateFile generates the aliases for data structures for this File.
func (f *File) GenerateAliasFile() string {
	var sb strings.Builder

	// Write headers
	sb.WriteString(header)
	sb.WriteString(newLine + newLine)

	// Add import
	sb.WriteString("import \"go.opentelemetry.io/collector/model/internal/pdata\"" + newLine + newLine)

	// Write all types and funcs
	for _, s := range f.structs {
		if ag, ok := s.(aliasGenerator); ok {
			ag.generateAlias(&sb)
		}
	}
	sb.WriteString(newLine)
	return sb.String()
}
