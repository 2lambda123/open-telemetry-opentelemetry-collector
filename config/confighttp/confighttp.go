// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package confighttp

import "encoding/json"
import "fmt"
import configauth "go.opentelemetry.io/collector/config/configauth"
import configcompression "go.opentelemetry.io/collector/config/configcompression"
import "go.opentelemetry.io/collector/config/configopaque"
import configtls "go.opentelemetry.io/collector/config/configtls"
import "time"

type ClientConfig struct {
	// TLSSetting corresponds to the JSON schema field "TLSSetting".
	TLSSetting *configtls.ClientConfig `mapstructure:"TLSSetting"`

	// Auth corresponds to the JSON schema field "auth".
	Auth *configauth.Authentication `mapstructure:"auth"`

	// Compression corresponds to the JSON schema field "compression".
	Compression configcompression.Type `mapstructure:"compression"`

	// Cookies corresponds to the JSON schema field "cookies".
	Cookies *ClientConfigCookies `mapstructure:"cookies"`

	// DisableKeepAlives corresponds to the JSON schema field "disable_keep_alives".
	DisableKeepAlives bool `mapstructure:"disable_keep_alives"`

	// Endpoint corresponds to the JSON schema field "endpoint".
	Endpoint string `mapstructure:"endpoint"`

	// Headers corresponds to the JSON schema field "headers".
	Headers map[string]configopaque.String `mapstructure:"headers"`

	// Http2PingTimeout corresponds to the JSON schema field "http2_ping_timeout".
	Http2PingTimeout time.Duration `mapstructure:"http2_ping_timeout"`

	// Http2ReadIdleTimeout corresponds to the JSON schema field
	// "http2_read_idle_timeout".
	Http2ReadIdleTimeout time.Duration `mapstructure:"http2_read_idle_timeout"`

	// IdleConnTimeout corresponds to the JSON schema field "idle_conn_timeout".
	IdleConnTimeout *time.Duration `mapstructure:"idle_conn_timeout"`

	// MaxConnsPerHost corresponds to the JSON schema field "max_conns_per_host".
	MaxConnsPerHost *int `mapstructure:"max_conns_per_host"`

	// MaxIdleConns corresponds to the JSON schema field "max_idle_conns".
	MaxIdleConns *int `mapstructure:"max_idle_conns"`

	// MaxIdleConnsPerHost corresponds to the JSON schema field
	// "max_idle_conns_per_host".
	MaxIdleConnsPerHost *int `mapstructure:"max_idle_conns_per_host"`

	// ProxyUrl corresponds to the JSON schema field "proxy_url".
	ProxyUrl string `mapstructure:"proxy_url"`

	// ReadBufferSize corresponds to the JSON schema field "read_buffer_size".
	ReadBufferSize int `mapstructure:"read_buffer_size"`

	// Timeout corresponds to the JSON schema field "timeout".
	Timeout time.Duration `mapstructure:"timeout"`

	// WriteBufferSize corresponds to the JSON schema field "write_buffer_size".
	WriteBufferSize int `mapstructure:"write_buffer_size"`
}

type ClientConfigCookies struct {
	// Enabled corresponds to the JSON schema field "enabled".
	Enabled bool `mapstructure:"enabled"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ClientConfigCookies) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ClientConfigCookies
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["enabled"]; !ok || v == nil {
		plain.Enabled = false
	}
	*j = ClientConfigCookies(plain)
	return nil
}

// SetDefaults sets the fields of ClientConfigCookies to their defaults.
// Fields which do not have a default value are left untouched.
func (c *ClientConfigCookies) SetDefaults() {
	c.Enabled = false
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ClientConfig) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ClientConfig
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["compression"]; !ok || v == nil {
		plain.Compression = "none"
	}
	if v, ok := raw["disable_keep_alives"]; !ok || v == nil {
		plain.DisableKeepAlives = false
	}
	if v, ok := raw["endpoint"]; !ok || v == nil {
		plain.Endpoint = ""
	}
	if v, ok := raw["http2_ping_timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field http2_ping_timeout:%w }", err)
		}
		plain.Http2PingTimeout = defaultDuration
	}
	if v, ok := raw["http2_read_idle_timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field http2_read_idle_timeout:%w }", err)
		}
		plain.Http2ReadIdleTimeout = defaultDuration
	}
	if v, ok := raw["proxy_url"]; !ok || v == nil {
		plain.ProxyUrl = ""
	}
	if v, ok := raw["read_buffer_size"]; !ok || v == nil {
		plain.ReadBufferSize = 0.0
	}
	if v, ok := raw["timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field timeout:%w }", err)
		}
		plain.Timeout = defaultDuration
	}
	if v, ok := raw["write_buffer_size"]; !ok || v == nil {
		plain.WriteBufferSize = 0.0
	}
	*j = ClientConfig(plain)
	return nil
}

// SetDefaults sets the fields of ClientConfig to their defaults.
// Fields which do not have a default value are left untouched.
func (c *ClientConfig) SetDefaults() {
	c.Compression = "none"
	c.DisableKeepAlives = false
	c.Endpoint = ""
	c.Http2PingTimeout = "PT33.3S"
	c.Http2ReadIdleTimeout = "PT33.3S"
	c.ProxyUrl = ""
	c.ReadBufferSize = 0.0
	c.Timeout = "PT33.3S"
	c.WriteBufferSize = 0.0
}

type Cors struct {
	// AllowedHeaders corresponds to the JSON schema field "allowed_headers".
	AllowedHeaders []string `mapstructure:"allowed_headers"`

	// AllowedOrigins corresponds to the JSON schema field "allowed_origins".
	AllowedOrigins []string `mapstructure:"allowed_origins"`

	// MaxAge corresponds to the JSON schema field "max_age".
	MaxAge int `mapstructure:"max_age"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Cors) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain Cors
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["max_age"]; !ok || v == nil {
		plain.MaxAge = 0.0
	}
	*j = Cors(plain)
	return nil
}

