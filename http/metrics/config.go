package metrics

import (
	"os"
)

const (
	envVarDevelopment = "T2_DEVELOPMENT"
)

// Config contains the metrics configuration options for the server.
type Config struct {
	// HTTP contains configuration for the optional /metrics HTTP endpoint.
	HTTP HTTPConfig `json:"http"`
}

// SetDefaults sets any missing values to their defaults or environs variable
// values.
func (c *Config) SetDefaults() {
	// TODO(jaypipes): Remove this when able to easily set up development
	// metrics.
	dev := os.Getenv(envVarDevelopment)
	if dev != "" {
		c.HTTP.Enabled = true
	}
	return
}

// HTTPConfig contains configuration for the optional /metrics HTTP endpoint.
// Generally, this endpoint should only be configured for local development and
// testing. Production deployments should use the OTEL metrics exporter
// functionality.
type HTTPConfig struct {
	Enabled bool `json:"enabled"`
}
