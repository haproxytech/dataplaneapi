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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCrossOriginProtectionMiddleware(t *testing.T) {
	h := CrossOriginProtectionMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	for _, tc := range []struct {
		name, method, fetchSite string
		want                    int
	}{
		{"cross-site form post (CSRF)", http.MethodPost, "cross-site", http.StatusForbidden},
		{"same-origin post", http.MethodPost, "same-origin", http.StatusAccepted},
		{"non-browser client post", http.MethodPost, "", http.StatusAccepted},
		{"cross-site get", http.MethodGet, "cross-site", http.StatusAccepted},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/v3/services/haproxy/configuration/raw", strings.NewReader("# x=y\r\n"))
			req.Header.Set("Content-Type", "text/plain")
			if tc.fetchSite != "" {
				req.Header.Set("Sec-Fetch-Site", tc.fetchSite)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("status = %d, want %d", rec.Code, tc.want)
			}
		})
	}
}
