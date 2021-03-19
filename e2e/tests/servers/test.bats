#!/usr/bin/env bats
#
# Copyright 2019 HAProxy Technologies
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

load '../../libs/auth_curl'
load "../../libs/get_json_path"
load '../../libs/version'

@test "servers: setup" {
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/backends/post.json")
	[ "${SC}" = 201 ]
}

@test "servers: Add a new server" {
	read -r SC BODY < <(auth_curl POST "/v2/services/haproxy/configuration/servers?backend=test_backend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/servers/post.json")
	[ "${SC}" = 201 ]
}

@test "servers: Return one server" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/configuration/servers/test_server?backend=test_backend")
	[ "${SC}" = 200 ]

	local NAME; NAME=$(get_json_path "$BODY" '.data.name')
	[ "${NAME}" = "test_server" ]
}

@test "servers: Replace a server" {
	read -r SC BODY < <(auth_curl PUT "/v2/services/haproxy/configuration/servers/test_server?backend=test_backend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/servers/put.json")
	[ "${SC}" = 200 ]
}

@test "servers: Return an array of servers" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/configuration/servers?backend=test_backend")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data[0].name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Return an array of runtime servers' settings" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/runtime/servers?backend=test_backend")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.[0].name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Replace server transient settings" {
	read -r SC BODY < <(auth_curl PUT "/v2/services/haproxy/runtime/servers/test_server?backend=test_backend" "@${E2E_DIR}/tests/servers/transient.json")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Return one server runtime settings" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/runtime/servers/test_server?backend=test_backend")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Delete a server" {
	read -r SC BODY < <(auth_curl DELETE "/v2/services/haproxy/configuration/servers/test_server?backend=test_backend&force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
}

@test "servers: teardown" {
	read -r SC < <(auth_curl DELETE "/v2/services/haproxy/configuration/backends/test_backend?force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
}
