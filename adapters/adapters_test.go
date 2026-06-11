// Copyright 2026 HAProxy Technologies
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
//

package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthMiddlewareSkip(t *testing.T) {
	called := false
	h := BasicAuthMiddleware(true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v3/info", nil))

	if !called {
		t.Error("handler was not called with skip=true")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestBasicAuthMiddlewareUnauthorized(t *testing.T) {
	// The global user store is empty in this test process, so any credentials
	// are rejected; both paths must produce the same 401 response shape.
	tests := []struct {
		name     string
		withAuth bool
	}{
		{name: "missing authorization header", withAuth: false},
		{name: "credentials not in user store", withAuth: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := BasicAuthMiddleware(false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			}))

			req := httptest.NewRequest(http.MethodGet, "/v3/info", nil)
			if tt.withAuth {
				req.SetBasicAuth("nosuchuser", "nosuchpass")
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if called {
				t.Error("handler was called despite failed authentication")
			}
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want 401", rec.Code)
			}
			if got, want := rec.Header().Get("WWW-Authenticate"), `Basic realm="API"`; got != want {
				t.Errorf("WWW-Authenticate = %q, want %q", got, want)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			var e struct {
				Code    int64  `json:"code"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil {
				t.Fatalf("body %q is not valid JSON: %v", rec.Body.String(), err)
			}
			if e.Code != http.StatusUnauthorized {
				t.Errorf("body code = %d, want 401", e.Code)
			}
			if e.Message == "" {
				t.Error("body message is empty")
			}
		})
	}
}
