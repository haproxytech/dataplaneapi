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
load '../../libs/version'

setup() {
	run dpa_curl POST "/services/haproxy/configuration/frontends?force_reload=true&version=$(version)" "/frontends_post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
	run dpa_curl POST "/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "/backends_post.json"
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

@test "http_request_rules: Add a new HTTP Request Rule to frontend" {
	run dpa_curl POST "/services/haproxy/configuration/http_request_rules?parent_type=frontend&parent_name=test_frontend&force_reload=true&version=$(version)" "../http_request_rules/unless.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

@test "http_request_rules: Add a new HTTP Request Rule to backend" {
	run dpa_curl POST "/services/haproxy/configuration/http_request_rules?parent_type=backend&parent_name=test_backend&force_reload=true&version=$(version)" "../http_request_rules/unless.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}
