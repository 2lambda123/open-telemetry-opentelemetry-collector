module go.opentelemetry.io/collector/config/configtls

go 1.20

require (
	github.com/fsnotify/fsnotify v1.6.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector/config/configopaque v0.83.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/collector/config/configopaque => ../configopaque
