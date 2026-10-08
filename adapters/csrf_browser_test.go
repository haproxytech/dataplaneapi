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
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/rs/cors"
)

// Cross-site form POSTs to state-changing endpoints must be rejected even when
// the browser holds cached basic-auth credentials for the API origin.
// Set CHROME to a chromium binary to run; everything stays on loopback.
func TestCrossSiteTextPlainPostRejected(t *testing.T) {
	chrome := os.Getenv("CHROME")
	if chrome == "" {
		t.Skip("CHROME not set")
	}

	apiLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	siteLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	apiPort := apiLn.Addr().(*net.TCPAddr).Port
	sitePort := siteLn.Addr().(*net.TCPAddr).Port
	api := "http://127.0.0.1:" + strconv.Itoa(apiPort)
	site := "http://localhost:" + strconv.Itoa(sitePort) // different site than 127.0.0.1

	type hit struct{ fetchSite, ctype, body string }
	hits := make(chan hit, 1)

	// stub for the router: records what passes the middleware chain
	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v3/info": // first, authenticated top-level visit by the admin
			http.Redirect(w, r, site+"/", http.StatusFound)
		case "/v3/services/haproxy/configuration/raw":
			b, _ := io.ReadAll(r.Body)
			hits <- hit{r.Header.Get("Sec-Fetch-Site"), r.Header.Get("Content-Type"), string(b)}
			w.WriteHeader(http.StatusAccepted)
		}
	})
	// stands in for BasicAuthMiddleware, which caches credentials in the browser
	auth := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || user != "admin" || pass != "dummy" {
				w.Header().Set("WWW-Authenticate", `Basic realm="HAProxy Data Plane API"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	// same chain order and CORS options as configureAPI
	h := auth(stub)
	h = CrossOriginProtectionMiddleware()(h)
	h = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(h)
	logged := h
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, hasAuth := r.BasicAuth()
		t.Logf("api: %s %s auth=%v origin=%q sec-fetch-site=%q", r.Method, r.URL.Path, hasAuth, r.Header.Get("Origin"), r.Header.Get("Sec-Fetch-Site"))
		logged.ServeHTTP(w, r)
	})
	go func() { _ = http.Serve(apiLn, h) }()

	form := `<form method="POST" enctype="text/plain" action="` + api +
		`/v3/services/haproxy/configuration/raw?skip_version=true">` +
		`<input type="hidden" name="# x" value="y"></form><script>document.forms[0].submit()</script>`
	go func() {
		_ = http.Serve(siteLn, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("site: %s %s", r.Method, r.URL.Path)
			w.Header().Set("Content-Type", "text/html")
			_, _ = io.WriteString(w, form)
		}))
	}()

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, chrome, "--no-sandbox", "--no-first-run",
		"--user-data-dir="+t.TempDir(), "--proxy-server=direct://", "--virtual-time-budget=10000", "--dump-dom",
		"http://admin:dummy@127.0.0.1:"+strconv.Itoa(apiPort)+"/v3/info")
	cmd.Env = []string{"HOME=" + t.TempDir()}
	out, _ := cmd.CombinedOutput()
	t.Logf("browser exited, last output: %.300s", string(out[max(0, len(out)-300):]))

	select {
	case got := <-hits:
		t.Fatalf("cross-site form POST reached the handler authenticated: Sec-Fetch-Site=%q Content-Type=%q body=%q",
			got.fetchSite, got.ctype, got.body)
	default:
		// not reached: rejected (or the browser did not attach credentials)
	}
}
