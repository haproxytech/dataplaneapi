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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

const testSpec = `
openapi: "3.0.1"
info:
  title: test
  version: "1"
servers:
  - url: /v3
paths:
  /things:
    get:
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name:
                  type: string
      responses:
        '201':
          description: created
          content:
            application/json:
              schema:
                type: object
  /things/{id}:
    get:
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                type: object
`

// newTestHandler mirrors the production setup: the validator wraps each
// handler on a subrouter registered with the spec's path strings, mounted
// under the spec's server URL (handlers/router.go mounts everything at /v3).
func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromData([]byte(testSpec))
	if err != nil {
		t.Fatal(err)
	}
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := NewValidator(spec)
	sub := chi.NewRouter()
	sub.Method(http.MethodGet, "/things", mw(ok))
	sub.Method(http.MethodPost, "/things", mw(ok))
	sub.Method(http.MethodGet, "/things/{id}", mw(ok))
	root := chi.NewRouter()
	root.Mount("/v3", sub)
	return root
}

func TestValidator(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		name        string
		method      string
		target      string
		body        string
		contentType string
		accept      string
		wantStatus  int
	}{
		{name: "valid GET", method: "GET", target: "/v3/things", wantStatus: 200},
		{name: "valid GET with Accept json", method: "GET", target: "/v3/things", accept: "application/json", wantStatus: 200},
		{name: "valid GET with Accept wildcard", method: "GET", target: "/v3/things", accept: "*/*", wantStatus: 200},
		{name: "valid GET with Accept type wildcard", method: "GET", target: "/v3/things", accept: "application/*", wantStatus: 200},
		{name: "Accept mismatch returns 406", method: "GET", target: "/v3/things", accept: "text/plain", wantStatus: 406},
		{name: "Accept list with one match passes", method: "GET", target: "/v3/things", accept: "text/plain, application/json;q=0.8", wantStatus: 200},
		{name: "invalid query param returns 400", method: "GET", target: "/v3/things?limit=0", wantStatus: 400},
		{name: "valid path param", method: "GET", target: "/v3/things/12", wantStatus: 200},
		{name: "invalid path param returns 400", method: "GET", target: "/v3/things/abc", wantStatus: 400},
		{name: "valid POST", method: "POST", target: "/v3/things", body: `{"name":"a"}`, contentType: "application/json", wantStatus: 200},
		{name: "wrong Content-Type returns 415", method: "POST", target: "/v3/things", body: `name=a`, contentType: "text/plain", wantStatus: 415},
		{name: "body schema violation returns 422", method: "POST", target: "/v3/things", body: `{}`, contentType: "application/json", wantStatus: 422},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *strings.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			} else {
				body = strings.NewReader("")
			}
			req := httptest.NewRequest(tt.method, tt.target, body)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

// TestValidatorBodyTooLarge verifies that a body hitting the MaxBodySize cap
// while the validator buffers it surfaces as 413 instead of a generic 400.
func TestValidatorBodyTooLarge(t *testing.T) {
	h := newTestHandler(t)
	limited := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 8)
		h.ServeHTTP(w, r)
	})
	req := httptest.NewRequest(http.MethodPost, "/v3/things", strings.NewReader(`{"name":"too long for the limit"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	limited.ServeHTTP(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("got status %d, want 413 (body: %s)", rec.Code, rec.Body.String())
	}
}
