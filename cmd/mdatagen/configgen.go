// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/atombender/go-jsonschema/pkg/generator"
	"github.com/atombender/go-jsonschema/pkg/schemas"
)

const (
	//TODO: Get rid of this?
	CONFIG_NAME = "config"
)

func GenerateConfig(conf any, output io.Writer) error {
	// load config
	jsonBytes, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed loading config %w", err)
	}
	var schema schemas.Schema
	if err := json.Unmarshal(jsonBytes, &schema); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// TODO: don't hardcode
	pkgName := "batchprocessor"

	// init generator
	cfg := generator.Config{
		Warner: func(message string) {
			logf("Warning: %s", message)
		},
		DefaultPackageName:  pkgName,
		DefaultOutputName:   "config",
		StructNameFromTitle: true,
		Tags:                []string{"json", "yaml", "mapstructure"},
		SchemaMappings:      []generator.SchemaMapping{},
		YAMLExtensions:      []string{".yaml", ".yml"},
		// YAMLPackage:         "gopkg.in/yaml.v3",
	}

	generator, err := generator.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}
	if err = generator.AddFile(CONFIG_NAME, &schema); err != nil {
		return fmt.Errorf("failed to add config: %w", err)
	}

	// hasUnsupportedValidations := len(generator.NotSupportedValidations) > 0

	tplVars := struct {
		ValidatorFuncName string
	}{
		ValidatorFuncName: "Validate",
	}
	// if hasUnsupportedValidations {
	// 	tplVars.ValidatorFuncName = "ValidateHelper"
	// }

	tpl := `
func (cfg *Config){{.ValidatorFuncName}}() error {
	b, err := json.Marshal(cfg)
	if err != nil {
			return err
	}
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
			return err
	}
	return nil
}`
	tmpl, err := template.New("validator").Parse(tpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	for _, source := range generator.Sources() {
		buf := bytes.NewBufferString("")
		if err = tmpl.Execute(buf, tplVars); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
		// only write custom validation if there are no unsupported validations
		// source = append(source, []byte(tpl)...)
		source = append(source, buf.Bytes()...)
		_, err = output.Write(source)

		if err != nil {
			return fmt.Errorf("failed writing file: %w", err)
		}
	}
	fmt.Println("done")
	return nil
}

func logf(format string, args ...interface{}) {
	fmt.Fprint(os.Stderr, "go-jsonschema: ")
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprint(os.Stderr, "\n")
}
