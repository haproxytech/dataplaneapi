// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/go-chi/chi/v5"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// NewValidator returns a middleware that validates incoming requests against spec.
// Call it with the result of a handler package's GetSpec() function.
//
// The middleware runs inside the oapi-codegen wrapper, after chi has already
// matched the route, so the OpenAPI operation is resolved from chi's route
// context instead of a second routing pass: the generated registrations use
// the spec's path strings verbatim and chi shares OpenAPI's {param}
// placeholder syntax, making the matched chi pattern (minus the server-URL
// mount prefix) the spec path key.
//
// On top of openapi3filter request validation it restores two go-swagger
// behaviours: 406 when the Accept header matches none of the operation's
// response content types, and 415 when the request Content-Type is not
// accepted by the operation.
func NewValidator(spec *openapi3.T) func(http.Handler) http.Handler {
	var serverPrefix string
	if len(spec.Servers) > 0 {
		serverPrefix = spec.Servers[0].URL
	}
	specRoutes := make(map[string]*routers.Route)
	for path, item := range spec.Paths.Map() {
		for method, op := range item.Operations() {
			specRoutes[method+" "+path] = &routers.Route{
				Spec:      spec,
				Path:      path,
				PathItem:  item,
				Method:    method,
				Operation: op,
			}
		}
	}
	filterOptions := &openapi3filter.Options{
		ExcludeReadOnlyValidations: true,
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, pathParams := findRoute(specRoutes, serverPrefix, r)
			if route == nil {
				// Unreachable in practice: this middleware is attached per-route
				// by the generated chi registrations, so the matched chi pattern
				// is always one of this package's spec paths. Keep a 404 fallback.
				writeError(w, http.StatusNotFound, routers.ErrPathNotFound.Error())
				return
			}

			if produced := responseContentTypes(route.Operation); len(produced) > 0 && !acceptable(r.Header.Get("Accept"), produced) {
				// Same message go-swagger returned for a non-negotiable Accept header.
				writeError(w, http.StatusNotAcceptable,
					fmt.Sprintf("unsupported media type requested, only %v are available", produced))
				return
			}

			err := openapi3filter.ValidateRequest(r.Context(), &openapi3filter.RequestValidationInput{
				Request:    r,
				PathParams: pathParams,
				Route:      route,
				Options:    filterOptions,
			})
			if err != nil {
				statusCode, msg := normalizeValidationError(err, http.StatusBadRequest)
				writeError(w, statusCode, msg)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// findRoute resolves the OpenAPI operation for an already-routed request from chi's
// route context. serverPrefix (the spec's server URL, e.g. "/v3") is stripped from
// the matched pattern to obtain the spec path key. Parameter values are
// path-unescaped because chi matches on the raw (still-encoded) path when the URL
// has one, while the validator must see the decoded values.
func findRoute(specRoutes map[string]*routers.Route, serverPrefix string, r *http.Request) (*routers.Route, map[string]string) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return nil, nil
	}
	pattern := strings.TrimPrefix(rctx.RoutePattern(), serverPrefix)
	route, ok := specRoutes[r.Method+" "+pattern]
	if !ok {
		return nil, nil
	}
	pathParams := make(map[string]string, len(rctx.URLParams.Keys))
	for i, k := range rctx.URLParams.Keys {
		v := rctx.URLParams.Values[i]
		if dec, err := url.PathUnescape(v); err == nil {
			v = dec
		}
		pathParams[k] = v
	}
	return route, pathParams
}

// writeError emits the validation error JSON shape. The body code mirrors the
// HTTP status, consistent with the rest of the API.
func writeError(w http.ResponseWriter, statusCode int, msg string) {
	code := int64(statusCode)
	respond.JSON(w, statusCode, &models.Error{Code: &code, Message: &msg})
}

// responseContentTypes returns the union of content types this operation can
// produce across all of its declared responses (the OpenAPI 3 equivalent of
// the swagger 2 "produces" list).
func responseContentTypes(op *openapi3.Operation) []string {
	if op == nil || op.Responses == nil {
		return nil
	}
	var out []string
	for _, ref := range op.Responses.Map() {
		if ref == nil || ref.Value == nil {
			continue
		}
		for ct := range ref.Value.Content {
			if !slices.Contains(out, ct) {
				out = append(out, ct)
			}
		}
	}
	slices.Sort(out)
	return out
}

// acceptable reports whether the Accept header allows at least one of the
// produced content types. An absent header accepts anything; q-values are
// ignored beyond presence.
func acceptable(accept string, produced []string) bool {
	accept = strings.TrimSpace(accept)
	if accept == "" {
		return true
	}
	for r := range strings.SplitSeq(accept, ",") {
		mediaRange, _, _ := strings.Cut(r, ";")
		mediaRange = strings.ToLower(strings.TrimSpace(mediaRange))
		switch {
		case mediaRange == "*/*":
			return true
		case strings.HasSuffix(mediaRange, "/*"):
			prefix := strings.TrimSuffix(mediaRange, "*")
			for _, ct := range produced {
				if strings.HasPrefix(ct, prefix) {
					return true
				}
			}
		default:
			if slices.Contains(produced, mediaRange) {
				return true
			}
		}
	}
	return false
}

// prefixInvalidCT is the reason openapi3filter sets on a RequestError when the
// request Content-Type is not declared by the operation (see
// openapi3filter/validate_request.go). go-swagger answered these with 415.
const prefixInvalidCT = "header Content-Type has unexpected value"

// normalizeValidationError extracts a concise error message from the verbose
// openapi3filter error chain, keeping it close to the go-swagger format.
// Unsupported request content types return 415 and body schema violations
// return 422 to preserve parity with the old go-swagger behaviour; parameter
// errors return the suggested statusCode.
func normalizeValidationError(err error, statusCode int) (int, string) {
	var reqErr *openapi3filter.RequestError
	if errors.As(err, &reqErr) {
		if strings.HasPrefix(reqErr.Reason, prefixInvalidCT) {
			return http.StatusUnsupportedMediaType, reqErr.Reason
		}
		var schemaErr *openapi3.SchemaError
		if errors.As(reqErr.Err, &schemaErr) {
			if reqErr.Parameter != nil {
				return statusCode, fmt.Sprintf("%s in %s: %s", reqErr.Parameter.Name, reqErr.Parameter.In, schemaErr.Reason)
			}
			// Body schema violation — go-swagger returned 422 for these.
			return http.StatusUnprocessableEntity, schemaErr.Error()
		}
		// No SchemaError — strip the generic "request body has an error: " prefix.
		msg := strings.TrimPrefix(reqErr.Error(), "request body has an error: ")
		return statusCode, msg
	}
	return statusCode, err.Error()
}

// ErrorHandler is the shared ChiServerOptions.ErrorHandlerFunc used by all handler
// packages. It converts oapi-codegen parameter-binding errors (RequiredParamError,
// InvalidParamFormatError, etc.) into JSON 400 responses.
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	respond.BadRequest(w, err.Error())
}
