package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/relexec/pkg/http/auth"
	"github.com/relexec/pkg/http/request"
)

// Authorized returns middleware that ensures there is an authenticated
// identity making the HTTP request and that they are authorized to perform the
// operation.
func Authorized(
	api huma.API,
	authz request.Authorizer,
) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		authorized, err := authz.Authorize(ctx)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, err.Error())
			return
		}

		ctx = huma.WithValue(
			ctx,
			request.ContextKeyRequestAuth,
			authorized,
		)
		rctx := ctx.Context()
		rctx = auth.IdentityToContext(rctx, authorized.Identity)
		ctx = huma.WithContext(ctx, rctx)

		next(ctx)
	}
}
