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
//

package misc

import (
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestRandomString(t *testing.T) {
	for range 1024 {
		size := rand.Intn(512)
		str, err := RandomString(size)
		if err != nil {
			t.Errorf("RandomString returned an error for size %d: %v", size, err)
		}
		if len(str) != size {
			t.Errorf("RandomString returned a string of length %d for size %d", len(str), size)
		}
	}
}

func TestIsUnixSocketAddr(t *testing.T) {
	tests := []struct {
		addr string
		want bool
	}{
		{addr: "", want: false},
		{addr: "/var/run/haproxy.sock", want: true},
		{addr: "unix@/var/run/haproxy.sock", want: true},
		{addr: "sockpair@7", want: false},
		{addr: "fd@3", want: false},
		{addr: "ipv4@127.0.0.1:1234", want: false},
		{addr: "ipv6@::1:1234", want: false},
		{addr: "127.0.0.1:1234", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			if got := IsUnixSocketAddr(tt.addr); got != tt.want {
				t.Errorf("IsUnixSocketAddr(%q) = %v, want %v", tt.addr, got, tt.want)
			}
		})
	}
}

// listenUnix binds a UNIX socket named name inside dir and returns its path.
func listenUnix(t *testing.T, dir, name string) string {
	t.Helper()

	socket := filepath.Join(dir, name)
	l, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatalf("cannot listen on %s: %v", socket, err)
	}
	t.Cleanup(func() { l.Close() })

	return socket
}

func TestMasterSocketFromEnv(t *testing.T) {
	// os.MkdirTemp instead of t.TempDir: the latter embeds the test name in the
	// path, which easily overflows the 104 bytes of sun_path on some systems.
	dir, err := os.MkdirTemp("", "dpapi")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	bound := listenUnix(t, dir, "master.sock")
	second := listenUnix(t, dir, "second.sock")
	missing := filepath.Join(dir, "missing.sock")
	regular := filepath.Join(dir, "regular")
	if err := os.WriteFile(regular, nil, 0o600); err != nil {
		t.Fatalf("cannot create regular file: %v", err)
	}

	tests := []struct {
		name   string
		value  string
		want   string
		wantOK bool
	}{
		{
			name:   "empty value",
			value:  "",
			want:   "",
			wantOK: false,
		},
		{
			name:   "only a sockpair",
			value:  "sockpair@7",
			want:   "",
			wantOK: false,
		},
		{
			name:   "unix socket followed by a sockpair",
			value:  "unix@" + bound + ";sockpair@7",
			want:   bound,
			wantOK: true,
		},
		{
			name:   "sockpair listed first",
			value:  "sockpair@7;unix@" + bound,
			want:   bound,
			wantOK: true,
		},
		{
			name:   "first bound socket wins",
			value:  "unix@" + missing + ";unix@" + second,
			want:   second,
			wantOK: true,
		},
		{
			name:   "a regular file is not a socket",
			value:  "unix@" + regular + ";unix@" + bound,
			want:   bound,
			wantOK: true,
		},
		{
			name:   "nothing bound yet falls back to the first candidate",
			value:  "unix@" + missing + ";sockpair@7",
			want:   missing,
			wantOK: true,
		},
		{
			name:   "bare path without the unix prefix",
			value:  bound,
			want:   bound,
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := MasterSocketFromEnv(tt.value)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("MasterSocketFromEnv(%q) = (%q, %v), want (%q, %v)", tt.value, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
