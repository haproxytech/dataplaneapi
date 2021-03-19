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

@test "sites: Add a site" {
	read -r SC BODY < <(auth_curl POST "/v2/services/haproxy/sites?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/sites/post.json")
	[ "${SC}" = 201 ]
}

@test "sites: Return a site" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/sites/test_site")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data.name')
	[ "${ACTUAL}" = "test_site" ]
}

@test "sites: Replace a site" {
	read -r SC BODY < <(auth_curl PUT "/v2/services/haproxy/sites/test_site?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/sites/put.json")
	[ "${SC}" = 200 ]
}

@test "sites: Return an array of sites" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/sites?force_reload=true&version=$(version)")
	[ "${SC}" = 200 ]

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data[0].name')
	[ "${ACTUAL}" = "test_site" ]
}

@test "sites: Delete a site" {
	read -r SC BODY < <(auth_curl DELETE "/v2/services/haproxy/sites/test_site?force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
}
