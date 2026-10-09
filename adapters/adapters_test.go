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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// selfSignedCert returns one ECDSA P-256 certificate usable as CA, server and
// client certificate, and a pool containing it.
func selfSignedCert(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}, pool
}

// TestBasicAuthMiddlewareSharedHandlerListeners reproduces the report: one
// handler is shared by a plain HTTP listener and an mTLS HTTPS listener, so
// Basic auth may be skipped only for requests carrying a verified client cert.
func TestBasicAuthMiddlewareSharedHandlerListeners(t *testing.T) {
	cert, pool := selfSignedCert(t)
	h := BasicAuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	plain := httptest.NewServer(h)
	t.Cleanup(plain.Close)
	mtls := httptest.NewUnstartedServer(h)
	mtls.TLS = &tls.Config{Certificates: []tls.Certificate{cert}, ClientCAs: pool, ClientAuth: tls.RequireAndVerifyClientCert}
	mtls.StartTLS()
	t.Cleanup(mtls.Close)
	mtlsClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool, Certificates: []tls.Certificate{cert}}}}

	const path = "/v3/services/haproxy/configuration/raw"
	tests := []struct {
		name     string
		url      string
		client   *http.Client
		withAuth bool
		want     int
	}{
		{name: "plain HTTP, no credentials", url: plain.URL, client: plain.Client(), want: http.StatusUnauthorized},
		{name: "plain HTTP, wrong credentials", url: plain.URL, client: plain.Client(), withAuth: true, want: http.StatusUnauthorized},
		{name: "mTLS HTTPS with verified client certificate, no credentials", url: mtls.URL, client: mtlsClient, want: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, tt.url+path, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			if tt.withAuth {
				req.SetBasicAuth("nosuchuser", "nosuchpass")
			}
			resp, err := tt.client.Do(req)
			if err != nil {
				t.Fatalf("request: %v", err)
			}
			resp.Body.Close()
			if resp.StatusCode != tt.want {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.want)
			}
		})
	}
}

func TestBasicAuthMiddlewareMTLS(t *testing.T) {
	tests := []struct {
		name       string
		tls        *tls.ConnectionState
		wantCalled bool
		wantCode   int
	}{
		{
			name:       "verified client certificate",
			tls:        &tls.ConnectionState{VerifiedChains: [][]*x509.Certificate{{&x509.Certificate{}}}},
			wantCalled: true,
			wantCode:   http.StatusOK,
		},
		{
			name:     "TLS without verified client certificate",
			tls:      &tls.ConnectionState{},
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := BasicAuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/v3/info", nil)
			req.TLS = tt.tls
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if called != tt.wantCalled {
				t.Errorf("handler called = %v, want %v", called, tt.wantCalled)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantCode)
			}
		})
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
			h := BasicAuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
