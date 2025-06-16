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

package runtime

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/haproxytech/dataplaneapi/runtime/commands"
)

func TestServe(t *testing.T) {
	comm := &Commands{}
	comm.Register(commands.DataplaneapiVersion{Version: "test-version"})
	help, ok := comm.Get("help")
	require.True(t, ok)
	helpResponse, err := help.Command([]string{"help"})
	require.NoError(t, err)
	helpText := string(helpResponse) + "\n"

	tests := []struct {
		name     string
		input    string
		response string
	}{
		{name: "empty", input: "", response: helpText},
		{name: "newline", input: "\n", response: helpText},
		{name: "whitespace", input: " \t\r\n", response: helpText},
		{name: "help", input: "help\n", response: helpText},
		{name: "unknown", input: "unknown\n", response: helpText},
		{name: "version", input: "version\n", response: "test-version\n"},
		{name: "exit", input: "exit\n", response: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := net.Pipe()
			defer client.Close()
			require.NoError(t, client.SetDeadline(time.Now().Add(5*time.Second)))
			go serve(comm, server)

			_, err := client.Write([]byte(tt.input))
			require.NoError(t, err)
			response, err := io.ReadAll(client)
			require.NoError(t, err)
			require.Equal(t, tt.response, string(response))
		})
	}
}
