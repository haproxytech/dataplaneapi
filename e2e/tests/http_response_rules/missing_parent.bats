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
load '../../libs/version'

@test "http_response_rules: Fail creating a HTTP Response rule when frontend doesn't exist" {
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/http_response_rules?parent_type=frontend&parent_name=ghost&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/http_response_rules/unless.json")
	[ "${SC}" = 400 ]
}

@test "http_response_rules: Fail creating a HTTP Response rule when backend doesn't exist" {
	read -r SC _ < <(auth_curl POST "/v2/services/haproxy/configuration/http_response_rules?parent_type=backend&parent_name=ghost&force_reload=true&version=$(version)" "@${E2E_DIR}/tests/http_response_rules/unless.json")
	[ "${SC}" = 400 ]
}
