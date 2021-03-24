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

@test "sites: Add a site" {
	run dpa_curl POST "/services/haproxy/sites?force_reload=true&version=$(version)" "../sites/post.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 201
}

@test "sites: Return a site" {
	run dpa_curl GET "/services/haproxy/sites/test_site"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data.name')
	[ "${ACTUAL}" = "test_site" ]
}

@test "sites: Replace a site" {
	run dpa_curl PUT "/services/haproxy/sites/test_site?force_reload=true&version=$(version)" "../sites/put.json"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200
}

@test "sites: Return an array of sites" {
	run dpa_curl GET "/services/haproxy/sites?force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 200

	local ACTUAL; ACTUAL=$(get_json_path "$BODY" '.data[0].name')
	[ "${ACTUAL}" = "test_site" ]
}

@test "sites: Delete a site" {
	run dpa_curl DELETE "/services/haproxy/sites/test_site?force_reload=true&version=$(version)"
	assert_success

	dpa_curl_status_body '$output'
	assert_equal $SC 204
}
