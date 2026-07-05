package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2/negotiation"
)

// HTMLDocs intercepts HTML requests on the root endpoint from a browser and
// redirects to the OpenAPI documentation published at /docs.
// the API from any origin.
func HTMLDocs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" &&
			strings.Contains(
				r.Header.Get("User-Agent"),
				"Mozilla",
			) &&
			negotiation.SelectQValueFast(
				r.Header.Get("Accept"),
				[]string{
					"text/html",
					"application/json",
					"application/cbor",
				},
			) == "text/html" {
			r.URL.Path = "/docs"
		}

		next.ServeHTTP(w, r)
	})
}
