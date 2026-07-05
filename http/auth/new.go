package auth

import (
	"context"
	"fmt"

	"github.com/relexec/pkg/http/auth/basic"
	"github.com/relexec/pkg/http/auth/oauth2/jose"
	"github.com/relexec/pkg/http/request"
)

// New returns an Authorizer configured using the supplied values.
func New(ctx context.Context, cfg Config) (request.Authorizer, error) {
	if cfg.Basic != nil {
		return basic.New(ctx, *cfg.Basic)
	}
	if cfg.OAuth2 == nil {
		return nil, fmt.Errorf(
			"either basic or oauth2 authorization must be configured",
		)
	}
	return jose.New(ctx, cfg.OAuth2.JOSE)
}
