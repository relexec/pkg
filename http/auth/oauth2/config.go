package oauth2

import "github.com/relexec/pkg/http/auth/oauth2/jose"

// Config contains configuration for the authn/authz setup when using the
// OAuth2 authorization framework over the Bearer HTTP authentication scheme.
type Config struct {
	// AuthorizationURL is the URL for the third-party OAuth2 service that
	// authorizes requests, e.g "https://my-tenant.auth0.com/authorize"
	AuthorizationURL string `json:"authorization_url,omitempty"`
	// TokenURL is the URL for the third-party OAuth2 service that is sent
	// authorization codes, client credentials or refresh tokens and vends
	// access tokens, e.g "https://my-tenant.auth0.com/oauth/token"
	TokenURL string `json:"token_url,omitempty"`
	// JOSE  contains configuration fields for Javascript Object Signing and
	// Encryption (JOSE) infrastructure
	JOSE jose.Config `json:"jose"`
}
