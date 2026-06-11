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

package respond

import (
	stdjson "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type errorBody struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func decodeErrorBody(t *testing.T, rec *httptest.ResponseRecorder) errorBody {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var e errorBody
	if err := stdjson.Unmarshal(rec.Body.Bytes(), &e); err != nil {
		t.Fatalf("body %q is not valid JSON: %v", rec.Body.String(), err)
	}
	return e
}

func TestNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	NotFound(rec, httptest.NewRequest(http.MethodGet, "/v3/nonexistent", nil))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	e := decodeErrorBody(t, rec)
	if e.Code != http.StatusNotFound {
		t.Errorf("body code = %d, want 404", e.Code)
	}
	if want := "path /v3/nonexistent was not found"; e.Message != want {
		t.Errorf("body message = %q, want %q", e.Message, want)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	ok := func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }
	inner := chi.NewRouter()
	inner.Get("/info", ok)
	inner.Put("/info", ok)
	inner.MethodNotAllowed(MethodNotAllowed(inner))
	root := chi.NewRouter()
	root.Mount("/v3", inner)

	rec := httptest.NewRecorder()
	root.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/v3/info", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != "GET, PUT" {
		t.Errorf("Allow = %q, want %q", allow, "GET, PUT")
	}
	e := decodeErrorBody(t, rec)
	if e.Code != http.StatusMethodNotAllowed {
		t.Errorf("body code = %d, want 405", e.Code)
	}
	if want := "method PATCH is not allowed for path /v3/info"; e.Message != want {
		t.Errorf("body message = %q, want %q", e.Message, want)
	}
}

func TestJSONNilPayload(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusOK, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "" {
		t.Errorf("Content-Type = %q, want unset", ct)
	}
}

func TestJSONEncodeFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusOK, map[string]any{"bad": func() {}})

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	e := decodeErrorBody(t, rec)
	if e.Code != http.StatusInternalServerError {
		t.Errorf("body code = %d, want 500", e.Code)
	}
	if e.Message == "" {
		t.Error("body message is empty")
	}
}
