package auth

import "context"

type contextKey string

const (
	contextKeyIdentity = "relexec.identity"
)

// IdentityToContext sets the supplied Identity into the supplied context and
// returns the adapted context.
func IdentityToContext(ctx context.Context, identity string) context.Context {
	return context.WithValue(ctx, contextKeyIdentity, identity)
}

// IdentityFromContext returns the Identity contained in the supplied context.
func IdentityFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(contextKeyIdentity); v != nil {
		return v.(string)
	}
	return ""
}
