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

package runtime

import (
	"net"
	"testing"

	"github.com/haproxytech/dataplaneapi/log"
)

// Whitespace-only input on the debug socket must not panic: serve runs in a
// goroutine without recover, so a panic kills the process.
func TestServe_EmptyCommand(t *testing.T) {
	if _, err := log.AppLogger(); err != nil {
		t.Fatal(err)
	}
	server, client := net.Pipe()
	go func() {
		_, _ = client.Write([]byte("\n"))
		_, _ = client.Read(make([]byte, 64))
		client.Close()
	}()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("serve panicked on empty input: %v", r)
		}
	}()
	serve(&Commands{}, server)
}
