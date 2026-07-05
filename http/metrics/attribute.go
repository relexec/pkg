package metrics

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	AttributeNameStatusCode = "status.code"
)

// AttributeStatusCode returns the HTTP status code attribute KeyValue with the
// value of the supplied code.
func AttributeStatusCode(code int) attribute.KeyValue {
	return attribute.Int(AttributeNameStatusCode, code)
}

const (
	AttributeNameMethod = "method"
)

// AttributeMethod returns the HTTP method attribute KeyValue with the value of
// the supplied HTTP method.
func AttributeMethod(method string) attribute.KeyValue {
	return attribute.String(AttributeNameMethod, method)
}

const (
	AttributeNameRoute = "route"
)

// AttributeRoute returns the route attribute KeyValue with the value of the
// supplied route.
func AttributeRoute(route string) attribute.KeyValue {
	return attribute.String(AttributeNameRoute, route)
}
