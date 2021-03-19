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
load '../../libs/get_json_path'
load '../../libs/version'

@test "frontends: Add a frontend" {
	read -r SC RES < <(auth_curl POST "/v2/services/haproxy/configuration/frontends?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/frontends/post.json")
	[ "${SC}" = 201 ]
}

@test "frontends: Return a frontend" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/configuration/frontends/test_frontend")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data.name')
	[ "${ACTUAL}" = "test_frontend" ]
}

@test "frontends: Replace a frontend" {
	read -r SC RES < <(auth_curl PUT "/v2/services/haproxy/configuration/frontends/test_frontend?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/frontends/put.json")
	[ "${SC}" = 200 ]
}

@test "frontends: Return an array of frontends" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/configuration/frontends")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data[0].name')
	[ "${ACTUAL}" = "test_frontend" ]
}

@test "frontends: Delete a frontend" {
	read -r SC RES < <(auth_curl DELETE "/v2/services/haproxy/configuration/frontends/test_frontend?force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
}
