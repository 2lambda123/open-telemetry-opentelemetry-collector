module go.opentelemetry.io/collector/exporter/exporterprofiles

go 1.21.0

require (
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/component v0.106.1
	go.opentelemetry.io/collector/consumer/consumertest v0.106.1
	go.opentelemetry.io/collector/exporter v0.106.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.106.1 // indirect
	go.opentelemetry.io/collector/consumer v0.106.1 // indirect
	go.opentelemetry.io/collector/consumer/consumerprofiles v0.106.1 // indirect
	go.opentelemetry.io/collector/pdata v1.13.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.106.1 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/collector/consumer/consumerprofiles => ../../consumer/consumerprofiles

replace go.opentelemetry.io/collector/internal/globalgates => ../../internal/globalgates

replace go.opentelemetry.io/collector/extension => ../../extension

replace go.opentelemetry.io/collector => ../..

replace go.opentelemetry.io/collector/pdata/pprofile => ../../pdata/pprofile

replace go.opentelemetry.io/collector/config/configtelemetry => ../../config/configtelemetry

replace go.opentelemetry.io/collector/config/configretry => ../../config/configretry

replace go.opentelemetry.io/collector/pdata/testdata => ../../pdata/testdata

replace go.opentelemetry.io/collector/featuregate => ../../featuregate

replace go.opentelemetry.io/collector/confmap => ../../confmap

replace go.opentelemetry.io/collector/component => ../../component

replace go.opentelemetry.io/collector/pdata => ../../pdata

replace go.opentelemetry.io/collector/receiver => ../../receiver

replace go.opentelemetry.io/collector/consumer/consumertest => ../../consumer/consumertest

replace go.opentelemetry.io/collector/consumer => ../../consumer

replace go.opentelemetry.io/collector/exporter => ../
