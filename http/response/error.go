package response

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// InternalServerError writes the supplied error to the HTTP response, setting
// a 500 Internal Server Error status code in the response.
func InternalServerError(
	api huma.API,
	hctx huma.Context,
	err error,
) {
	_ = huma.WriteErr(
		api, hctx, http.StatusInternalServerError,
		"server error", err,
	)
}

// NotFound writes a standardized error to the HTTP response, setting a 404
// NotFound status code in the response.
func NotFound(
	api huma.API,
	hctx huma.Context,
) {
	_ = huma.WriteErr(
		api, hctx, http.StatusNotFound,
		"not found", fmt.Errorf("not found"),
	)
}

// BadRequest writes the supplied error to the HTTP response, setting a 400
// BadRequest status code in the response.
func BadRequest(
	api huma.API,
	hctx huma.Context,
	err error,
) {
	_ = huma.WriteErr(
		api, hctx, http.StatusBadRequest,
		"bad request", err,
	)
}