// SetDefaults sets the fields of Cors to their defaults.
// Fields which do not have a default value are left untouched.
func (c *Cors) SetDefaults() {
	c.MaxAge = 0.0
}

type ServerConfig struct {
	// TLSSetting corresponds to the JSON schema field "TLSSetting".
	TLSSetting *configtls.ServerConfig `mapstructure:"TLSSetting"`

	// Auth corresponds to the JSON schema field "auth".
	Auth *ServerConfigAuth `mapstructure:"auth"`

	// CompressionAlgorithms corresponds to the JSON schema field
	// "compression_algorithms".
	CompressionAlgorithms configcompression.Type `mapstructure:"compression_algorithms"`

	// Cors corresponds to the JSON schema field "cors".
	Cors *Cors `mapstructure:"cors"`

	// Endpoint corresponds to the JSON schema field "endpoint".
	Endpoint string `mapstructure:"endpoint"`

	// IdleTimeout corresponds to the JSON schema field "idle_timeout".
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`

	// IncludeMetadata corresponds to the JSON schema field "include_metadata".
	IncludeMetadata bool `mapstructure:"include_metadata"`

	// MaxRequestBodySize corresponds to the JSON schema field
	// "max_request_body_size".
	MaxRequestBodySize int `mapstructure:"max_request_body_size"`

	// ReadHeaderTimeout corresponds to the JSON schema field "read_header_timeout".
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`

	// ReadTimeout corresponds to the JSON schema field "read_timeout".
	ReadTimeout time.Duration `mapstructure:"read_timeout"`

	// ResponseHeaders corresponds to the JSON schema field "response_headers".
	ResponseHeaders map[string]configopaque.String `mapstructure:"response_headers"`

	// WriteTimeout corresponds to the JSON schema field "write_timeout".
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type ServerConfigAuth struct {
	// RequestParams corresponds to the JSON schema field "request_params".
	RequestParams []string `mapstructure:"request_params"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ServerConfig) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ServerConfig
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["compression_algorithms"]; !ok || v == nil {
		plain.CompressionAlgorithms = "none"
	}
	if v, ok := raw["endpoint"]; !ok || v == nil {
		plain.Endpoint = ""
	}
	if v, ok := raw["idle_timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field idle_timeout:%w }", err)
		}
		plain.IdleTimeout = defaultDuration
	}
	if v, ok := raw["include_metadata"]; !ok || v == nil {
		plain.IncludeMetadata = false
	}
	if v, ok := raw["max_request_body_size"]; !ok || v == nil {
		plain.MaxRequestBodySize = 0.0
	}
	if v, ok := raw["read_header_timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field read_header_timeout:%w }", err)
		}
		plain.ReadHeaderTimeout = defaultDuration
	}
	if v, ok := raw["read_timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field read_timeout:%w }", err)
		}
		plain.ReadTimeout = defaultDuration
	}
	if v, ok := raw["write_timeout"]; !ok || v == nil {
		defaultDuration, err := time.ParseDuration("33.3s")
		if err != nil {
			return fmt.Errorf("failed to parse the \"33.3s\" default value for field write_timeout:%w }", err)
		}
		plain.WriteTimeout = defaultDuration
	}
	*j = ServerConfig(plain)
	return nil
}

// SetDefaults sets the fields of ServerConfig to their defaults.
// Fields which do not have a default value are left untouched.
func (c *ServerConfig) SetDefaults() {
	c.CompressionAlgorithms = "none"
	c.Endpoint = ""
	c.IdleTimeout = "PT33.3S"
	c.IncludeMetadata = false
	c.MaxRequestBodySize = 0.0
	c.ReadHeaderTimeout = "PT33.3S"
	c.ReadTimeout = "PT33.3S"
	c.WriteTimeout = "PT33.3S"
}
