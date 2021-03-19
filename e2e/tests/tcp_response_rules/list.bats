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
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "@${E2E_DIR}/tests/backends/post.json")
	[ "${SC}" = 201 ]
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/tcp_response_rules?backend=test_backend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/tcp_response_rules/if.json")
	[ "${SC}" = 201 ]
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/tcp_response_rules?backend=test_backend&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/tcp_response_rules/unless.json")
	[ "${SC}" = 201 ]
}

teardown() {
	read -r SC _ < <(auth_curl DELETE "/v2/services/haproxy/configuration/backends/test_backend?force_reload=true&version=$(version)")
	[ "${SC}" = 204 ]
}

@test "tcp_response_rules: Return one TCP Response Rule from backend" {
	read -r SC BODY < <(auth_curl GET "/v2/services/haproxy/configuration/tcp_response_rules?backend=test_backend")
	[ "${SC}" = 200 ]
	[ "$(get_json_path "${BODY}" ".data | length")" = 2 ]
	[ "$(get_json_path "${BODY}" ".data[0].cond")" = "unless" ]
	[ "$(get_json_path "${BODY}" ".data[0].cond_test")" = "{ src 10.0.0.0/8 }" ]
	[ "$(get_json_path "${BODY}" ".data[1].cond")" = "if" ]
	[ "$(get_json_path "${BODY}" ".data[1].cond_test")" = "{ src 192.168.0.0/16 }" ]
}
