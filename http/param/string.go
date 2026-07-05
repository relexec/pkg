package param

import (
	"fmt"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

// String returns a URL-unescaped string parameter value for a supplied query
// parameter key. If the parameter key does not exist in the request context or
// is an empty string, returns an error.
func String(hctx huma.Context, key string) (string, error) {
	param, err := url.QueryUnescape(hctx.Param(key))
	if err != nil {
		return "", fmt.Errorf("could not parse %q query param: %w", key, err)
	}
	if param == "" {
		return "", fmt.Errorf("empty %q query param", key)
	}
	return param, nil
}
