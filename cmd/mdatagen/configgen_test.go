package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: An **incomplete** list of things to fix with go-jsonschema:
// * It doesn't use "uint" anywhere.
// * For "integer", it does't use "minimum: 0" inside UnmarshalJSON.
// * Durations show up as float64 instead of as time.Duration.
func TestXxx(t *testing.T) {
	inputDir := `./testdata/config_gen/input_schema`
	outputDir := `./testdata/config_gen/expected_golang_output/`

	inputFiles, err := os.ReadDir(inputDir)
	require.NoError(t, err)

	for _, inputFile := range inputFiles {
		if inputFile.IsDir() {
			continue
		}

		md, err := loadMetadata(filepath.Join(inputDir, inputFile.Name()))
		require.NoError(t, err)

		buf := new(bytes.Buffer)

		err = GenerateConfig(md.Config, buf)
		require.NoError(t, err)

		actual := buf.String()

		expectedOutputFile := filepath.Join(outputDir, inputFile.Name())
		var expectedOutput []byte
		expectedOutput, err = os.ReadFile(expectedOutputFile)
		require.NoError(t, err)

		require.Equal(t, string(expectedOutput), actual)
	}
}
