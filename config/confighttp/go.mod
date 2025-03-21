module go.opentelemetry.io/collector/config/confighttp

go 1.22.0
toolchain go1.23.7

require (
	github.com/golang/snappy v0.0.4
	github.com/klauspost/compress v1.17.11
	github.com/pierrec/lz4/v4 v4.1.21
	github.com/rs/cors v1.11.1
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/client v1.17.0
	go.opentelemetry.io/collector/component v0.111.0
	go.opentelemetry.io/collector/config/configauth v0.111.0
	go.opentelemetry.io/collector/config/configcompression v1.17.0
	go.opentelemetry.io/collector/config/configopaque v1.17.0
	go.opentelemetry.io/collector/config/configtelemetry v0.111.0
	go.opentelemetry.io/collector/config/configtls v1.17.0
	go.opentelemetry.io/collector/config/internal v0.111.0
	go.opentelemetry.io/collector/extension/auth v0.111.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.56.0
	go.opentelemetry.io/otel v1.31.0
	go.opentelemetry.io/otel/metric v1.31.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.36.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/collector/extension v0.111.0 // indirect
	go.opentelemetry.io/collector/pdata v1.17.0 // indirect
	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240822170219-fc7c04adadcd // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/collector/config/configauth => ../configauth

replace go.opentelemetry.io/collector/config/configcompression => ../configcompression

replace go.opentelemetry.io/collector/config/configopaque => ../configopaque

replace go.opentelemetry.io/collector/config/configtls => ../configtls

replace go.opentelemetry.io/collector/config/configtelemetry => ../configtelemetry

replace go.opentelemetry.io/collector/config/internal => ../internal

replace go.opentelemetry.io/collector/extension => ../../extension

replace go.opentelemetry.io/collector/extension/auth => ../../extension/auth

replace go.opentelemetry.io/collector/pdata => ../../pdata

replace go.opentelemetry.io/collector/component => ../../component

replace go.opentelemetry.io/collector/consumer => ../../consumer

replace go.opentelemetry.io/collector/client => ../../client

replace go.opentelemetry.io/collector/pdata/testdata => ../../pdata/testdata

replace go.opentelemetry.io/collector/pdata/pprofile => ../../pdata/pprofile
