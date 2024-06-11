#!/usr/bin/env bats
#
# Copyright 2022 HAProxy Technologies
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
load '../../libs/haproxy_config_setup'
load '../../libs/haproxy_version'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "http_error_rules: Return an array of all HTTP Error Rules from defaults" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_get "$_ERR_RULES_BASE_PATH" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 200
	assert_equal 1 "$(get_json_path "$BODY" ". | length")"
	assert_equal "$(get_json_path "$BODY" ".[0].type")" "status"
	assert_equal "$(get_json_path "$BODY" ".[0].status")" 503
	assert_equal "$(get_json_path "$BODY" ".[0].return_content_type")" "\"application/json\""
	assert_equal "$(get_json_path "$BODY" ".[0].return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".[0].return_content")" "\"Default 503 content\""
}

@test "http_error_rules: Return an array of all HTTP Error Rules from frontend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_get "$_ERR_RULES_BASE_PATH" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 200
	assert_equal 1 "$(get_json_path "$BODY" ". | length")"
	assert_equal "$(get_json_path "$BODY" ".[0].type")" "status"
	assert_equal "$(get_json_path "$BODY" ".[0].status")" 400
}

@test "http_error_rules: Return an array of all HTTP Error Rules from backend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_get "$_ERR_RULES_BASE_PATH" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal 2 "$(get_json_path "$BODY" ". | length")"
	assert_equal "$(get_json_path "$BODY" ".[0].type")" "status"
	assert_equal "$(get_json_path "$BODY" ".[0].status")" 200
	assert_equal "$(get_json_path "$BODY" ".[0].return_content_type")" "\"text/plain\""
	assert_equal "$(get_json_path "$BODY" ".[0].return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".[0].return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".[0].return_hdrs[0].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".[0].return_hdrs[0].fmt")" "value"
	assert_equal "$(get_json_path "$BODY" ".[1].type")" "status"
	assert_equal "$(get_json_path "$BODY" ".[1].status")" 503
	assert_equal "$(get_json_path "$BODY" ".[1].return_content_type")" "application/json"
	assert_equal "$(get_json_path "$BODY" ".[1].return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".[1].return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".[1].return_hdrs[0].name")" "Additional-Header"
	assert_equal "$(get_json_path "$BODY" ".[1].return_hdrs[0].fmt")" "value1"
	assert_equal "$(get_json_path "$BODY" ".[1].return_hdrs[1].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".[1].return_hdrs[1].fmt")" "value"
}
