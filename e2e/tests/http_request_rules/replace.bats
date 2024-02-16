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

@test "http_request_rules: Replace a HTTP Request Rule of frontend" {
  resource_put "$_REQ_RULES_BASE_PATH/0" "data/put.json" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_REQ_RULES_BASE_PATH/0" "parent_type=frontend&parent_name=test_frontend&force_reload=true"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".hdr_format")" "%T"
	assert_equal "$(get_json_path "$BODY" ".hdr_name")" "X-Haproxy-Current-Date"
	assert_equal "$(get_json_path "$BODY" ".type")" "add-header"
}

@test "http_request_rules: Replace a HTTP Request Rule of backend" {
	resource_put "$_REQ_RULES_BASE_PATH/0" "data/put.json" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 200

	resource_get "$_REQ_RULES_BASE_PATH/0" "parent_type=backend&parent_name=test_backend&force_reload=true"
	assert_equal "$SC" 200
	assert_equal "$(get_json_path "$BODY" ".cond")" "if"
	assert_equal "$(get_json_path "$BODY" ".cond_test")" "{ src 192.168.0.0/16 }"
	assert_equal "$(get_json_path "$BODY" ".hdr_format")" "%T"
	assert_equal "$(get_json_path "$BODY" ".hdr_name")" "X-Haproxy-Current-Date"
	assert_equal "$(get_json_path "$BODY" ".type")" "add-header"
}
