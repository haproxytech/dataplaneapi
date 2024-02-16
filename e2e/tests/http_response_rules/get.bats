#!/usr/bin/env bats
#
# Copyright 2021 HAProxy Technologies
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
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "http_response_rules: Return one HTTP Response Rule from frontend" {
  resource_get "$_RES_RULES_BASE_PATH/0" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".type")" "add-header"
	assert_equal "$(get_json_path "$BODY" ".hdr_name")" "X-Add-Frontend"
	assert_equal "$(get_json_path "$BODY" ".cond")" "unless"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 192.168.0.0/16 }"

	resource_get "$_RES_RULES_BASE_PATH/1" "parent_type=frontend&parent_name=test_frontend"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".type")" "del-header"
	assert_equal "$(get_json_path "$BODY" ".hdr_name")" "X-Del-Frontend"
	assert_equal "$(get_json_path "$BODY" ".cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 10.1.0.0/16 }"
}

@test "http_response_rules: Return one HTTP Response Rule from backend" {
  resource_get "$_RES_RULES_BASE_PATH/0" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".type")" "add-header"
	assert_equal "$(get_json_path "$BODY" ".hdr_name")" "X-Add-Backend"
	assert_equal "$(get_json_path "$BODY" ".cond")" "unless"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 192.168.0.0/16 }"

	resource_get "$_RES_RULES_BASE_PATH/1" "parent_type=backend&parent_name=test_backend"
	assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" ".type")" "del-header"
	assert_equal "$(get_json_path "$BODY" ".hdr_name")" "X-Del-Backend"
	assert_equal "$(get_json_path "$BODY" ".cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 10.1.0.0/16 }"
}
