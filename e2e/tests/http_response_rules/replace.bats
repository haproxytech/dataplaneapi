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

load '../../libs/dataplaneapi'
load '../../libs/get_json_path'
load '../../libs/version'

setup() {
	# creating frontend and related HTTP Response rule
	run dpa_curl POST "/services/haproxy/configuration/frontends?force_reload=true&version=$(version)" "/frontends_post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
	run dpa_curl POST "/services/haproxy/configuration/http_response_rules?parent_type=frontend&parent_name=test_frontend&force_reload=true&version=$(version)" "../http_response_rules/unless.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
	# creating backend and related HTTP Response rule
	run dpa_curl POST "/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "/backends_post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
	run dpa_curl POST "/services/haproxy/configuration/http_response_rules?parent_type=backend&parent_name=test_backend&force_reload=true&version=$(version)" "../http_response_rules/unless.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

teardown() {
	run dpa_curl DELETE "/services/haproxy/configuration/frontends/test_frontend?force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 204
	run dpa_curl DELETE "/services/haproxy/configuration/backends/test_backend?force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 204
}

@test "http_response_rules: Replace a HTTP Response Rule of frontend" {
	run dpa_curl PUT "/services/haproxy/configuration/http_response_rules/0?parent_type=frontend&parent_name=test_frontend&force_reload=true&version=$(version)" "../http_response_rules/put.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200
	[ "$(get_json_path "${BODY}" ".cond")" = "if" ]
}

@test "http_response_rules: Replace a HTTP Response Rule of backend" {
	run dpa_curl PUT "/services/haproxy/configuration/http_response_rules/0?parent_type=backend&parent_name=test_backend&force_reload=true&version=$(version)" "../http_response_rules/put.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200
	[ "$(get_json_path "${BODY}" ".cond")" = "if" ]
}
