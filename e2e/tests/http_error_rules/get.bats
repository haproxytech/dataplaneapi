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

@test "http_error_rules: Return one HTTP Error Rule from frontend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	PARENT_NAME="test_frontend"
	resource_get "$_FRONTEND_BASE_PATH/$PARENT_NAME/http_error_rules/0"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 400
}

@test "http_error_rules: Return one HTTP Error Rule from backend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	PARENT_NAME="test_backend"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/0"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 200
	assert_equal "$(get_json_path "$BODY" ".return_content_type")" "\"text/plain\""
	assert_equal "$(get_json_path "$BODY" ".return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].fmt")" "value"

	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/1"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 503
	assert_equal "$(get_json_path "$BODY" ".return_content_type")" "application/json"
	assert_equal "$(get_json_path "$BODY" ".return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".return_content")" "\"My content\""
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].name")" "Additional-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[0].fmt")" "value1"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].name")" "Some-Header"
	assert_equal "$(get_json_path "$BODY" ".return_hdrs[1].fmt")" "value"
}

@test "http_error_rules: Return one HTTP Error Rule from defaults" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	PARENT_NAME="mydefaults"
	resource_get "$_DEFAULTS_BASE_PATH/$PARENT_NAME/http_error_rules/0"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "status"
	assert_equal "$(get_json_path "$BODY" ".status")" 503
	assert_equal "$(get_json_path "$BODY" ".return_content_type")" "\"application/json\""
	assert_equal "$(get_json_path "$BODY" ".return_content_format")" "string"
	assert_equal "$(get_json_path "$BODY" ".return_content")" "\"Default 503 content\""
}

@test "http_error_rules: Fail to return a HTTP Error Rule when backend does not exist" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	PARENT_NAME="ghost"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/0"
	assert_equal "$SC" 400
}

@test "http_error_rules: Fail to return a backend HTTP Error Rule that does not exist" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	PARENT_NAME="test_backend"
	resource_get "$_BACKEND_BASE_PATH/$PARENT_NAME/http_error_rules/1000" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 404
}
