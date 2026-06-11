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

package respond

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	jsoniter "github.com/json-iterator/go"

	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ReloadIDHeader is the response header carrying the reload identifier for 202
// responses. The spec spells it "Reload-ID", but header names are
// case-insensitive and http.Header.Set canonicalizes any spelling to
// "Reload-Id" on the wire, so the constant uses the canonical form directly.
const ReloadIDHeader = "Reload-Id"

// JSON writes status and payload as JSON. A nil payload produces an empty body
// with no Content-Type, matching the previous go-swagger responses.
// The payload is encoded before any byte is written, so a marshaling failure
// still produces a 500 instead of a truncated body with a success status.
func JSON(w http.ResponseWriter, status int, payload any) {
	if payload == nil {
		w.WriteHeader(status)
		return
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		log.Errorf("encoding %T response payload: %v", payload, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":500,"message":"failed to encode response payload"}` + "\n"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Errorf("writing response body: %v", err)
	}
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Write writes a raw (non-JSON) body after the header has been written,
// logging the error if the write fails. The status can no longer change at
// that point, so logging is the only remaining evidence of a truncated body.
func Write(w http.ResponseWriter, b []byte) {
	if _, err := w.Write(b); err != nil {
		log.Errorf("writing response body: %v", err)
	}
}

// Copy streams src to w after the header has been written, logging the error
// if the copy fails. See Write.
func Copy(w http.ResponseWriter, src io.Reader) {
	if _, err := io.Copy(w, src); err != nil {
		log.Errorf("writing response body: %v", err)
	}
}

// Accepted writes a 202 Accepted response with a Reload-ID header and JSON payload.
// A nil payload produces an empty body, matching the previous go-swagger responses.
func Accepted(w http.ResponseWriter, reloadID string, payload any) {
	if reloadID != "" {
		w.Header().Set(ReloadIDHeader, reloadID)
	}
	JSON(w, http.StatusAccepted, payload)
}

// Error translates err to a models.Error via misc.HandleError and writes the appropriate status.
func Error(w http.ResponseWriter, err error) {
	e := misc.HandleError(err)
	JSON(w, int(*e.Code), e)
}

// RuntimeError translates err using misc.GetHTTPStatusFromErr (which maps
// client_errors.ErrGeneral to 400) and writes the resulting status. Runtime and
// storage operation handlers use this to preserve the pre-migration status codes,
// where general runtime/storage failures return 400 rather than 500.
func RuntimeError(w http.ResponseWriter, err error) {
	status := misc.GetHTTPStatusFromErr(err)
	code := int64(status)
	msg := err.Error()
	JSON(w, status, &models.Error{Code: &code, Message: &msg})
}

// MultipartError writes the error from parsing a multipart upload body. With
// body validation excluded for multipart operations, the first read of the
// body happens in the handler's ParseMultipartForm call, so an upload that
// hit the server's max-body-size cap surfaces here as *http.MaxBytesError and
// maps to 413, matching the validator's handling of buffered bodies. Any
// other parse failure is a plain 400.
func MultipartError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		code := int64(http.StatusRequestEntityTooLarge)
		msg := fmt.Sprintf("request body exceeds the maximum allowed size of %d bytes", maxBytesErr.Limit)
		JSON(w, http.StatusRequestEntityTooLarge, &models.Error{Code: &code, Message: &msg})
		return
	}
	BadRequest(w, err.Error())
}

// BadRequest writes a 400 response with msg as the error message.
func BadRequest(w http.ResponseWriter, msg string) {
	code := misc.ErrHTTPBadRequest
	JSON(w, http.StatusBadRequest, &models.Error{Code: &code, Message: &msg})
}

// NotFound is the router's chi.NotFound handler. It replaces chi's default
// plain-text "404 page not found" body with the JSON shape the previous
// go-swagger server returned for unmatched routes:
// {"code":404,"message":"path <path> was not found"}.
func NotFound(w http.ResponseWriter, r *http.Request) {
	code := misc.ErrHTTPNotFound
	msg := "path " + r.URL.Path + " was not found"
	JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
}

// MethodNotAllowed returns the router's chi.MethodNotAllowed handler. It
// replaces chi's default plain-text "405 method not allowed" body with a JSON
// error and sets the Allow header RFC 9110 requires on 405 responses. chi does
// not expose the allowed-method list to custom 405 handlers, so the handler
// probes routes for each standard method; routes must be the router the
// handler is registered on, since the probe path is relative to it.
func MethodNotAllowed(routes chi.Routes) http.HandlerFunc {
	methods := []string{
		http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete, http.MethodOptions,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Inside a mounted sub-router the registered patterns do not carry the
		// mount prefix; chi records the trimmed path in RoutePath.
		path := r.URL.Path
		if rctx := chi.RouteContext(r.Context()); rctx != nil && rctx.RoutePath != "" {
			path = rctx.RoutePath
		}
		var allowed []string
		for _, m := range methods {
			if routes.Match(chi.NewRouteContext(), m, path) {
				allowed = append(allowed, m)
			}
		}
		if len(allowed) > 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
		}
		code := int64(http.StatusMethodNotAllowed)
		msg := "method " + r.Method + " is not allowed for path " + r.URL.Path
		JSON(w, http.StatusMethodNotAllowed, &models.Error{Code: &code, Message: &msg})
	}
}

// DecodeBody decodes JSON from r.Body into out then calls out.Validate(strfmt.Default).
// On any failure it writes the appropriate error response and returns false.
// Pass a pointer to the target value: respond.DecodeBody(r, w, &data).
// All write handlers (POST/PUT) must use this instead of raw json.Decode;
// generated inline body structs without a Validate method use DecodeJSON.
func DecodeBody[T interface {
	Validate(registry strfmt.Registry) error
}](r *http.Request, w http.ResponseWriter, out T) bool {
	if !DecodeJSON(r, w, out) {
		return false
	}
	if err := out.Validate(strfmt.Default); err != nil {
		Unprocessable(w, err)
		return false
	}
	return true
}

// DecodeJSON decodes JSON from r.Body into out, writing a 400 response and
// returning false on failure. It is DecodeBody for the generated inline body
// structs that have no Validate method; client-native models embedded in such
// a body should still be validated individually (see Unprocessable).
func DecodeJSON(r *http.Request, w http.ResponseWriter, out any) bool {
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		BadRequest(w, err.Error())
		return false
	}
	return true
}

// Unprocessable writes a 422 Unprocessable Entity response with a cleaned
// message from a go-openapi CompositeError. The body code matches the HTTP
// status, consistent with the rest of the API.
func Unprocessable(w http.ResponseWriter, err error) {
	code := int64(http.StatusUnprocessableEntity)
	msg := flattenValidationError(err)
	JSON(w, http.StatusUnprocessableEntity, &models.Error{Code: &code, Message: &msg})
}

// flattenValidationError strips the repetitive "validation failure list:" wrapper
// lines that go-openapi/errors CompositeError nests for each schema level, leaving
// only the leaf error messages.
func flattenValidationError(err error) string {
	const prefix = "validation failure list:"
	var keep []string
	for line := range strings.SplitSeq(err.Error(), "\n") {
		if line != prefix {
			keep = append(keep, line)
		}
	}
	if len(keep) == 0 {
		return err.Error()
	}
	return strings.Join(keep, "\n")
}
