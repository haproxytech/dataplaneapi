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

package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/runtime/security"
)

func TestBasicAuthenticator(t *testing.T) {
	leaf := &x509.Certificate{}
	verified := &tls.ConnectionState{VerifiedChains: [][]*x509.Certificate{{leaf}}}

	tests := []struct {
		name          string
		tls           *tls.ConnectionState
		user, pass    string
		wantApplies   bool
		wantPrincipal any
		wantErr       bool
		wantCalls     int
	}{
		{name: "plain HTTP without credentials"},
		{name: "plain HTTP with wrong credentials", user: "admin", pass: "wrong", wantApplies: true, wantErr: true, wantCalls: 1},
		{name: "plain HTTP with valid credentials", user: "admin", pass: "secret", wantApplies: true, wantPrincipal: "admin", wantCalls: 1},
		{name: "mTLS verified client certificate", tls: verified, wantApplies: true, wantPrincipal: leaf},
		{name: "TLS without verified client certificate", tls: &tls.ConnectionState{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			auth := basicAuthenticator(func(user, pass string) (any, error) {
				calls++
				if user == "admin" && pass == "secret" {
					return "admin", nil
				}
				return nil, errors.New("invalid credentials")
			})

			req := httptest.NewRequest(http.MethodGet, "/v2/info", nil)
			req.TLS = tt.tls
			if tt.user != "" {
				req.SetBasicAuth(tt.user, tt.pass)
			}

			applies, principal, err := auth.Authenticate(&security.ScopedAuthRequest{Request: req})
			if applies != tt.wantApplies {
				t.Errorf("applies = %v, want %v", applies, tt.wantApplies)
			}
			if principal != tt.wantPrincipal {
				t.Errorf("principal = %v, want %v", principal, tt.wantPrincipal)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, want error %v", err, tt.wantErr)
			}
			if calls != tt.wantCalls {
				t.Errorf("authentication calls = %d, want %d", calls, tt.wantCalls)
			}
		})
	}
}
