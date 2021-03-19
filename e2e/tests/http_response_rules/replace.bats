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

setup() {
	# creating frontend and related HTTP Response rule
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/frontends?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/frontends/post.json")
	[ "${SC}" = 201 ]
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/http_response_rules?parent_type=frontend&parent_name=test_frontend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/http_response_rules/unless.json")
	[ "${SC}" = 201 ]
	# creating backend and related HTTP Response rule
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/backends/post.json")
	[ "${SC}" = 201 ]
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/http_response_rules?parent_type=backend&parent_name=test_backend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/http_response_rules/unless.json")
	[ "${SC}" = 201 ]
}

teardown() {
	read -r SC _ < <(auth_curl DELETE "/v2/services/haproxy/configuration/frontends/test_frontend?force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
	read -r SC _ < <(auth_curl DELETE "/v2/services/haproxy/configuration/backends/test_backend?force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
}

@test "http_response_rules: Replace a HTTP Response Rule of frontend" {
	read -r SC BODY < <(auth_curl PUT "/v2/services/haproxy/configuration/http_response_rules/0?parent_type=frontend&parent_name=test_frontend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/http_response_rules/put.json")
	[ "${SC}" = 200 ]
	[ "$(get_json_path "${BODY}" ".cond")" = "if" ]
}

@test "http_response_rules: Replace a HTTP Response Rule of backend" {
	read -r SC BODY < <(auth_curl PUT "/v2/services/haproxy/configuration/http_response_rules/0?parent_type=backend&parent_name=test_backend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/http_response_rules/put.json")
	[ "${SC}" = 200 ]
	[ "$(get_json_path "${BODY}" ".cond")" = "if" ]
}
