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

@test "http_error_rules: Replace a HTTP Error Rule of frontend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/0" "data/put.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_ERR_RULES_BASE_PATH/0" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 429
	assert_equal "$(get_json_path "$BODY" ".return_content_type")" "application/json"
	assert_equal "$(get_json_path "$BODY" ".return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].name")" "Additional-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].fmt")" "value1"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].fmt")" "value"
}

@test "http_error_rules: Replace a HTTP Error Rule of backend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/1" "data/put.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_ERR_RULES_BASE_PATH/1" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 429
	assert_equal "$(get_json_path "$BODY" ".return_content_type")" "application/json"
	assert_equal "$(get_json_path "$BODY" ".return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].name")" "Additional-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].fmt")" "value1"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].fmt")" "value"
}

@test "http_error_rules: Replace a HTTP Error Rule of defaults" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/0" "data/put.json" "parent_type=defaults&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_ERR_RULES_BASE_PATH/0" "parent_type=defaults"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 429
	assert_equal "$(get_json_path "$BODY" ".return_content_type")" "application/json"
	assert_equal "$(get_json_path "$BODY" ".return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].name")" "Additional-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].fmt")" "value1"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].fmt")" "value"
}

@test "http_error_rules: Fail to replace a HTTP Error rule when backend does not exist" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/0" "data/put.json" "parent_type=backend&parent_name=ghost&force_reload=true"
	assert_equal "$SC" 400
}

@test "http_error_rules: Fail to replace a frontend HTTP Error Rule with no status code" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/0" "data/invalid_no_status_code.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 422
}

@test "http_error_rules: Fail to replace a backend HTTP Error Rule with unsupported status code" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/0" "data/invalid_bad_status_code.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 422
}

@test "http_error_rules: Fail to replace a frontend HTTP Error Rule that does not exist" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/1" "data/put.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 404
}

@test "http_error_rules: Fail to replace a backend HTTP Error Rule that does not exist" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_put "$_ERR_RULES_BASE_PATH/1000" "data/put.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 404
}
