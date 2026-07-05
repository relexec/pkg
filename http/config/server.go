package httpconfig

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

const (
	DefaultBindHost              = "0.0.0.0"
	envVarHost                   = "BIND_HOST"
	DefaultBindPort              = 8888
	envVarPort                   = "BIND_PORT"
	DefaultHTTPReadTimeout       = 5 * time.Second
	DefaultHTTPReadHeaderTimeout = 1 * time.Second
	DefaultHTTPWriteTimeout      = 5 * time.Minute
	DefaultHTTPIdleTimeout       = 30 * time.Second
)

// ServerConfig contains configuration for an HTTP server.
type ServerConfig struct {
	// Host is the host address to listen on. Defaults to "0.0.0.0".
	Host string `json:"host,omitempty"`
	// Port is the host port to listen on.
	Port int `json:"port,omitzero"`
	// ReadTimeout is the timeout duration for reads. You can use Go duration
	// string formats, e.g. "5s".
	ReadTimeout         string `json:"read_timeout,omitempty"`
	readTimeoutDuration time.Duration
	// ReadHeaderTimeout is the timeout duration for reading headers. You can
	// use Go duration string formats, e.g. "1s".
	ReadHeaderTimeout         string `json:"read_header_timeout,omitempty"`
	readHeaderTimeoutDuration time.Duration
	// WriteTimeout is the timeout duration for writes. You can use Go duration
	// string formats, e.g. "5m".
	WriteTimeout         string `json:"write_timeout,omitempty"`
	writeTimeoutDuration time.Duration
	// IdleTimeout is the timeout duration for idle connections. You can use Go
	// duration string formats, e.g. "30s".
	IdleTimeout         string `json:"idle_timeout,omitempty"`
	idleTimeoutDuration time.Duration
}

// Validate checks for invalid settings.
func (c ServerConfig) Validate() error {
	if c.ReadTimeout != "" {
		dur, err := time.ParseDuration(c.ReadTimeout)
		if err != nil {
			return fmt.Errorf("failed parsing http.read_timeout: %w", err)
		}
		c.readTimeoutDuration = dur
	}
	if c.ReadHeaderTimeout != "" {
		dur, err := time.ParseDuration(c.ReadHeaderTimeout)
		if err != nil {
			return fmt.Errorf("failed parsing http.read_header_timeout: %w", err)
		}
		c.readHeaderTimeoutDuration = dur
	}
	if c.WriteTimeout != "" {
		dur, err := time.ParseDuration(c.WriteTimeout)
		if err != nil {
			return fmt.Errorf("failed parsing http.write_timeout: %w", err)
		}
		c.writeTimeoutDuration = dur
	}
	if c.IdleTimeout != "" {
		dur, err := time.ParseDuration(c.IdleTimeout)
		if err != nil {
			return fmt.Errorf("failed parsing http.idle_timeout: %w", err)
		}
		c.idleTimeoutDuration = dur
	}
	return nil
}

// BindFlags bings the supplied flagset to the Config's fields.
func (c *ServerConfig) BindFlags(fs *pflag.FlagSet) {
	return
}

// SetDefaults sets any missing values to their defaults or environs variable
// values.
func (c *ServerConfig) SetDefaults() {
	if c.Host == "" {
		v := os.Getenv(envVarHost)
		if v == "" {
			v = DefaultBindHost
		}
		c.Host = v
	}
	if c.Port == 0 {
		var err error
		port := DefaultBindPort
		v := os.Getenv(envVarPort)
		if v != "" {
			port, err = strconv.Atoi(v)
			if err != nil {
				port = DefaultBindPort
			}
		}
		c.Port = port
	}
	if c.ReadTimeout == "" {
		c.readTimeoutDuration = DefaultHTTPReadTimeout
	}
	if c.ReadHeaderTimeout == "" {
		c.readHeaderTimeoutDuration = DefaultHTTPReadHeaderTimeout
	}
	if c.WriteTimeout == "" {
		c.writeTimeoutDuration = DefaultHTTPWriteTimeout
	}
	if c.IdleTimeout == "" {
		c.idleTimeoutDuration = DefaultHTTPIdleTimeout
	}
}

// ReadTimeoutDuration returns the valid time.Duration for HTTP read timeouts.
func (c ServerConfig) ReadTimeoutDuration() time.Duration {
	return c.readTimeoutDuration
}

// ReadHeaderTimeoutDuration returns the valid time.Duration for HTTP header
// read timeouts.
func (c ServerConfig) ReadHeaderTimeoutDuration() time.Duration {
	return c.readHeaderTimeoutDuration
}

// WriteTimeoutDuration returns the valid time.Duration for HTTP write
// timeouts.
func (c ServerConfig) WriteTimeoutDuration() time.Duration {
	return c.writeTimeoutDuration
}

// IdleTimeoutDuration returns the valid time.Duration for HTTP idle connection
// timeouts.
func (c ServerConfig) IdleTimeoutDuration() time.Duration {
	return c.idleTimeoutDuration
}
