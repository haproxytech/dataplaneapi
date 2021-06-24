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
load "../../libs/get_json_path"
load '../../libs/version'

@test "server_templates: setup" {
    run dpa_curl POST "/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "/backends_post.json"
    assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

@test "server_templates: Add a new server template" {
	run dpa_curl POST "/services/haproxy/configuration/server_templates?backend=test_backend&force_reload=true&version=$(version)" "../server_templates/post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

@test "server_templates: Return one server template" {
    run dpa_curl GET "/services/haproxy/configuration/server_templates/srv?backend=test_backend"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local PREFIX; PREFIX=$(get_json_path "$BODY" '.data.prefix')
	[ "${PREFIX}" = "srv" ]
}

@test "servers: Replace a server template" {
	run dpa_curl PUT "/services/haproxy/configuration/server_templates/srv?backend=test_backend&force_reload=true&version=$(version)" "../server_templates/put.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200
}

@test "servers: Delete a server template" {
	run dpa_curl DELETE "/services/haproxy/configuration/server_templates/srv?backend=test_backend&force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 204
}
