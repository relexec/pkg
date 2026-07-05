package auth

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/relexec/pkg/http/auth/basic"
	"github.com/relexec/pkg/http/auth/oauth2"
)

const (
	envVarDevelopment = "DEVELOPMENT"
)

// Config contains configuration for the authentication/authorization
// infrastructure used in the server.
type Config struct {
	// Basic contains configuration for the HTTP Basic Authorizer, which should
	// only be used for local developer testing.
	Basic *basic.Config `json:"basic,omitempty"`
	// OAuth2 contains configuration for the OAuth2 Authorizer.
	OAuth2 *oauth2.Config `json:"oauth2,omitempty"`
}

// SetDefaults sets any missing values to their defaults or environs variable
// values.
func (c *Config) SetDefaults() {
	// TODO(jaypipes): Remove this when able to easily set up development
	// authz.
	dev := os.Getenv(envVarDevelopment)
	if dev != "" {
		a := basic.Development()
		c.Basic = &a
	}
	return
}

// BindFlags binds the supplied flagset to the Config's fields.
func (c *Config) BindFlags(fs *pflag.FlagSet) {
	return
}

// Validate checks for invalid settings.
func (c Config) Validate() error {
	if c.Basic == nil && c.OAuth2 == nil {
		return fmt.Errorf(
			"either basic or oauth2 authorization must be configured",
		)
	}
	return nil
}

// ScopedSecurity returns the OpenAPI Operation.Security information scoped to
// the supplied permissions.
func (c Config) ScopedSecurity(scopes ...string) []map[string][]string {
	switch {
	case c.Basic != nil:
		return []map[string][]string{
			{
				"basic": scopes,
			},
		}
	case c.OAuth2 != nil:
		return []map[string][]string{
			{
				"oauth2": scopes,
			},
		}
	default:
		return nil
	}
}
