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

@test "http_error_rules: Add a new HTTP Error Rule to frontend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_post "$_ERR_RULES_BASE_PATH/0" "data/post.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 201

	resource_get "$_ERR_RULES_BASE_PATH/0" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".status")" 425
}

@test "http_error_rules: Add a new HTTP Error Rule to backend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

	resource_post "$_ERR_RULES_BASE_PATH/0" "data/post.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 201

	resource_get "$_ERR_RULES_BASE_PATH/0" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".status")" 425
}

@test "http_error_rules: Add a new HTTP Error Rule to defaults" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

    resource_post "$_ERR_RULES_BASE_PATH/0" "data/post.json" "parent_type=defaults&parent_name=mydefaults&force_reload=true"
	assert_equal "$SC" 201

	resource_get "$_ERR_RULES_BASE_PATH/0" "parent_type=defaults&parent_name=mydefaults"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".status")" 425
}

@test "http_error_rules: Fail to add a new HTTP Error Rule without status code to backend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

    resource_post "$_ERR_RULES_BASE_PATH/0" "data/invalid_no_status_code.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 422
}

@test "http_error_rules: Fail to add a new HTTP Error Rule with unsupported status code to frontend" {
	haproxy_version_ge $_ERR_SUPPORTED_HAPROXY_VERSION || skip "requires HAProxy $_ERR_SUPPORTED_HAPROXY_VERSION+"

    resource_post "$_ERR_RULES_BASE_PATH/0" "data/invalid_bad_status_code.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 422
}
