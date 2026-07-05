package request

import (
	"github.com/danielgtaylor/huma/v2"
)

// Auth contains information about the authenticated identity.
type Auth struct {
	// Identity is the string identifier for the authenticated identity.
	Identity string
}

// Authorizer is the authorization admission handler for HTTP requests.
type Authorizer interface {
	// Authorize performs authorization on the supplied HTTP request and
	// returns an Auth with the validated identity. Authorize returns an error
	// if the supplied HTTP request cannot be verified for the given
	// conditions.
	Authorize(huma.Context) (*Auth, error)
}
