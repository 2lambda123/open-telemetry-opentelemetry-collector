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

package configcompression // import "go.opentelemetry.io/collector/config/configcompression"

import "fmt"

type CompressionType string

const (
	CompressionGzip    CompressionType = "gzip"
	CompressionZlib    CompressionType = "zlib"
	CompressionDeflate CompressionType = "deflate"
	CompressionSnappy  CompressionType = "snappy"
	CompressionZstd    CompressionType = "zstd"
	compressionNone    CompressionType = "none"
	compressionEmpty   CompressionType = ""
)

func IsCompressed(compressionType CompressionType) bool {
	return compressionType != compressionEmpty && compressionType != compressionNone
}

func (ct *CompressionType) UnmarshalText(in []byte) error {
	switch typ := CompressionType(in); typ {
	case CompressionGzip,
		CompressionZlib,
		CompressionDeflate,
		CompressionSnappy,
		CompressionZstd,
		compressionNone,
		compressionEmpty:
		*ct = typ
		return nil
	default:
		return fmt.Errorf("unsupported compression type %q", typ)
	}
}
