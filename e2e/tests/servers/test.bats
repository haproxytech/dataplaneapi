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

@test "servers: setup" {
	run dpa_curl POST "/services/haproxy/configuration/backends?force_reload=true&version=$(version)" "/backends_post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

@test "servers: Add a new server" {
	run dpa_curl POST "/services/haproxy/configuration/servers?backend=test_backend&force_reload=true&version=$(version)" "../servers/post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

@test "servers: Return one server" {
	run dpa_curl GET "/services/haproxy/configuration/servers/test_server?backend=test_backend"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local NAME; NAME=$(get_json_path "$BODY" '.data.name')
	[ "${NAME}" = "test_server" ]
}

@test "servers: Replace a server" {
	run dpa_curl PUT "/services/haproxy/configuration/servers/test_server?backend=test_backend&force_reload=true&version=$(version)" "../servers/put.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200
}

@test "servers: Return an array of servers" {
	run dpa_curl GET "/services/haproxy/configuration/servers?backend=test_backend"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data[0].name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Return an array of runtime servers' settings" {
	run dpa_curl GET "/services/haproxy/runtime/servers?backend=test_backend"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.[0].name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Replace server transient settings" {
	run dpa_curl PUT "/services/haproxy/runtime/servers/test_server?backend=test_backend" "../servers/transient.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Return one server runtime settings" {
	run dpa_curl GET "/services/haproxy/runtime/servers/test_server?backend=test_backend"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.name')
	[ "${ACTUAL}" = "test_server" ]
}

@test "servers: Delete a server" {
	run dpa_curl DELETE "/services/haproxy/configuration/servers/test_server?backend=test_backend&force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 204
}

@test "servers: teardown" {
	run dpa_curl DELETE "/services/haproxy/configuration/backends/test_backend?force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 204
}
